package utils

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

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
	State    string `json:"state"`
	Provider string `json:"provider"`
}

func RedirectToProvider(c *fiber.Ctx, provider string) error {
	oauthConfig, err := configs.GetOAuthConfig(provider)
	if err != nil {
		return fmt.Errorf("Provider not found")
	}

	state := OAuth2State{
		State:    uuid.NewString(),
		Provider: provider,
	}
	val, err := json.Marshal(state)
	c.Cookie(&fiber.Cookie{
		Name:     pendingAuth,
		Value:    string(val),
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteLaxMode,
	})

	return c.Redirect(oauthConfig.AuthCodeURL(state.State))
}

func GetOAuthToken(c *ctx.Base) (*oauth2.Token, string, error) {
	stateQuery := c.Query("state")
	stateCookie := c.Cookies(pendingAuth)
	c.ClearCookie(pendingAuth)

	state := OAuth2State{}
	err := json.Unmarshal([]byte(stateCookie), &state)
	if err != nil {
		return nil, "", fmt.Errorf("Unable to get mid auth state: %v", err)
	}

	if stateQuery == "" || stateCookie == "" || stateCookie != stateQuery {
		return nil, "", fmt.Errorf("Oauth2 state from query (%s) does not match state from cookie (%s)", stateQuery, stateCookie)
	}

	oauthConfig, err := configs.GetOAuthConfig(state.Provider)
	if err != nil {
		return nil, "", fmt.Errorf("Provider %v not found", state.Provider)
	}

	token, err := oauthConfig.Exchange(c.UserContext(), c.Query("code"))
	return token, state.Provider, err
}

type githubUser struct {
	Id    string `json:"id"`
	Email string `json:"email"`
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

	db, err := database.GetDbConnection()
	if err != nil {
		return nil, err
	}

	user, err := db.UsersStore.GetUserByProviderId(id)
	if err != nil {
		user, err = db.UsersStore.CreateOAuthUser(id, provider, email, repository.UserRoleName)
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
	res, err := client.Get("https://api.github.com/user")

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("Github API request failed: %s", res.Status)
	}

	var user githubUser
	err = json.Unmarshal(body, &user)
	if err != nil {
		return "", "", err
	}

	return user.Id, user.Email, nil
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
