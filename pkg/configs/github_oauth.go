package configs

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

var GithubOAuthConfig = oauth2.Config{
	ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
	ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
	RedirectURL:  os.Getenv("APP_HOST") + "/auth/redirect",
	Endpoint:     endpoints.GitHub,
	Scopes:       []string{"read:user", "user:email"},
}
