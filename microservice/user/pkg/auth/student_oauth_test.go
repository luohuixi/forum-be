package auth

import (
	"net/url"
	"testing"
)

func TestBuildStudentOAuthLoginURL(t *testing.T) {
	cfg := StudentOAuthConfig{
		ClientID:         "51f03389-2a18-4941-ba73-c85d08201d42",
		CASLoginURL:      "https://account.ccnu.edu.cn/cas/login",
		OAuthCallbackURL: "https://pass.muxixyz.com/auth/api/oauth/cas/callback",
	}

	loginURL, err := BuildStudentOAuthLoginURL(cfg, "http://localhost:8081/login")
	if err != nil {
		t.Fatalf("BuildStudentOAuthLoginURL() error = %v", err)
	}

	parsedLogin, err := url.Parse(loginURL)
	if err != nil {
		t.Fatalf("parse login url: %v", err)
	}
	if parsedLogin.Scheme+"://"+parsedLogin.Host+parsedLogin.Path != cfg.CASLoginURL {
		t.Fatalf("unexpected cas login url: %s", loginURL)
	}

	serviceURL := parsedLogin.Query().Get("service")
	if serviceURL == "" {
		t.Fatal("service query is blank")
	}

	parsedService, err := url.Parse(serviceURL)
	if err != nil {
		t.Fatalf("parse service url: %v", err)
	}
	if parsedService.Scheme+"://"+parsedService.Host+parsedService.Path != cfg.OAuthCallbackURL {
		t.Fatalf("unexpected service url: %s", serviceURL)
	}
	if got := parsedService.Query().Get("callback_url"); got != "http://localhost:8081/login" {
		t.Fatalf("callback_url = %q", got)
	}
	if got := parsedService.Query().Get("client_id"); got != cfg.ClientID {
		t.Fatalf("client_id = %q", got)
	}
}
