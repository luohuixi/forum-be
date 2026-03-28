package auth

import (
	"testing"

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
