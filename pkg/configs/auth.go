package configs

import (
	"fmt"

	"golang.org/x/oauth2"
)

func GetOAuthConfig(identifier string) (oauth2.Config, error) {
	switch identifier {
	case "github":
		return GithubOAuthConfig, nil
	case "google":
		return GoogleOAuthConfig, nil
	default:
		return oauth2.Config{}, fmt.Errorf("Error: OAuth not valid")
	}
}
