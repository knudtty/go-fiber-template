package utils

import (
	"fmt"

	"my_project/pkg/configs"
	ctx "my_project/pkg/context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

var pendingAuth = "pending-auth"

func GetOAuthToken(c *ctx.Base, providerUrl string) (*oauth2.Token, error) {
	stateQuery := c.Query("state")
	stateCookie := c.Cookies(pendingAuth)

	c.ClearCookie(pendingAuth)
	if stateQuery == "" || stateCookie == "" || stateCookie != stateQuery {
		return nil, fmt.Errorf("Oauth2 state from query (%s) does not match state from cookie (%s)", stateQuery, stateCookie)
	}

    oauthConfig, err := configs.GetOAuthConfig(providerUrl)
	if err != nil {
		return nil, fmt.Errorf("Provider not found")
	}

	return oauthConfig.Exchange(c.UserContext(), c.Query("code"))
}

func RedirectToProvider(c *fiber.Ctx, providerUrl string) error {
    oauthConfig, err := configs.GetOAuthConfig(providerUrl)
	if err != nil {
		return fmt.Errorf("Provider not found")
	}

	state := uuid.NewString()
	c.Cookie(&fiber.Cookie{
		Name:     pendingAuth,
		Value:    state,
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteLaxMode,
	})

	return c.Redirect(oauthConfig.AuthCodeURL(state))
}
