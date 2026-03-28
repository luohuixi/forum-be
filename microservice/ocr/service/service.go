package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	pb "forum-ocr/proto"
	logger "forum/log"
	"forum/pkg/errno"

	"github.com/spf13/viper"
)

const (
	defaultModelscopePythonBin = "/opt/forum-ocr/.venvs/modelscope-ocr/bin/python"
	defaultModelscopeWorkspace = "/opt/forum-ocr/workspace"
	defaultModelscopeModel     = "xiaolv/ocr_small"
	defaultModelscopeTimeout   = 180 * time.Second
	defaultSelfCheckTimeout    = 120 * time.Second

	envModelscopePythonBin   = "FORUM_OCR_OCR_MODELSCOPE_PYTHON_BIN"
	envModelscopeWorkspace   = "FORUM_OCR_OCR_MODELSCOPE_WORKSPACE"
	envModelscopeRuntimeDir  = "FORUM_OCR_OCR_MODELSCOPE_RUNTIME_DIR"
	envModelscopeSkipSelfChk = "FORUM_OCR_OCR_MODELSCOPE_SKIP_SELF_CHECK"
)

const workerScriptTemplate = `from modelscope.pipelines import pipeline
from modelscope.utils.constant import Tasks
import json
import re
import sys

try:
    ocr = pipeline(Tasks.ocr_recognition, model=%q)
    print(json.dumps({"status": "ready"}, ensure_ascii=False), flush=True)
except Exception as exc:
    print(json.dumps({"status": "error", "error": str(exc)}, ensure_ascii=False), flush=True)
    raise

for line in sys.stdin:
    line = line.strip()
    if not line:
        continue
    try:
        payload = json.loads(line)
        raw = ''.join(ocr(payload["image_path"]).get("text") or [])
        captcha = re.sub(r'[^A-Za-z0-9]', '', raw)[-4:]
        print(json.dumps({"ok": True, "raw": raw, "captcha": captcha}, ensure_ascii=False), flush=True)
    except Exception as exc:
        print(json.dumps({"ok": False, "error": str(exc)}, ensure_ascii=False), flush=True)
`

type OCRService struct {
	engine *modelscopeEngine
}

type modelscopeEngine struct {
	pythonBin     string
	workspace     string
	runtimeDir    string
	model         string
	timeout       time.Duration
	skipSelfCheck bool

	helperOnce sync.Once
	helperPath string
	helperErr  error
	// TODO 优化锁粒度，提高并发度
	mu     sync.Mutex
	worker *modelscopeWorker
}

type modelscopeWorker struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Reader
	stderr *lockedBuffer
	waitCh chan error
}

type lockedBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

type recognizeResult struct {
	Raw     string `json:"raw"`
	Captcha string `json:"captcha"`
}

type workerMessage struct {
	Status  string `json:"status"`
	Ok      bool   `json:"ok"`
	Error   string `json:"error"`
	Raw     string `json:"raw"`
	Captcha string `json:"captcha"`
}

type lineResult struct {
	line string
	err  error
}

func New() (*OCRService, error) {
	engine := newModelscopeEngine()
	if err := engine.SelfCheck(context.Background()); err != nil {
		return nil, err
	}
	return &OCRService{engine: engine}, nil
}

func (s *OCRService) Close() error {
	if s == nil || s.engine == nil {
		return nil
	}
	return s.engine.Close()
}

func (s *OCRService) RecognizeCaptcha(ctx context.Context, req *pb.RecognizeCaptchaRequest, resp *pb.RecognizeCaptchaResponse) error {
	logger.Info("OCRService RecognizeCaptcha")

	imageBase64 := strings.TrimSpace(req.GetImageBase64())
	if imageBase64 == "" {
		return errno.ServerErr(errno.ErrBadRequest, "image_base64 is required")
	}

	text, err := s.engine.Recognize(ctx, imageBase64)
	if err != nil {
		return errno.ServerErr(errno.InternalServerError, err.Error())
	}
	resp.Text = text
	return nil
}

func newModelscopeEngine() *modelscopeEngine {
	pythonBin := defaultModelscopePythonBin
	if override := strings.TrimSpace(os.Getenv(envModelscopePythonBin)); override != "" {
		pythonBin = override
	}

	workspace := defaultModelscopeWorkspace
	if override := strings.TrimSpace(os.Getenv(envModelscopeWorkspace)); override != "" {
		workspace = override
	}

	runtimeDir := filepath.Join(os.TempDir(), "forum-ocr")
	if override := strings.TrimSpace(os.Getenv(envModelscopeRuntimeDir)); override != "" {
		runtimeDir = override
	}

	model := strings.TrimSpace(viper.GetString("ocr.modelscope.model"))
	if model == "" {
		model = defaultModelscopeModel
	}

	timeout := defaultModelscopeTimeout
	if timeoutMS := viper.GetInt("ocr.modelscope.timeout_ms"); timeoutMS > 0 {
		timeout = time.Duration(timeoutMS) * time.Millisecond
	}

	skipSelfCheck := false
	if override := strings.TrimSpace(os.Getenv(envModelscopeSkipSelfChk)); override != "" {
		if parsed, err := strconv.ParseBool(override); err == nil {
			skipSelfCheck = parsed
		} else {
			logger.Error(fmt.Sprintf("invalid %s value: %q", envModelscopeSkipSelfChk, override))
		}
	}

	return &modelscopeEngine{
		pythonBin:     pythonBin,
		workspace:     workspace,
		runtimeDir:    runtimeDir,
		model:         model,
		timeout:       timeout,
		skipSelfCheck: skipSelfCheck,
	}
}

func (e *modelscopeEngine) SelfCheck(parent context.Context) error {
	if e == nil {
		return errors.New("ocr engine is nil")
	}
	if e.skipSelfCheck {
		logger.Info("OCR startup self-check skipped by config")
		return nil
	}
	if err := e.validateEnvironment(); err != nil {
		return err
	}

	timeout := defaultSelfCheckTimeout
	if timeoutMS := viper.GetInt("ocr.modelscope.self_check_timeout_ms"); timeoutMS > 0 {
		timeout = time.Duration(timeoutMS) * time.Millisecond
	}

	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()

	startedAt := time.Now()
	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.ensureWorkerLocked(ctx); err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("OCR startup self-check passed in %s", time.Since(startedAt)))
	return nil
}

func (e *modelscopeEngine) Close() error {
	if e == nil {
		return nil
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.stopWorkerLocked()
}

func (e *modelscopeEngine) Recognize(_ context.Context, imageBase64 string) (string, error) {
	if e == nil {
		return "", errors.New("ocr engine is nil")
	}
	if err := e.ensureHelperScript(); err != nil {
		return "", err
	}

	imagePath, cleanup, err := e.writeImage(imageBase64)
	if err != nil {
		return "", err
	}
	defer cleanup()

	ctx, cancel := e.newWorkerRequestContext()
	defer cancel()

	e.mu.Lock()
	defer e.mu.Unlock()

	if err := e.ensureWorkerReadyForRequestLocked(); err != nil {
		return "", err
	}

	result, err := e.runWorkerRequestLocked(ctx, imagePath)
	if err != nil {
		return "", err
	}

	text := sanitizeCaptchaText(result.Captcha)
	if text == "" {
		text = sanitizeCaptchaText(result.Raw)
	}
	if text == "" {
		return "", errors.New("ocr response text is empty")
	}
	return text, nil
}

func (e *modelscopeEngine) ensureWorkerReadyForRequestLocked() error {
	startCtx, cancel := e.newWorkerStartContext()
	defer cancel()
	return e.ensureWorkerLocked(startCtx)
}

func (e *modelscopeEngine) workerStartTimeout() time.Duration {
	timeout := defaultSelfCheckTimeout
	if timeoutMS := viper.GetInt("ocr.modelscope.self_check_timeout_ms"); timeoutMS > 0 {
		timeout = time.Duration(timeoutMS) * time.Millisecond
	}
	if e.timeout > timeout {
		timeout = e.timeout
	}
	return timeout
}

func (e *modelscopeEngine) newWorkerStartContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), e.workerStartTimeout())
}

func (e *modelscopeEngine) newWorkerRequestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), e.timeout)
}

func (e *modelscopeEngine) validateEnvironment() error {
	if strings.TrimSpace(e.pythonBin) == "" {
		return errors.New("ocr python_bin is empty")
	}
	if _, err := exec.LookPath(e.pythonBin); err != nil {
		return fmt.Errorf("ocr python_bin is not executable: %w", err)
	}
	if e.workspace != "" {
		if err := os.MkdirAll(e.workspace, 0o755); err != nil {
			return fmt.Errorf("create ocr workspace failed: %w", err)
		}
	}
	if err := e.ensureHelperScript(); err != nil {
		return fmt.Errorf("prepare helper script failed: %w", err)
	}
	return nil
}

func (e *modelscopeEngine) ensureWorkerLocked(ctx context.Context) error {
	if err := e.validateEnvironment(); err != nil {
		return err
	}

	if err := e.checkWorkerExitedLocked(); err != nil {
		logger.Error(fmt.Sprintf("OCR worker exited unexpectedly: %v", err))
	}

	if e.worker != nil {
		return nil
	}
	return e.startWorkerLocked(ctx)
}

func (e *modelscopeEngine) checkWorkerExitedLocked() error {
	if e.worker == nil || e.worker.waitCh == nil {
		return nil
	}
	select {
	case err, ok := <-e.worker.waitCh:
		if !ok {
			err = nil
		}
		stderrText := e.worker.stderr.String()
		e.worker = nil
		if err == nil {
			if stderrText != "" {
				return fmt.Errorf("ocr worker stopped: %s", stderrText)
			}
			return errors.New("ocr worker stopped")
		}
		if stderrText != "" {
			return fmt.Errorf("%w: %s", err, stderrText)
		}
		return err
	default:
		return nil
	}
}

func (e *modelscopeEngine) startWorkerLocked(ctx context.Context) error {
	cmd := exec.Command(e.pythonBin, "-u", e.helperPath)
	if e.workspace != "" {
		cmd.Dir = e.workspace
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("create ocr worker stdin failed: %w", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("create ocr worker stdout failed: %w", err)
	}
	stderr := &lockedBuffer{}
	cmd.Stderr = stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start ocr worker failed: %w", err)
	}

	worker := &modelscopeWorker{
		cmd:    cmd,
		stdin:  stdin,
		stdout: bufio.NewReader(stdout),
		stderr: stderr,
		waitCh: make(chan error, 1),
	}
	go func() {
		worker.waitCh <- cmd.Wait()
		close(worker.waitCh)
	}()

	e.worker = worker

	message, err := e.readWorkerMessageLocked(ctx)
	if err != nil {
		stderrText := worker.stderr.String()
		_ = e.stopWorkerLocked()
		if stderrText != "" {
			return fmt.Errorf("start ocr worker failed: %w: %s", err, stderrText)
		}
		return fmt.Errorf("start ocr worker failed: %w", err)
	}
	if message.Status != "ready" {
		stderrText := worker.stderr.String()
		_ = e.stopWorkerLocked()
		if message.Error != "" {
			if stderrText != "" {
				return fmt.Errorf("ocr worker ready handshake failed: %s: %s", message.Error, stderrText)
			}
			return fmt.Errorf("ocr worker ready handshake failed: %s", message.Error)
		}
		if stderrText != "" {
			return fmt.Errorf("ocr worker ready handshake failed: %s", stderrText)
		}
		return errors.New("ocr worker ready handshake failed")
	}
	worker.stderr.Reset()
	return nil
}

func (e *modelscopeEngine) stopWorkerLocked() error {
	if e.worker == nil {
		return nil
	}

	worker := e.worker
	e.worker = nil

	if worker.stdin != nil {
		_ = worker.stdin.Close()
	}
	if worker.cmd != nil && worker.cmd.Process != nil {
		_ = worker.cmd.Process.Kill()
	}
	return nil
}

func (e *modelscopeEngine) runWorkerRequestLocked(ctx context.Context, imagePath string) (recognizeResult, error) {
	if e.worker == nil {
		return recognizeResult{}, errors.New("ocr worker is not running")
	}

	reqBody, err := json.Marshal(map[string]string{"image_path": imagePath})
	if err != nil {
		return recognizeResult{}, err
	}
	if _, err := fmt.Fprintln(e.worker.stdin, string(reqBody)); err != nil {
		stderrText := e.worker.stderr.String()
		_ = e.stopWorkerLocked()
		if stderrText != "" {
			return recognizeResult{}, fmt.Errorf("send ocr request failed: %w: %s", err, stderrText)
		}
		return recognizeResult{}, fmt.Errorf("send ocr request failed: %w", err)
	}

	message, err := e.readWorkerMessageLocked(ctx)
	if err != nil {
		stderrText := e.worker.stderr.String()
		_ = e.stopWorkerLocked()
		if errors.Is(err, context.DeadlineExceeded) {
			if stderrText != "" {
				return recognizeResult{}, fmt.Errorf("modelscope ocr timed out after %s: %s", e.timeout, stderrText)
			}
			return recognizeResult{}, fmt.Errorf("modelscope ocr timed out after %s", e.timeout)
		}
		if stderrText != "" {
			return recognizeResult{}, fmt.Errorf("modelscope ocr failed: %w: %s", err, stderrText)
		}
		return recognizeResult{}, fmt.Errorf("modelscope ocr failed: %w", err)
	}
	if !message.Ok {
		if strings.TrimSpace(message.Error) != "" {
			return recognizeResult{}, fmt.Errorf("modelscope ocr failed: %s", strings.TrimSpace(message.Error))
		}
		return recognizeResult{}, errors.New("modelscope ocr failed")
	}

	return recognizeResult{
		Raw:     message.Raw,
		Captcha: message.Captcha,
	}, nil
}

func (e *modelscopeEngine) readWorkerMessageLocked(ctx context.Context) (workerMessage, error) {
	for {
		line, err := e.readLineLocked(ctx)
		if err != nil {
			return workerMessage{}, err
		}

		trimmed := strings.TrimSpace(line)
		if trimmed == "" || !strings.HasPrefix(trimmed, "{") {
			continue
		}

		var message workerMessage
		if err := json.Unmarshal([]byte(trimmed), &message); err != nil {
			continue
		}
		return message, nil
	}
}

func (e *modelscopeEngine) readLineLocked(ctx context.Context) (string, error) {
	if e.worker == nil || e.worker.stdout == nil {
		return "", errors.New("ocr worker stdout is not available")
	}

	resultCh := make(chan lineResult, 1)
	go func(reader *bufio.Reader) {
		line, err := reader.ReadString('\n')
		resultCh <- lineResult{line: line, err: err}
	}(e.worker.stdout)

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case result := <-resultCh:
		if result.err != nil {
			return "", result.err
		}
		return result.line, nil
	}
}

func (e *modelscopeEngine) ensureHelperScript() error {
	e.helperOnce.Do(func() {
		if err := os.MkdirAll(e.runtimeDir, 0o755); err != nil {
			e.helperErr = err
			return
		}

		helperPath := filepath.Join(e.runtimeDir, "ocr_worker.py")
		script := fmt.Sprintf(workerScriptTemplate, e.model)
		if err := os.WriteFile(helperPath, []byte(script), 0o644); err != nil {
			e.helperErr = err
			return
		}
		e.helperPath = helperPath
	})
	return e.helperErr
}

func (e *modelscopeEngine) writeImage(imageBase64 string) (string, func(), error) {
	imageBytes, err := decodeBase64Image(imageBase64)
	if err != nil {
		return "", func() {}, err
	}

	imageDir := filepath.Join(e.runtimeDir, "images")
	if err := os.MkdirAll(imageDir, 0o755); err != nil {
		return "", func() {}, err
	}

	file, err := os.CreateTemp(imageDir, "captcha-*.png")
	if err != nil {
		return "", func() {}, err
	}
	cleanup := func() {
		_ = os.Remove(file.Name())
	}

	if _, err := file.Write(imageBytes); err != nil {
		_ = file.Close()
		cleanup()
		return "", func() {}, err
	}
	if err := file.Close(); err != nil {
		cleanup()
		return "", func() {}, err
	}

	return file.Name(), cleanup, nil
}

func parseRecognizeResult(output []byte) (recognizeResult, error) {
	lines := strings.Split(string(output), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line == "" || !strings.HasPrefix(line, "{") {
			continue
		}

		var result recognizeResult
		if err := json.Unmarshal([]byte(line), &result); err == nil {
			return result, nil
		}
	}
	return recognizeResult{}, errors.New("ocr json payload not found")
}

func decodeBase64Image(raw string) ([]byte, error) {
	raw = strings.TrimSpace(raw)
	if idx := strings.Index(raw, ","); idx >= 0 && strings.Contains(raw[:idx], "base64") {
		raw = raw[idx+1:]
	}

	decoders := []*base64.Encoding{
		base64.StdEncoding,
		base64.RawStdEncoding,
		base64.URLEncoding,
		base64.RawURLEncoding,
	}
	for _, encoding := range decoders {
		decoded, err := encoding.DecodeString(raw)
		if err == nil {
			return decoded, nil
		}
	}
	return nil, errors.New("invalid base64 image payload")
}

func sanitizeCaptchaText(raw string) string {
	var builder strings.Builder
	for _, r := range raw {
		if r >= '0' && r <= '9' || r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' {
			builder.WriteRune(r)
		}
	}

	text := builder.String()
	if len(text) > 4 {
		text = text[len(text)-4:]
	}
	return text
}

func (b *lockedBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *lockedBuffer) String() string {
	if b == nil {
		return ""
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	return strings.TrimSpace(b.buf.String())
}

func (b *lockedBuffer) Reset() {
	if b == nil {
		return
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	b.buf.Reset()
}
