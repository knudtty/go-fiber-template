package configs

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"

	"github.com/coreos/go-oidc"
)

var GoogleOAuthConfig = oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	RedirectURL:  os.Getenv("APP_HOST") + "/auth/redirect",
	Endpoint:     endpoints.Google,
	Scopes:       []string{oidc.ScopeOpenID, "email", "profile"},
}
