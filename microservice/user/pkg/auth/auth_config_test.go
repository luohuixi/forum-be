package auth

import (
	"testing"

	"github.com/spf13/viper"
)

func TestInitVarRequiresAuthHost(t *testing.T) {
	viper.Reset()
	if err := InitVar(); err == nil {
		t.Fatal("expected error when auth_server.host is blank")
	}
}

func TestInitVarBuildsPassURLs(t *testing.T) {
	viper.Reset()
	viper.Set("auth_server.host", "pass.muxixyz.com")
	viper.Set("auth_server.client_id", "client-a")
	viper.Set("auth_server.client_secret", "secret-a")

	if err := InitVar(); err != nil {
		t.Fatalf("InitVar() error = %v", err)
	}

	if OauthTokenURL != "https://pass.muxixyz.com/auth/api/oauth/token" {
		t.Fatalf("unexpected token url: %s", OauthTokenURL)
	}
	if UserInfoURL != "https://pass.muxixyz.com/auth/api/user" {
		t.Fatalf("unexpected user info url: %s", UserInfoURL)
	}
}
