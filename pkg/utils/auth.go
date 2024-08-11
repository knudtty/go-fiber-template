package utils

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"my_project/app/models"
	"my_project/pkg/configs"
	ctx "my_project/pkg/context"
	"my_project/pkg/repository"
	"my_project/platform/database"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

var pendingAuth = "pending-auth"

type OAuth2State struct {
	State    string `json:"oauthState"`
	Provider string `json:"provider"`
}

func RedirectToProvider(c *fiber.Ctx, provider string) error {
	oauthConfig, err := configs.GetOAuthConfig(provider)
	if err != nil {
		return fmt.Errorf("Provider not found")
	}

	oauthState := OAuth2State{
		State:    uuid.NewString(),
		Provider: provider,
	}
	val, err := json.Marshal(oauthState)
	c.Cookie(&fiber.Cookie{
		Name:     pendingAuth,
		Value:    string(val),
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteLaxMode,
	})

	return c.Redirect(oauthConfig.AuthCodeURL(oauthState.State))
}

func GetOAuthToken(c *ctx.Base) (*oauth2.Token, string, error) {
	var oauthState OAuth2State

	stateQuery := c.Query("state")
	stateCookie := c.Cookies(pendingAuth)
	c.ClearCookie(pendingAuth)

	err := json.Unmarshal([]byte(stateCookie), &oauthState)
	if err != nil {
		return nil, "", fmt.Errorf("Unable to get mid auth oauthState: %v", err)
	}

	if stateQuery == "" || oauthState.State == "" || oauthState.State != stateQuery {
		return nil, "", fmt.Errorf("Oauth2 oauthState from query (%s) does not match oauthState from cookie (%s)", stateQuery, stateCookie)
	}

	oauthConfig, err := configs.GetOAuthConfig(oauthState.Provider)
	if err != nil {
		return nil, "", fmt.Errorf("Provider %v not found", oauthState.Provider)
	}

	token, err := oauthConfig.Exchange(c.UserContext(), c.Query("code"))
	return token, oauthState.Provider, err
}

type googleUser struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
}

func GetOrCreateUser(oauth2Token *oauth2.Token, provider string) (*models.User, error) {
	config, err := configs.GetOAuthConfig(provider)
	if err != nil {
		return nil, err
	}

	client := config.Client(context.Background(), oauth2Token)
	var id, email string
	switch provider {
	case "github":
		id, email, err = getGithubInfo(client)
	case "google":
		id, email, err = getGoogleInfo(client)
	default:
		return nil, fmt.Errorf("Couldn't find provider %v", provider)
	}
	if err != nil {
		return nil, fmt.Errorf("Couldn't get id and email from provider %s: %s ", provider, err)
	}

	db, err := database.GetDbConnection()
	if err != nil {
		return nil, err
	}

	user, err := db.GetUserByProviderId(id)
	if err != nil {
		user, err = db.CreateOAuthUser(id, provider, email, repository.UserRoleName)
		if err != nil {
			return nil, err
		}
	}

	if user.Email != email {
		err = db.UsersStore.UpdateUserEmail(user.ID, email)
		if err != nil {
			log.Printf("Failed to update email for user %v from %v to %v: %v", id, user.Email, email, err)
		}
	}

	return user, nil
}

func getGithubInfo(client *http.Client) (string, string, error) {
	type githubId struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
	}

	type githubEmail struct {
		Email   string `json:"email"`
		Primary bool   `json:"primary"`
	}

	// Get Id
	res, err := client.Get("https://api.github.com/user")
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", "", fmt.Errorf("getGithubInfo: %s", err)
	}

	if res.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("getGithubInfo: Github API request failed %s", res.Status)
	}

	var id githubId
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

func getGoogleInfo(client *http.Client) (string, string, error) {
	res, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("Google API request failed: %s", res.Status)
	}

	var user googleUser
	err = json.Unmarshal(body, &user)
	if err != nil {
		return "", "", err
	}

	return user.Sub, user.Email, nil
}
