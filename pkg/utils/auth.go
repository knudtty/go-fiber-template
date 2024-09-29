package utils

import (
	"fmt"
	"log"

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

type userInfo struct {
	id        string
	email     string
	avatarURL string
	name      string
}

func GetOrCreateOAuthUser(c *ctx.WebCtx, oauth2Token *oauth2.Token, provider string) error {
	var err error
	var ui userInfo

	switch provider {
	case "github":
		ui, err = getGithubInfo(oauth2Token)
	case "google":
		ui, err = getGoogleInfo(oauth2Token)
	default:
		return fmt.Errorf("Couldn't find provider %v", provider)
	}
	if err != nil {
		return fmt.Errorf("Couldn't get id and email from provider %s: %s ", provider, err)
	}

	db, err := database.GetDbConnection()
	if err != nil {
		return err
	}

	user, oauthAccount, err := db.GetUserByProviderId(ui.id)
	if err != nil {
		user, oauthAccount, err = db.CreateOAuthUser(ui.id, provider, ui.name, ui.email, repository.UserRoleName, ui.avatarURL)
		if err != nil {
			return err
		}
	}

	// TODO: Update for profile pic
	if user.Email != ui.email {
		err = db.UpdateUserEmail(user.ID, ui.email)
		if err != nil {
			log.Printf("Failed to update email for user %v from %v to %v: %v", ui.id, user.Email, ui.email, err)
		}
	}
	if user.AvatarURL != ui.avatarURL {
		err = db.UpdateUserAvatar(user.ID, ui.avatarURL)
		if err != nil {
			log.Printf("Failed to update email for user %v from %v to %v: %v", ui.id, user.Email, ui.email, err)
		}
	}

	c.Doer = user
	c.OAuthAccount = oauthAccount
	return nil
}
