package auth

import (
	"context"
	"errors"
	"strings"

	pb "forum-ocr/proto"
	ocrclient "forum/client"

	"github.com/spf13/viper"
)

func ccnuCaptchaAutoRetry() int {
	retries := viper.GetInt("ccnu_login.captcha_auto_retry")
	if retries < 0 {
		return 0
	}
	if retries == 0 {
		return 3
	}
	return retries
}

func ccnuAutoOCRAvailable() bool {
	return ocrclient.OCRClient != nil
}

func (m *ccnuLoginManager) recognizeCaptcha(imageBase64 string) (string, error) {
	if ocrclient.OCRClient == nil {
		return "", errors.New("ocr rpc client is not initialized")
	}

	resp, err := ocrclient.OCRClient.RecognizeCaptcha(context.Background(), &pb.RecognizeCaptchaRequest{
		ImageBase64: imageBase64,
	})
	if err != nil {
		return "", err
	}

	text := sanitizeCaptchaText(strings.TrimSpace(resp.GetText()))
	if text == "" {
		return "", errors.New("ocr response text is empty")
	}
	return text, nil
}
