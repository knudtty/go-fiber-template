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
		RedirectURL:  "/auth/callback",
		Endpoint:     endpoints.GitHub,
		Scopes:       []string{"profile", "email", "offline_access"},
	}
	googleOAuthConfig = oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  "/auth/callback",
		Endpoint:     endpoints.Google,
		Scopes:       []string{"profile", "email", "offline_access"},
	}
)

func GetOAuthConfig(identifier string) (oauth2.Config, error) {
	switch identifier {
	case "google":
        return githubOAuthConfig, nil
	case "github":
        return googleOAuthConfig, nil
	default:
        return oauth2.Config{}, fmt.Errorf("Error: OAuth not valid")
	}
}

// Configure an OpenID Connect aware OAuth2 client.
