package auth

import (
	"errors"
	"strings"

	"github.com/spf13/viper"
)

const (
	authBasicPath    = "/auth/api"
	registerPath     = "/signup"
	userInfoPath     = "/user"
	authPath         = "/oauth"
	tokenPath        = "/oauth/token"
	refreshTokenPath = "/oauth/token/refresh"
	clientStorePath  = "/oauth/store"
)

var (
	// muxi-auth-server request url
	RegisterURL     string
	UserInfoURL     string
	OauthURL        string
	OauthTokenURL   string
	OauthRefreshURL string
	ClientStoreURL  string

	clientID     string
	clientSecret string
)

func InitVar() error {
	authHost := strings.TrimSpace(viper.GetString("auth_server.host"))
	if authHost == "" {
		return errors.New("auth_server.host is blank")
	}

	basicURL := withHTTP(authHost) + authBasicPath
	RegisterURL = basicURL + registerPath
	UserInfoURL = basicURL + userInfoPath
	OauthURL = basicURL + authPath
	OauthTokenURL = basicURL + tokenPath
	OauthRefreshURL = basicURL + refreshTokenPath
	ClientStoreURL = basicURL + clientStorePath

	clientID = strings.TrimSpace(viper.GetString("auth_server.client_id"))
	clientSecret = strings.TrimSpace(viper.GetString("auth_server.client_secret"))
	return nil
}

func withHTTP(host string) string {
	if strings.HasPrefix(host, "http://") || strings.HasPrefix(host, "https://") {
		return strings.TrimRight(host, "/")
	}
	return "https://" + strings.TrimRight(host, "/")
}
