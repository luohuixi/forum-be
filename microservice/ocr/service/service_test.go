package service

import "testing"

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
