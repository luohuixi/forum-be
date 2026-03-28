package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	pb "forum-ocr/proto"
	ocrclient "forum/client"

	"github.com/spf13/viper"
	microclient "go-micro.dev/v4/client"
)

const defaultCCNUCaptchaAutoTimeout = 3 * time.Second

func ccnuCaptchaAutoRetry() int {
	if !viper.IsSet("ccnu_login.captcha_auto_retry") {
		return 3
	}

	retries := viper.GetInt("ccnu_login.captcha_auto_retry")
	if retries < 0 {
		return 0
	}
	return retries
}

func ccnuAutoOCRAvailable() bool {
	return ocrclient.OCRClient != nil
}

func ccnuCaptchaAutoTimeout() time.Duration {
	timeout := defaultCCNUCaptchaAutoTimeout
	if timeoutMS := viper.GetInt("ccnu_login.captcha_auto_timeout_ms"); timeoutMS > 0 {
		timeout = time.Duration(timeoutMS) * time.Millisecond
	}
	return timeout
}

func (m *ccnuLoginManager) recognizeCaptcha(imageBase64 string, timeout time.Duration) (string, error) {
	if ocrclient.OCRClient == nil {
		return "", errors.New("ocr rpc client is not initialized")
	}
	if timeout <= 0 {
		return "", context.DeadlineExceeded
	}

	resp, err := ocrclient.OCRClient.RecognizeCaptcha(context.Background(), &pb.RecognizeCaptchaRequest{
		ImageBase64: imageBase64,
	}, microclient.WithRequestTimeout(timeout))
	if err != nil {
		return "", err
	}

	text := sanitizeCaptchaText(strings.TrimSpace(resp.GetText()))
	if text == "" {
		return "", errors.New("ocr response text is empty")
	}
	return text, nil
}
