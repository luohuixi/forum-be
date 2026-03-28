package service

import (
	"testing"

	"github.com/spf13/viper"
)

func TestParseRecognizeResult(t *testing.T) {
	output := []byte("warming up\n{\"raw\":\"a-b12\",\"captcha\":\"ab12\"}\n")

	result, err := parseRecognizeResult(output)
	if err != nil {
		t.Fatalf("parseRecognizeResult returned error: %v", err)
	}
	if result.Raw != "a-b12" {
		t.Fatalf("unexpected raw result: %q", result.Raw)
	}
	if result.Captcha != "ab12" {
		t.Fatalf("unexpected captcha result: %q", result.Captcha)
	}
}

func TestSanitizeCaptchaText(t *testing.T) {
	text := sanitizeCaptchaText("!ab-12XYZ")
	if text != "2XYZ" {
		t.Fatalf("unexpected sanitized captcha: %q", text)
	}
}

func TestDecodeBase64ImageDataURI(t *testing.T) {
	imageBytes, err := decodeBase64Image("data:image/png;base64,aGVsbG8=")
	if err != nil {
		t.Fatalf("decodeBase64Image returned error: %v", err)
	}
	if string(imageBytes) != "hello" {
		t.Fatalf("unexpected decoded payload: %q", string(imageBytes))
	}
}

func TestNewModelscopeEngineUsesConfiguredValues(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)

	t.Setenv(envModelscopePythonBin, "/tmp/python3")
	t.Setenv(envModelscopeWorkspace, "/tmp/workspace")
	t.Setenv(envModelscopeRuntimeDir, "/tmp/runtime")
	t.Setenv(envModelscopeSkipSelfChk, "true")
	viper.Set("ocr.modelscope.model", "custom/model")
	viper.Set("ocr.modelscope.timeout_ms", 1234)

	engine := newModelscopeEngine()

	if engine.pythonBin != "/tmp/python3" {
		t.Fatalf("unexpected python bin: %q", engine.pythonBin)
	}
	if engine.workspace != "/tmp/workspace" {
		t.Fatalf("unexpected workspace: %q", engine.workspace)
	}
	if engine.runtimeDir != "/tmp/runtime" {
		t.Fatalf("unexpected runtime dir: %q", engine.runtimeDir)
	}
	if engine.model != "custom/model" {
		t.Fatalf("unexpected model: %q", engine.model)
	}
	if engine.timeout.Milliseconds() != 1234 {
		t.Fatalf("unexpected timeout: %s", engine.timeout)
	}
	if !engine.skipSelfCheck {
		t.Fatal("expected skip self-check to be enabled")
	}
}
