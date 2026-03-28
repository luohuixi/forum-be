package auth

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"golang.org/x/net/publicsuffix"
)

func TestParseCCNULoginPage(t *testing.T) {
	htmlContent := mustReadFixture(t, "ccnuNew-index.html")
	page, err := parseCCNUPage(ccnuLoginURL, htmlContent)
	if err != nil {
		t.Fatalf("parseCCNUPage() error = %v", err)
	}
	if page.Form == nil {
		t.Fatal("expected login form")
	}
	payload := extractFormValues(page.Form)
	for _, key := range []string{"lt", "execution", "_eventId"} {
		if payload.Get(key) == "" {
			t.Fatalf("expected form field %q", key)
		}
	}
	captchaURL := findCaptchaURL(page)
	if captchaURL == "" {
		t.Fatal("expected captcha url")
	}
}

func TestDetectSecondAuthSMSPage(t *testing.T) {
	htmlContent := mustReadFixture(t, "sms-post-send", "secondlogin-before-send.html")
	page, err := parseCCNUPage("https://account.ccnu.edu.cn/account/secondlogin.jsf", htmlContent)
	if err != nil {
		t.Fatalf("parseCCNUPage() error = %v", err)
	}

	methods, currentMethod, switchURLs := detectSecondAuthMethods(page)
	if !reflect.DeepEqual(methods, []string{studentSecondAuthMethodSMS, studentSecondAuthMethodEmail}) {
		t.Fatalf("unexpected methods: %#v", methods)
	}
	if currentMethod != studentSecondAuthMethodSMS {
		t.Fatalf("unexpected current method: %q", currentMethod)
	}
	if switchURLs[studentSecondAuthMethodEmail] == "" {
		t.Fatal("expected email switch url")
	}
	if _, _, ok := findSecondAuthSendButton(page.Form); !ok {
		t.Fatal("expected send button")
	}
	if findSecondAuthCodeField(page.Form) == "" {
		t.Fatal("expected second auth code field")
	}
}

func TestDetectSecondAuthEmailPage(t *testing.T) {
	htmlContent := mustReadFixture(t, "email-send", "email-page-before-send.html")
	page, err := parseCCNUPage("https://account.ccnu.edu.cn/account/secondloginm.jsf", htmlContent)
	if err != nil {
		t.Fatalf("parseCCNUPage() error = %v", err)
	}

	methods, currentMethod, switchURLs := detectSecondAuthMethods(page)
	if !reflect.DeepEqual(methods, []string{studentSecondAuthMethodSMS, studentSecondAuthMethodEmail}) {
		t.Fatalf("unexpected methods: %#v", methods)
	}
	if currentMethod != studentSecondAuthMethodEmail {
		t.Fatalf("unexpected current method: %q", currentMethod)
	}
	if switchURLs[studentSecondAuthMethodSMS] == "" {
		t.Fatal("expected sms switch url")
	}
}

func TestSecondAuthVerifyMessage(t *testing.T) {
	htmlContent := mustReadFixture(t, "sms-post-verify", "after-verify.html")
	page, err := parseCCNUPage("https://account.ccnu.edu.cn/account/secondlogin.jsf", htmlContent)
	if err != nil {
		t.Fatalf("parseCCNUPage() error = %v", err)
	}
	if !isCCNUSecondAuthPage(page) {
		t.Fatal("expected second auth page")
	}
	if got := ccnuSecondAuthVerifyMessage(page); got != "验证码错误，请重新输入！" {
		t.Fatalf("unexpected verify message: %q", got)
	}
}

func TestFindSecondAuthCodeFieldPrefersForm1Input(t *testing.T) {
	form := &ccnuHTMLForm{
		Inputs: []ccnuHTMLInput{
			{Name: "search", Type: "text"},
			{Name: "form1:hidcount", ID: "form1:hidcount", Type: "text"},
			{Name: "form1:verifyCode", ID: "form1:verifyCode", Type: "text"},
		},
	}
	if got := findSecondAuthCodeField(form); got != "form1:verifyCode" {
		t.Fatalf("unexpected code field: %q", got)
	}
}

func TestIsCCNUCaptchaChallengePage(t *testing.T) {
	loginHTML := mustReadFixture(t, "step1-login-page.html")
	page, err := parseCCNUPage("https://one.ccnu.edu.cn/cas/login?service=test", loginHTML)
	if err != nil {
		t.Fatalf("parseCCNUPage() error = %v", err)
	}
	if !isCCNUCaptchaChallengePage(page) {
		t.Fatal("expected login page to be treated as captcha challenge page")
	}

	unexpected := &ccnuPageSnapshot{
		CurrentURL: "https://account.ccnu.edu.cn/error",
		Title:      "系统异常",
		BodyText:   "服务异常，请稍后再试",
	}
	if isCCNUCaptchaChallengePage(unexpected) {
		t.Fatal("unexpected page should not be treated as captcha challenge page")
	}
}

func TestSnapshotAndRestoreCCNUCookiesPreservesAccountScope(t *testing.T) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		t.Fatalf("cookiejar.New() error = %v", err)
	}

	secondAuthURL := "https://account.ccnu.edu.cn/account/secondlogin.jsf;jsessionid=ABC123"
	parsedSecondAuthURL, err := url.Parse(secondAuthURL)
	if err != nil {
		t.Fatalf("url.Parse() error = %v", err)
	}
	jar.SetCookies(parsedSecondAuthURL, []*http.Cookie{
		{Name: "JSESSIONID", Value: "abc123"},
	})

	page := &ccnuPageSnapshot{
		CurrentURL: secondAuthURL,
		Form: &ccnuHTMLForm{
			Action: "/account/secondlogin.jsf",
			Method: "POST",
		},
		Links: []ccnuHTMLLink{
			{Href: "secondloginm.jsf"},
		},
	}
	serialized := snapshotCCNUCookies(jar, page)
	if len(serialized) == 0 {
		t.Fatal("expected serialized cookies")
	}

	foundSecondAuthTarget := false
	for _, item := range serialized {
		if strings.Contains(item.URL, "/account/secondlogin") {
			foundSecondAuthTarget = true
			break
		}
	}
	if !foundSecondAuthTarget {
		t.Fatalf("expected second auth target url in serialized cookies: %#v", serialized)
	}

	restoredJar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		t.Fatalf("cookiejar.New() error = %v", err)
	}
	if err := restoreCCNUCookies(restoredJar, serialized); err != nil {
		t.Fatalf("restoreCCNUCookies() error = %v", err)
	}

	rootURL, err := url.Parse("https://account.ccnu.edu.cn/")
	if err != nil {
		t.Fatalf("url.Parse() error = %v", err)
	}
	if cookies := restoredJar.Cookies(rootURL); len(cookies) != 0 {
		t.Fatalf("expected no root-scoped cookies, got %#v", cookies)
	}

	emailURL, err := url.Parse("https://account.ccnu.edu.cn/account/secondloginm.jsf")
	if err != nil {
		t.Fatalf("url.Parse() error = %v", err)
	}
	if cookies := restoredJar.Cookies(emailURL); len(cookies) == 0 {
		t.Fatal("expected restored account-scoped cookie for second auth page")
	}
}

func TestSanitizeCaptchaText(t *testing.T) {
	if got := sanitizeCaptchaText(" aB-12c "); got != "B12c" {
		t.Fatalf("unexpected sanitized captcha: %q", got)
	}
}

func TestLoggedInPageByURL(t *testing.T) {
	page := &ccnuPageSnapshot{
		CurrentURL: "https://one.ccnu.edu.cn/index.html",
		BodyText:   "We're sorry but 网上办事服务大厅 doesn't work properly without JavaScript enabled.",
	}
	if !isCCNULoggedInPage(page) {
		t.Fatal("expected logged in page")
	}
}

func mustReadFixture(t *testing.T, parts ...string) string {
	t.Helper()
	fixturePath := filepath.Join(append([]string{"..", "..", "..", "..", "..", "ccnu-mvp", "probe-artifacts"}, parts...)...)
	content, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatalf("read fixture %s: %v", fixturePath, err)
	}
	return string(content)
}
