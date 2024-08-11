package configs

import (
	"fmt"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

var (
	githubOAuthConfig = oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("APP_HOST") + "/auth/redirect",
		Endpoint:     endpoints.GitHub,
		Scopes:       []string{"read:user", "user:email"},
	}
	googleOAuthConfig = oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("APP_HOST") + "/auth/redirect",
		Endpoint:     endpoints.Google,
		Scopes:       []string{"profile", "email", "offline_access"},
	}
)

func GetOAuthConfig(identifier string) (oauth2.Config, error) {
	switch identifier {
	case "github":
		return githubOAuthConfig, nil
	case "google":
		return googleOAuthConfig, nil
	default:
		return oauth2.Config{}, fmt.Errorf("Error: OAuth not valid")
	}
}
