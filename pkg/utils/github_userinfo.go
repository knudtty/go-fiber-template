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

type githubEmail struct {
	Email   string `json:"email"`
	Primary bool   `json:"primary"`
}

func getGithubInfo(oauth2Token *oauth2.Token) (userInfo, error) {
	var ui userInfo
	var githubInfo struct {
		Id        int    `json:"id"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
		Name      string `json:"name"`
	}

	client := configs.GithubOAuthConfig.Client(context.Background(), oauth2Token)

	// Get Id
	res, err := client.Get("https://api.github.com/user")
	if err != nil || res.StatusCode != http.StatusOK {
		return ui, fmt.Errorf("getGithubInfo: Github API request failed with status %s: %s", res.Status, err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return ui, fmt.Errorf("getGithubInfo: %s", err)
	}

	err = json.Unmarshal(body, &githubInfo)
	if err != nil {
		return ui, fmt.Errorf("getGithubInfo: failed to unmarshal: %s", err)
	}

	ui.id = strconv.Itoa(githubInfo.Id)
	ui.avatarURL = githubInfo.AvatarURL
	ui.name = githubInfo.Name
	if githubInfo.Email != "" {
		// Email was included in user query, return now
		ui.email = githubInfo.Email
		return ui, nil
	}

	// Email was not returned in user query, must use emails endpoint
	res, err = client.Get("https://api.github.com/user/emails")
	body, err = io.ReadAll(res.Body)
	if err != nil {
		return ui, err
	}

	var emails []githubEmail
	var email string

	err = json.Unmarshal(body, &emails)
	if err != nil {
		return ui, fmt.Errorf("getGithubInfo: failed to unmarshal: %s", err)
	}

	for _, e := range emails {
		if e.Primary {
			email = e.Email
			break
		}
	}
	if email == "" {
		return ui, fmt.Errorf("getGithubInfo: No primary email found: %s", err)
	}
	ui.email = email

	return ui, nil
}
