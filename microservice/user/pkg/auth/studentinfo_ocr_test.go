package auth

import (
	"testing"
	"time"

	"github.com/spf13/viper"
)

func TestCCNUCaptchaAutoRetry(t *testing.T) {
	key := "ccnu_login.captcha_auto_retry"

	t.Run("default when unset", func(t *testing.T) {
		viper.Set(key, nil)
		viper.SetDefault(key, nil)
		viper.Reset()
		if got := ccnuCaptchaAutoRetry(); got != 3 {
			t.Fatalf("expected default retries to be 3, got %d", got)
		}
	})

	t.Run("explicit zero disables auto ocr", func(t *testing.T) {
		viper.Reset()
		viper.Set(key, 0)
		if got := ccnuCaptchaAutoRetry(); got != 0 {
			t.Fatalf("expected retries to be 0, got %d", got)
		}
	})

	t.Run("negative retries clamp to zero", func(t *testing.T) {
		viper.Reset()
		viper.Set(key, -1)
		if got := ccnuCaptchaAutoRetry(); got != 0 {
			t.Fatalf("expected retries to be 0, got %d", got)
		}
	})

	t.Run("positive retries are respected", func(t *testing.T) {
		viper.Reset()
		viper.Set(key, 5)
		if got := ccnuCaptchaAutoRetry(); got != 5 {
			t.Fatalf("expected retries to be 5, got %d", got)
		}
	})

	viper.Reset()
}

func TestCCNUCaptchaAutoTimeout(t *testing.T) {
	key := "ccnu_login.captcha_auto_timeout_ms"

	t.Run("default when unset", func(t *testing.T) {
		viper.Reset()
		if got := ccnuCaptchaAutoTimeout(); got != 3*time.Second {
			t.Fatalf("expected default timeout to be 3s, got %s", got)
		}
	})

	t.Run("configured timeout is respected", func(t *testing.T) {
		viper.Reset()
		viper.Set(key, 1500)
		if got := ccnuCaptchaAutoTimeout(); got != 1500*time.Millisecond {
			t.Fatalf("expected timeout to be 1500ms, got %s", got)
		}
	})

	viper.Reset()
}
