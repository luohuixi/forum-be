package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	pb "forum-user/proto"
	"forum/model"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"golang.org/x/net/publicsuffix"
)

const (
	studentLoginActionStart            = "start"
	studentLoginActionCaptcha          = "captcha"
	studentLoginActionSecondAuthSend   = "second_auth_send"
	studentLoginActionSecondAuthVerify = "second_auth_verify"

	studentLoginStatusLoggedIn             = "logged_in"
	studentLoginStatusNeedCaptcha          = "need_captcha"
	studentLoginStatusNeedSecondAuthMethod = "need_second_auth_method"
	studentLoginStatusNeedSecondAuthCode   = "need_second_auth_code"

	studentSecondAuthMethodSMS   = "sms"
	studentSecondAuthMethodEmail = "email"

	ccnuLoginURL           = "https://one.ccnu.edu.cn/"
	ccnuLoginSessionTTL    = 10 * time.Minute
	ccnuRequestTimeout     = 30 * time.Second
	ccnuLoginSessionPrefix = "TeaHouse:ccnu_login:"
)

var (
	ErrStudentCredentialInvalid   = errors.New("student credentials invalid")
	ErrStudentLoginSessionExpired = errors.New("student login session expired")
	ErrStudentLoginBadRequest     = errors.New("student login request invalid")
	ErrStudentLoginUnexpectedPage = errors.New("student login unexpected page")
)

type StudentLoginState struct {
	SessionID                  string
	Status                     string
	Message                    string
	CaptchaImageBase64         string
	AvailableSecondAuthMethods []string
	CurrentSecondAuthMethod    string
}

type StudentLoginResult struct {
	State     StudentLoginState
	StudentID string
	Password  string
}

type ccnuLoginManager struct {
	redis *redis.Client
}

type ccnuLoginSession struct {
	SessionID                  string             `json:"session_id"`
	StudentID                  string             `json:"student_id"`
	Password                   string             `json:"password"`
	Status                     string             `json:"status"`
	Message                    string             `json:"message"`
	CurrentURL                 string             `json:"current_url"`
	CurrentHTML                string             `json:"current_html"`
	CaptchaImageBase64         string             `json:"captcha_image_base64"`
	AvailableSecondAuthMethods []string           `json:"available_second_auth_methods"`
	CurrentSecondAuthMethod    string             `json:"current_second_auth_method"`
	Cookies                    []serializedCookie `json:"cookies"`
}

type serializedCookie struct {
	URL   string `json:"url"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

func HandleStudentLogin(req *pb.StudentLoginRequest) (*StudentLoginResult, error) {
	manager, err := newCCNULoginManager()
	if err != nil {
		return nil, err
	}
	return manager.handle(req)
}

func newCCNULoginManager() (*ccnuLoginManager, error) {
	if model.RedisDB == nil || model.RedisDB.Self == nil {
		return nil, errors.New("redis is not initialized")
	}
	return &ccnuLoginManager{
		redis: model.RedisDB.Self,
	}, nil
}

func newCCNUHTTPClient() *http.Client {
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	return &http.Client{
		Timeout: ccnuRequestTimeout,
		Jar:     jar,
	}
}

func (m *ccnuLoginManager) handle(req *pb.StudentLoginRequest) (*StudentLoginResult, error) {
	action := normalizeStudentLoginAction(req.GetAction())
	switch action {
	case studentLoginActionStart:
		studentID := strings.TrimSpace(req.GetStudentId())
		password := req.GetPassword()
		if studentID == "" || strings.TrimSpace(password) == "" {
			return nil, fmt.Errorf("%w: student_id and password are required", ErrStudentLoginBadRequest)
		}
		return m.start(studentID, password)
	case studentLoginActionCaptcha:
		if strings.TrimSpace(req.GetSessionId()) == "" || strings.TrimSpace(req.GetCaptcha()) == "" {
			return nil, fmt.Errorf("%w: session_id and captcha are required", ErrStudentLoginBadRequest)
		}
		return m.submitCaptcha(req.GetSessionId(), req.GetCaptcha())
	case studentLoginActionSecondAuthSend:
		if strings.TrimSpace(req.GetSessionId()) == "" {
			return nil, fmt.Errorf("%w: session_id is required", ErrStudentLoginBadRequest)
		}
		return m.sendSecondAuthCode(req.GetSessionId(), req.GetSecondAuthMethod())
	case studentLoginActionSecondAuthVerify:
		if strings.TrimSpace(req.GetSessionId()) == "" || strings.TrimSpace(req.GetSecondAuthCode()) == "" {
			return nil, fmt.Errorf("%w: session_id and second_auth_code are required", ErrStudentLoginBadRequest)
		}
		return m.verifySecondAuthCode(req.GetSessionId(), req.GetSecondAuthCode())
	default:
		return nil, fmt.Errorf("%w: unsupported action %q", ErrStudentLoginBadRequest, action)
	}
}

func normalizeStudentLoginAction(action string) string {
	action = strings.TrimSpace(action)
	if action == "" {
		return studentLoginActionStart
	}
	return action
}

func (m *ccnuLoginManager) start(studentID, password string) (*StudentLoginResult, error) {
	client := newCCNUHTTPClient()
	session := &ccnuLoginSession{
		SessionID: uuid.NewString(),
		StudentID: studentID,
		Password:  password,
	}

	page, err := m.fetchCCNUPage(client, http.MethodGet, ccnuLoginURL, nil, "")
	if err != nil {
		return nil, err
	}

	if findCaptchaURL(page) == "" {
		pageAfter, err := m.submitLoginForm(client, session, page, "")
		if err != nil {
			return nil, err
		}
		return m.resultFromLoginPage(client, session, pageAfter, ccnuCaptchaMessage(pageAfter), true)
	}

	if !ccnuAutoOCRAvailable() || ccnuCaptchaAutoRetry() == 0 {
		return m.buildCaptchaChallenge(client, session, page, "请输入验证码后继续登录。")
	}

	for attempt := 0; attempt < ccnuCaptchaAutoRetry(); attempt++ {
		captchaBase64, err := m.fetchCaptchaBase64(client, page)
		if err != nil {
			if attempt == ccnuCaptchaAutoRetry()-1 {
				return m.buildCaptchaChallenge(client, session, page, "请输入验证码后继续登录。")
			}
			page, err = m.fetchCCNUPage(client, http.MethodGet, ccnuLoginURL, nil, "")
			if err != nil {
				return nil, err
			}
			continue
		}

		captchaText, err := m.recognizeCaptcha(captchaBase64)
		if err != nil {
			if attempt == ccnuCaptchaAutoRetry()-1 {
				return m.buildCaptchaChallenge(client, session, page, "请输入验证码后继续登录。")
			}
			page, err = m.fetchCCNUPage(client, http.MethodGet, ccnuLoginURL, nil, "")
			if err != nil {
				return nil, err
			}
			continue
		}

		pageAfter, err := m.submitLoginForm(client, session, page, captchaText)
		if err != nil {
			return nil, err
		}

		if isCCNUCredentialInvalidPage(pageAfter) {
			return nil, ErrStudentCredentialInvalid
		}
		if isCCNULoggedInPage(pageAfter) || isCCNUSecondAuthPage(pageAfter) {
			return m.resultFromLoginPage(client, session, pageAfter, "", true)
		}
		if !isCCNUCaptchaChallengePage(pageAfter) {
			return m.unexpectedLoginPageResult(session, pageAfter)
		}
		if attempt == ccnuCaptchaAutoRetry()-1 {
			return m.buildCaptchaChallenge(client, session, pageAfter, "OCR 未通过，请输入验证码继续登录。")
		}
		page = pageAfter
	}

	return m.buildCaptchaChallenge(client, session, page, "请输入验证码后继续登录。")
}

func (m *ccnuLoginManager) submitCaptcha(sessionID, captcha string) (*StudentLoginResult, error) {
	session, client, page, err := m.loadSessionPage(sessionID)
	if err != nil {
		return nil, err
	}

	pageAfter, err := m.submitLoginForm(client, session, page, strings.TrimSpace(captcha))
	if err != nil {
		return nil, err
	}
	return m.resultFromLoginPage(client, session, pageAfter, ccnuCaptchaMessage(pageAfter), true)
}

func (m *ccnuLoginManager) sendSecondAuthCode(sessionID, method string) (*StudentLoginResult, error) {
	session, client, page, err := m.loadSessionPage(sessionID)
	if err != nil {
		return nil, err
	}

	page, err = m.switchSecondAuthMethodIfNeeded(client, session, page, method)
	if err != nil {
		return nil, err
	}

	form := page.Form
	if form == nil {
		return nil, errors.New("second auth form not found")
	}
	sendName, sendValue, ok := findSecondAuthSendButton(form)
	if !ok || sendName == "" {
		return nil, errors.New("send second auth button not found")
	}

	payload := extractFormValues(form)
	payload.Set(sendName, sendValue)
	for _, input := range form.Inputs {
		if input.Type == "submit" && input.Name != sendName && input.Name != "" {
			payload.Del(input.Name)
		}
	}

	pageAfter, err := m.fetchCCNUPage(client, http.MethodPost, resolveCCNUURL(page.CurrentURL, form.Action), strings.NewReader(payload.Encode()), "application/x-www-form-urlencoded")
	if err != nil {
		return nil, err
	}

	if isCCNULoggedInPage(pageAfter) {
		return m.loggedInResult(session)
	}
	if !isCCNUSecondAuthPage(pageAfter) {
		return nil, errors.New("unexpected second auth response")
	}
	return m.buildSecondAuthCodeChallenge(client, session, pageAfter, ccnuSecondAuthSendMessage(pageAfter, session.CurrentSecondAuthMethod))
}

func (m *ccnuLoginManager) verifySecondAuthCode(sessionID, code string) (*StudentLoginResult, error) {
	session, client, page, err := m.loadSessionPage(sessionID)
	if err != nil {
		return nil, err
	}

	form := page.Form
	if form == nil {
		return nil, errors.New("second auth form not found")
	}
	codeField := findSecondAuthCodeField(form)
	confirmName, confirmValue, ok := findSecondAuthConfirmButton(form)
	if codeField == "" || !ok || confirmName == "" {
		return nil, errors.New("second auth verify fields not found")
	}

	payload := extractFormValues(form)
	payload.Set(codeField, strings.TrimSpace(code))
	payload.Set(confirmName, confirmValue)
	for _, input := range form.Inputs {
		if input.Type == "submit" && input.Name != confirmName && input.Name != "" {
			payload.Del(input.Name)
		}
	}

	pageAfter, err := m.fetchCCNUPage(client, http.MethodPost, resolveCCNUURL(page.CurrentURL, form.Action), strings.NewReader(payload.Encode()), "application/x-www-form-urlencoded")
	if err != nil {
		return nil, err
	}

	if isCCNULoggedInPage(pageAfter) {
		return m.loggedInResult(session)
	}
	if !isCCNUSecondAuthPage(pageAfter) {
		return nil, errors.New("unexpected second auth verify response")
	}
	return m.buildSecondAuthCodeChallenge(client, session, pageAfter, ccnuSecondAuthVerifyMessage(pageAfter))
}

func (m *ccnuLoginManager) resultFromLoginPage(client *http.Client, session *ccnuLoginSession, page *ccnuPageSnapshot, fallbackMessage string, persist bool) (*StudentLoginResult, error) {
	switch {
	case isCCNUCredentialInvalidPage(page):
		_ = m.deleteSession(session.SessionID)
		return nil, ErrStudentCredentialInvalid
	case isCCNULoggedInPage(page):
		return m.loggedInResult(session)
	case isCCNUSecondAuthPage(page):
		message := fallbackMessage
		if strings.TrimSpace(message) == "" {
			message = "已进入二次认证，请选择认证方式。"
		}
		return m.buildSecondAuthMethodChallenge(client, session, page, message, persist)
	default:
		if !isCCNUCaptchaChallengePage(page) {
			return m.unexpectedLoginPageResult(session, page)
		}
		message := fallbackMessage
		if strings.TrimSpace(message) == "" {
			message = ccnuCaptchaMessage(page)
		}
		return m.buildCaptchaChallenge(client, session, page, message)
	}
}

func (m *ccnuLoginManager) loggedInResult(session *ccnuLoginSession) (*StudentLoginResult, error) {
	_ = m.deleteSession(session.SessionID)
	return &StudentLoginResult{
		StudentID: session.StudentID,
		Password:  session.Password,
		State: StudentLoginState{
			Status:  studentLoginStatusLoggedIn,
			Message: "登录成功。",
		},
	}, nil
}

func (m *ccnuLoginManager) buildCaptchaChallenge(client *http.Client, session *ccnuLoginSession, page *ccnuPageSnapshot, message string) (*StudentLoginResult, error) {
	captchaBase64, err := m.fetchCaptchaBase64(client, page)
	if err != nil {
		return nil, err
	}
	session.Status = studentLoginStatusNeedCaptcha
	session.Message = message
	session.CurrentURL = page.CurrentURL
	session.CurrentHTML = page.HTML
	session.CaptchaImageBase64 = captchaBase64
	session.AvailableSecondAuthMethods = nil
	session.CurrentSecondAuthMethod = ""
	session.Cookies = snapshotCCNUCookies(client.Jar, page)
	if err := m.saveSession(session); err != nil {
		return nil, err
	}
	return &StudentLoginResult{
		StudentID: session.StudentID,
		Password:  session.Password,
		State:     session.publicState(),
	}, nil
}

func (m *ccnuLoginManager) buildSecondAuthMethodChallenge(client *http.Client, session *ccnuLoginSession, page *ccnuPageSnapshot, message string, persist bool) (*StudentLoginResult, error) {
	methods, currentMethod, _ := detectSecondAuthMethods(page)
	if len(methods) == 0 {
		methods = []string{studentSecondAuthMethodSMS}
	}
	if currentMethod == "" {
		currentMethod = methods[0]
	}

	session.Status = studentLoginStatusNeedSecondAuthMethod
	session.Message = message
	session.CurrentURL = page.CurrentURL
	session.CurrentHTML = page.HTML
	session.CaptchaImageBase64 = ""
	session.AvailableSecondAuthMethods = methods
	session.CurrentSecondAuthMethod = currentMethod
	session.Cookies = snapshotCCNUCookies(client.Jar, page)
	if persist {
		if err := m.saveSession(session); err != nil {
			return nil, err
		}
	}

	return &StudentLoginResult{
		StudentID: session.StudentID,
		Password:  session.Password,
		State:     session.publicState(),
	}, nil
}

func (m *ccnuLoginManager) buildSecondAuthCodeChallenge(client *http.Client, session *ccnuLoginSession, page *ccnuPageSnapshot, message string) (*StudentLoginResult, error) {
	methods, currentMethod, _ := detectSecondAuthMethods(page)
	if len(methods) == 0 {
		methods = session.AvailableSecondAuthMethods
	}
	if len(methods) == 0 {
		methods = []string{studentSecondAuthMethodSMS}
	}
	if currentMethod == "" {
		currentMethod = normalizeStudentSecondAuthMethod(session.CurrentSecondAuthMethod)
	}
	if currentMethod == "" {
		currentMethod = methods[0]
	}

	session.Status = studentLoginStatusNeedSecondAuthCode
	session.Message = message
	session.CurrentURL = page.CurrentURL
	session.CurrentHTML = page.HTML
	session.CaptchaImageBase64 = ""
	session.AvailableSecondAuthMethods = methods
	session.CurrentSecondAuthMethod = currentMethod
	session.Cookies = snapshotCCNUCookies(client.Jar, page)
	if err := m.saveSession(session); err != nil {
		return nil, err
	}

	return &StudentLoginResult{
		StudentID: session.StudentID,
		Password:  session.Password,
		State:     session.publicState(),
	}, nil
}

func (m *ccnuLoginManager) submitLoginForm(client *http.Client, session *ccnuLoginSession, page *ccnuPageSnapshot, captcha string) (*ccnuPageSnapshot, error) {
	if page == nil || page.Form == nil {
		return nil, errors.New("login form not found")
	}
	payload := extractFormValues(page.Form)
	payload.Set("username", session.StudentID)
	payload.Set("password", session.Password)
	if strings.TrimSpace(captcha) != "" {
		payload.Set("captcha", strings.TrimSpace(captcha))
	}
	return m.fetchCCNUPage(client, http.MethodPost, resolveCCNUURL(page.CurrentURL, page.Form.Action), strings.NewReader(payload.Encode()), "application/x-www-form-urlencoded")
}

func (m *ccnuLoginManager) switchSecondAuthMethodIfNeeded(client *http.Client, session *ccnuLoginSession, page *ccnuPageSnapshot, requestedMethod string) (*ccnuPageSnapshot, error) {
	methods, currentMethod, switchURLs := detectSecondAuthMethods(page)
	if len(methods) == 0 {
		methods = session.AvailableSecondAuthMethods
	}
	requestedMethod = normalizeStudentSecondAuthMethod(requestedMethod)
	if requestedMethod == "" {
		requestedMethod = normalizeStudentSecondAuthMethod(session.CurrentSecondAuthMethod)
	}
	if requestedMethod == "" && len(methods) > 0 {
		requestedMethod = methods[0]
	}
	if requestedMethod == "" {
		return nil, fmt.Errorf("%w: second_auth_method is required", ErrStudentLoginBadRequest)
	}
	if len(methods) > 0 && !containsString(methods, requestedMethod) {
		return nil, fmt.Errorf("%w: second_auth_method %q is unavailable", ErrStudentLoginBadRequest, requestedMethod)
	}
	if currentMethod == "" {
		currentMethod = normalizeStudentSecondAuthMethod(session.CurrentSecondAuthMethod)
	}
	if requestedMethod == currentMethod {
		session.CurrentSecondAuthMethod = requestedMethod
		return page, nil
	}

	switchURL := strings.TrimSpace(switchURLs[requestedMethod])
	if switchURL == "" {
		return nil, fmt.Errorf("%w: second_auth_method %q switch url not found", ErrStudentLoginBadRequest, requestedMethod)
	}
	nextPage, err := m.fetchCCNUPage(client, http.MethodGet, switchURL, nil, "")
	if err != nil {
		return nil, err
	}
	session.CurrentSecondAuthMethod = requestedMethod
	return nextPage, nil
}

func (m *ccnuLoginManager) loadSessionPage(sessionID string) (*ccnuLoginSession, *http.Client, *ccnuPageSnapshot, error) {
	session, err := m.loadSession(sessionID)
	if err != nil {
		return nil, nil, nil, err
	}
	client := newCCNUHTTPClient()
	if err := restoreCCNUCookies(client.Jar, session.Cookies); err != nil {
		return nil, nil, nil, err
	}
	page, err := parseCCNUPage(session.CurrentURL, session.CurrentHTML)
	if err != nil {
		return nil, nil, nil, err
	}
	return session, client, page, nil
}

func (m *ccnuLoginManager) loadSession(sessionID string) (*ccnuLoginSession, error) {
	value, err := m.redis.Get(ccnuLoginSessionKey(sessionID)).Result()
	if errors.Is(err, redis.Nil) {
		return nil, ErrStudentLoginSessionExpired
	}
	if err != nil {
		return nil, err
	}

	var session ccnuLoginSession
	if err := json.Unmarshal([]byte(value), &session); err != nil {
		return nil, err
	}
	if session.SessionID == "" {
		session.SessionID = sessionID
	}
	return &session, nil
}

func (m *ccnuLoginManager) saveSession(session *ccnuLoginSession) error {
	payload, err := json.Marshal(session)
	if err != nil {
		return err
	}
	return m.redis.Set(ccnuLoginSessionKey(session.SessionID), payload, ccnuLoginSessionTTL).Err()
}

func (m *ccnuLoginManager) deleteSession(sessionID string) error {
	if strings.TrimSpace(sessionID) == "" {
		return nil
	}
	return m.redis.Del(ccnuLoginSessionKey(sessionID)).Err()
}

func (m *ccnuLoginManager) fetchCCNUPage(client *http.Client, method, rawURL string, body io.Reader, contentType string) (*ccnuPageSnapshot, error) {
	req, err := http.NewRequest(method, rawURL, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 forum-ccnu-login")
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return parseCCNUPage(resp.Request.URL.String(), string(bodyBytes))
}

func (m *ccnuLoginManager) fetchCaptchaBase64(client *http.Client, page *ccnuPageSnapshot) (string, error) {
	captchaURL := findCaptchaURL(page)
	if captchaURL == "" {
		return "", errors.New("captcha image not found")
	}

	req, err := http.NewRequest(http.MethodGet, captchaURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 forum-ccnu-login")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bodyBytes), nil
}

func (m *ccnuLoginManager) unexpectedLoginPageResult(session *ccnuLoginSession, page *ccnuPageSnapshot) (*StudentLoginResult, error) {
	_ = m.deleteSession(session.SessionID)
	return nil, fmt.Errorf("%w: %s", ErrStudentLoginUnexpectedPage, ccnuUnexpectedLoginMessage(page))
}

func (s *ccnuLoginSession) publicState() StudentLoginState {
	methods := make([]string, len(s.AvailableSecondAuthMethods))
	copy(methods, s.AvailableSecondAuthMethods)
	return StudentLoginState{
		SessionID:                  s.SessionID,
		Status:                     s.Status,
		Message:                    s.Message,
		CaptchaImageBase64:         s.CaptchaImageBase64,
		AvailableSecondAuthMethods: methods,
		CurrentSecondAuthMethod:    s.CurrentSecondAuthMethod,
	}
}

func ccnuLoginSessionKey(sessionID string) string {
	return ccnuLoginSessionPrefix + strings.TrimSpace(sessionID)
}

func snapshotCCNUCookies(jar http.CookieJar, page *ccnuPageSnapshot) []serializedCookie {
	if jar == nil {
		return nil
	}
	targets := ccnuCookieTargetURLs(page)
	cookies := make([]serializedCookie, 0, 8)
	seen := map[string]struct{}{}
	for _, rawURL := range targets {
		parsed, err := url.Parse(rawURL)
		if err != nil {
			continue
		}
		for _, cookie := range jar.Cookies(parsed) {
			key := rawURL + "|" + cookie.Name + "|" + cookie.Value
			if _, exists := seen[key]; exists {
				continue
			}
			seen[key] = struct{}{}
			cookies = append(cookies, serializedCookie{
				URL:   rawURL,
				Name:  cookie.Name,
				Value: cookie.Value,
			})
		}
	}
	return cookies
}

func restoreCCNUCookies(jar http.CookieJar, cookies []serializedCookie) error {
	if jar == nil || len(cookies) == 0 {
		return nil
	}
	grouped := make(map[string][]*http.Cookie)
	for _, cookie := range cookies {
		if strings.TrimSpace(cookie.URL) == "" || strings.TrimSpace(cookie.Name) == "" {
			continue
		}
		grouped[cookie.URL] = append(grouped[cookie.URL], &http.Cookie{
			Name:  cookie.Name,
			Value: cookie.Value,
		})
	}
	for rawURL, items := range grouped {
		parsed, err := url.Parse(rawURL)
		if err != nil {
			return err
		}
		jar.SetCookies(parsed, items)
	}
	return nil
}

func ccnuCookieTargetURLs(page *ccnuPageSnapshot) []string {
	targets := make([]string, 0, 12)
	seen := map[string]struct{}{}
	addTarget := func(rawURL string) {
		trimmed := strings.TrimSpace(rawURL)
		if trimmed == "" {
			return
		}
		parsed, err := url.Parse(trimmed)
		if err != nil {
			return
		}
		if parsed.Scheme != "http" && parsed.Scheme != "https" {
			return
		}
		host := strings.ToLower(parsed.Hostname())
		if host != "one.ccnu.edu.cn" && host != "account.ccnu.edu.cn" {
			return
		}
		if _, exists := seen[trimmed]; exists {
			return
		}
		seen[trimmed] = struct{}{}
		targets = append(targets, trimmed)
	}

	for _, rawURL := range []string{
		"https://one.ccnu.edu.cn/",
		"https://one.ccnu.edu.cn/cas/login",
		"https://account.ccnu.edu.cn/",
		"https://account.ccnu.edu.cn/cas/login",
		"https://account.ccnu.edu.cn/account/secondlogin.jsf",
		"https://account.ccnu.edu.cn/account/secondloginm.jsf",
	} {
		addTarget(rawURL)
	}
	if page == nil {
		return targets
	}

	addTarget(page.CurrentURL)
	if page.Form != nil {
		addTarget(resolveCCNUURL(page.CurrentURL, page.Form.Action))
	}
	addTarget(findCaptchaURL(page))
	for _, link := range page.Links {
		addTarget(resolveCCNUURL(page.CurrentURL, link.Href))
	}
	return targets
}

func containsString(list []string, target string) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}
