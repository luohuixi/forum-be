package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"hash/adler32"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"path"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	defaultFeishuUploadURL    = "https://open.feishu.cn/open-apis/drive/v1/medias/upload_all"
	defaultFeishuUploadParent = "bitable_image"
	feedbackServiceTimeout    = 15 * time.Second
	maxFeedbackImageSize      = 20 << 20
)

type FeedbackRecordRequest struct {
	TableIdentify string         `json:"table_identify"`
	StudentID     string         `json:"student_id"`
	Content       string         `json:"content"`
	Images        []string       `json:"images"`
	ContactInfo   string         `json:"contact_info"`
	ExtraRecord   map[string]any `json:"extra_record"`
}

type feedbackEnvelope[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Msg     string `json:"msg"`
	Data    T      `json:"data"`
}

type feedbackTableTokenResp struct {
	AccessToken string `json:"access_token"`
}

type feedbackCreateRecordResp struct {
	RecordID string `json:"record_id"`
}

type feedbackTenantTokenResp struct {
	AccessToken string `json:"access_token"`
}

type feishuUploadResp struct {
	FileToken string `json:"file_token"`
}

type feedbackImageFile struct {
	name        string
	contentType string
	data        []byte
}

func CreateFeedbackRecord(ctx context.Context, req FeedbackRecordRequest) error {
	baseURL, err := feedbackServiceBaseURL()
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: feedbackServiceTimeout}
	token, err := getFeedbackTableToken(ctx, client, baseURL, req.TableIdentify)
	if err != nil {
		return err
	}

	_, err = postFeedbackJSON[feedbackCreateRecordResp](ctx, client, baseURL+"/api/v1/sheet/records", token, req)
	return err
}

func UploadFeedbackImage(ctx context.Context, fileHeader *multipart.FileHeader) (string, error) {
	baseURL, err := feedbackServiceBaseURL()
	if err != nil {
		return "", err
	}

	parentNode := strings.TrimSpace(viper.GetString("feedback_service.upload_parent_node"))
	if parentNode == "" {
		return "", fmt.Errorf("feedback_service.upload_parent_node 未配置")
	}

	client := &http.Client{Timeout: feedbackServiceTimeout}
	tenantToken, err := getFeedbackTenantToken(ctx, client, baseURL)
	if err != nil {
		return "", err
	}

	file, err := feedbackImageFromMultipart(fileHeader)
	if err != nil {
		return "", err
	}

	return uploadFeedbackImageToFeishu(ctx, client, tenantToken, parentNode, file)
}

func feedbackServiceBaseURL() (string, error) {
	baseURL := strings.TrimRight(strings.TrimSpace(viper.GetString("feedback_service.base_url")), "/")
	if baseURL == "" {
		return "", fmt.Errorf("feedback_service.base_url 未配置")
	}
	return baseURL, nil
}

func getFeedbackTableToken(ctx context.Context, client *http.Client, baseURL string, tableIdentify string) (string, error) {
	resp, err := postFeedbackJSON[feedbackTableTokenResp](ctx, client, baseURL+"/api/v1/auth/table-config/token", "", map[string]string{
		"table_identify": tableIdentify,
	})
	if err != nil {
		return "", err
	}
	if resp.AccessToken == "" {
		return "", fmt.Errorf("反馈服务未返回访问凭证")
	}
	return resp.AccessToken, nil
}

func getFeedbackTenantToken(ctx context.Context, client *http.Client, baseURL string) (string, error) {
	resp, err := postFeedbackJSON[feedbackTenantTokenResp](ctx, client, baseURL+"/api/v1/auth/tenant/token", "", nil)
	if err != nil {
		return "", err
	}
	if resp.AccessToken == "" {
		return "", fmt.Errorf("反馈服务未返回飞书上传凭证")
	}
	return resp.AccessToken, nil
}

func feedbackImageFromMultipart(fileHeader *multipart.FileHeader) (feedbackImageFile, error) {
	if fileHeader == nil {
		return feedbackImageFile{}, fmt.Errorf("反馈图片不能为空")
	}
	if fileHeader.Size > maxFeedbackImageSize {
		return feedbackImageFile{}, fmt.Errorf("反馈图片超过大小限制")
	}

	src, err := fileHeader.Open()
	if err != nil {
		return feedbackImageFile{}, err
	}
	defer src.Close()

	data, err := io.ReadAll(io.LimitReader(src, maxFeedbackImageSize+1))
	if err != nil {
		return feedbackImageFile{}, err
	}
	if len(data) == 0 {
		return feedbackImageFile{}, fmt.Errorf("反馈图片不能为空")
	}
	if len(data) > maxFeedbackImageSize {
		return feedbackImageFile{}, fmt.Errorf("反馈图片超过大小限制")
	}

	contentType := feedbackImageContentType(fileHeader, data)
	if !strings.HasPrefix(contentType, "image/") {
		return feedbackImageFile{}, fmt.Errorf("反馈图片类型不支持")
	}

	return feedbackImageFile{
		name:        feedbackImageFileName(fileHeader.Filename, contentType),
		contentType: contentType,
		data:        data,
	}, nil
}

func feedbackImageContentType(fileHeader *multipart.FileHeader, data []byte) string {
	contentType := strings.TrimSpace(fileHeader.Header.Get("Content-Type"))
	if idx := strings.Index(contentType, ";"); idx >= 0 {
		contentType = strings.TrimSpace(contentType[:idx])
	}
	if !strings.HasPrefix(contentType, "image/") {
		contentType = http.DetectContentType(data)
	}
	return contentType
}

func feedbackImageFileName(fileName string, contentType string) string {
	fileName = strings.ReplaceAll(fileName, "\\", "/")
	fileName = path.Base(fileName)
	if fileName == "." || fileName == "/" || fileName == "" {
		fileName = "feedback-image"
	}
	if path.Ext(fileName) == "" {
		if extensions, err := mime.ExtensionsByType(contentType); err == nil && len(extensions) > 0 {
			fileName += extensions[0]
		}
	}
	return fileName
}

func uploadFeedbackImageToFeishu(ctx context.Context, client *http.Client, token string, parentNode string, file feedbackImageFile) (string, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	fields := map[string]string{
		"file_name":   file.name,
		"parent_type": stringWithDefault(viper.GetString("feedback_service.upload_parent_type"), defaultFeishuUploadParent),
		"parent_node": parentNode,
		"size":        fmt.Sprintf("%d", len(file.data)),
		"checksum":    fmt.Sprintf("%d", adler32.Checksum(file.data)),
	}
	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			return "", err
		}
	}

	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, escapeMultipartFilename(file.name)))
	header.Set("Content-Type", file.contentType)
	part, err := writer.CreatePart(header)
	if err != nil {
		return "", err
	}
	if _, err := part.Write(file.data); err != nil {
		return "", err
	}
	if err := writer.Close(); err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, stringWithDefault(viper.GetString("feedback_service.feishu_upload_url"), defaultFeishuUploadURL), &body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var envelope feedbackEnvelope[feishuUploadResp]
	if err := json.Unmarshal(respBody, &envelope); err != nil {
		return "", err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices || envelope.Code != 0 {
		message := envelope.Message
		if message == "" {
			message = envelope.Msg
		}
		if message == "" {
			message = string(respBody)
		}
		return "", fmt.Errorf("飞书图片上传失败: status=%d code=%d message=%s", resp.StatusCode, envelope.Code, message)
	}
	if envelope.Data.FileToken == "" {
		return "", fmt.Errorf("飞书图片上传未返回 file_token")
	}

	return envelope.Data.FileToken, nil
}

func postFeedbackJSON[T any](ctx context.Context, client *http.Client, url string, token string, body any) (T, error) {
	var zero T

	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return zero, err
		}
		reader = bytes.NewReader(payload)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, reader)
	if err != nil {
		return zero, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	httpResp, err := client.Do(req)
	if err != nil {
		return zero, err
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return zero, err
	}

	var envelope feedbackEnvelope[T]
	if err := json.Unmarshal(respBody, &envelope); err != nil {
		return zero, err
	}

	if httpResp.StatusCode < http.StatusOK || httpResp.StatusCode >= http.StatusMultipleChoices || envelope.Code != 0 {
		message := envelope.Message
		if message == "" {
			message = envelope.Msg
		}
		if message == "" {
			message = string(respBody)
		}
		return zero, fmt.Errorf("反馈服务请求失败: status=%d code=%d message=%s", httpResp.StatusCode, envelope.Code, message)
	}

	return envelope.Data, nil
}

func stringWithDefault(value string, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}

func escapeMultipartFilename(fileName string) string {
	fileName = strings.ReplaceAll(fileName, "\\", "\\\\")
	fileName = strings.ReplaceAll(fileName, `"`, "\\\"")
	return fileName
}
