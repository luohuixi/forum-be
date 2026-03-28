package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
)

const helperScriptTemplate = `from modelscope.pipelines import pipeline
from modelscope.utils.constant import Tasks
import json
import re
import sys

ocr = pipeline(Tasks.ocr_recognition, model=%q)
raw = ''.join(ocr(sys.argv[1]).get('text') or [])
captcha = re.sub(r'[^A-Za-z0-9]', '', raw)[-4:]
print(json.dumps({'raw': raw, 'captcha': captcha}, ensure_ascii=False))
`

type OCRService struct {
	engine *modelscopeEngine
}

type modelscopeEngine struct {
	pythonBin  string
	workspace  string
	runtimeDir string
	model      string
	timeout    time.Duration

	helperOnce sync.Once
	helperPath string
	helperErr  error
}

type recognizeResult struct {
	Raw     string `json:"raw"`
	Captcha string `json:"captcha"`
}

func New() (*OCRService, error) {
	engine := newModelscopeEngine()
	if err := engine.SelfCheck(context.Background()); err != nil {
		return nil, err
	}
	return &OCRService{engine: engine}, nil
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
	runtimeDir := strings.TrimSpace(viper.GetString("ocr.modelscope.runtime_dir"))
	if runtimeDir == "" {
		runtimeDir = filepath.Join(os.TempDir(), "forum-ocr")
	}

	model := strings.TrimSpace(viper.GetString("ocr.modelscope.model"))
	if model == "" {
		model = defaultModelscopeModel
	}

	timeout := defaultModelscopeTimeout
	if timeoutMS := viper.GetInt("ocr.modelscope.timeout_ms"); timeoutMS > 0 {
		timeout = time.Duration(timeoutMS) * time.Millisecond
	}

	return &modelscopeEngine{
		pythonBin:  defaultModelscopePythonBin,
		workspace:  defaultModelscopeWorkspace,
		runtimeDir: runtimeDir,
		model:      model,
		timeout:    timeout,
	}
}

func (e *modelscopeEngine) SelfCheck(parent context.Context) error {
	if e == nil {
		return errors.New("ocr engine is nil")
	}
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

	timeout := defaultSelfCheckTimeout
	if timeoutMS := viper.GetInt("ocr.modelscope.self_check_timeout_ms"); timeoutMS > 0 {
		timeout = time.Duration(timeoutMS) * time.Millisecond
	}

	ctx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()

	startedAt := time.Now()
	cmd := exec.CommandContext(ctx, e.pythonBin, "-c", fmt.Sprintf(
		"from modelscope.pipelines import pipeline\nfrom modelscope.utils.constant import Tasks\npipeline(Tasks.ocr_recognition, model=%q)\nprint('ok')\n",
		e.model,
	))
	if e.workspace != "" {
		cmd.Dir = e.workspace
	}

	output, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("ocr self-check timed out after %s", timeout)
	}
	if err != nil {
		return fmt.Errorf("ocr self-check failed: %w: %s", err, strings.TrimSpace(string(output)))
	}

	logger.Info(fmt.Sprintf("OCR startup self-check passed in %s", time.Since(startedAt)))
	return nil
}

func (e *modelscopeEngine) Recognize(ctx context.Context, imageBase64 string) (string, error) {
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

	result, err := e.runHelper(ctx, imagePath)
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

func (e *modelscopeEngine) ensureHelperScript() error {
	e.helperOnce.Do(func() {
		if err := os.MkdirAll(e.runtimeDir, 0o755); err != nil {
			e.helperErr = err
			return
		}

		helperPath := filepath.Join(e.runtimeDir, "ocr_once.py")
		script := fmt.Sprintf(helperScriptTemplate, e.model)
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

func (e *modelscopeEngine) runHelper(parent context.Context, imagePath string) (recognizeResult, error) {
	ctx, cancel := context.WithTimeout(parent, e.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, e.pythonBin, e.helperPath, imagePath)
	if e.workspace != "" {
		cmd.Dir = e.workspace
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return recognizeResult{}, fmt.Errorf("modelscope ocr failed: %w: %s", err, strings.TrimSpace(string(output)))
	}

	result, err := parseRecognizeResult(output)
	if err != nil {
		return recognizeResult{}, fmt.Errorf("parse ocr result failed: %w: %s", err, strings.TrimSpace(string(output)))
	}
	return result, nil
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
