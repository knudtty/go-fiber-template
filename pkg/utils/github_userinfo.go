package utils

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"my_project/pkg/configs"

	"github.com/goccy/go-json"
	"golang.org/x/oauth2"
)

type githubId struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type githubEmail struct {
	Email   string `json:"email"`
	Primary bool   `json:"primary"`
}

func getGithubInfo(oauth2Token *oauth2.Token) (string, string, error) {
	var id githubId

	client := configs.GithubOAuthConfig.Client(context.Background(), oauth2Token)

	// Get Id
	res, err := client.Get("https://api.github.com/user")
	if err != nil || res.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("getGithubInfo: Github API request failed with status %s: %s", res.Status, err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", "", fmt.Errorf("getGithubInfo: %s", err)
	}

	err = json.Unmarshal(body, &id)
	if err != nil {
		return "", "", fmt.Errorf("getGithubInfo: failed to unmarshal: %s", err)
	}

	if id.Email != "" {
		// Email was included in user query, return now
		return strconv.Itoa(id.Id), id.Email, nil
	}

	// Email was not returned in user query, must use emails endpoint
	res, err = client.Get("https://api.github.com/user/emails")
	body, err = io.ReadAll(res.Body)
	if err != nil {
		return "", "", err
	}

	var emails []githubEmail
	var email string

	err = json.Unmarshal(body, &emails)
	if err != nil {
		return "", "", fmt.Errorf("getGithubInfo: failed to unmarshal: %s", err)
	}

	for _, e := range emails {
		if e.Primary {
			email = e.Email
			break
		}
	}
	if email == "" {
		return "", "", fmt.Errorf("getGithubInfo: No primary email found: %s", err)
	}

	return strconv.Itoa(id.Id), email, nil
}
