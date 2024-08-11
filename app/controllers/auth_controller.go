package controllers

import (
	"log"

	ctx "my_project/pkg/context"
	"my_project/pkg/utils"
	"my_project/platform/database"

	"github.com/gofiber/fiber/v2"
)

func UserSignUp(c *ctx.ApiCtx) error {
	return nil
}

func AuthRedirectFromProvider(c *ctx.WebCtx) error {
	oauth2Token, provider, err := utils.GetOAuthToken(c.Base)
	if err != nil {
		log.Println("Couldn't get oauth token", err)
		return c.Status(400).Redirect("/login")
	}

	user, err := utils.GetOrCreateUser(oauth2Token, provider)
	if err != nil {
		log.Println("Couldn't get user info: ", err)
		return c.Status(400).Redirect("/login")
	}

	tokens, err := utils.GenerateNewTokens(user)
	if err != nil {
		log.Println("Couldn't generate token: ", err)
		return c.Status(400).Redirect("/login")
	}

	store, err := database.GetDbConnection()
	err = store.UpdateUserRefreshToken(user.ID, tokens.Refresh)
	if err != nil {
		log.Println("Couldn't update refresh token: ", err)
		return c.Status(400).Redirect("/login")
	}

    err = store.SetUserOAuthTokens(oauth2Token.AccessToken, oauth2Token.RefreshToken, oauth2Token.Expiry)
	if err != nil {
		log.Println("Couldn't set refresh and access tokens: ", err)
		return c.Status(400).Redirect("/login")
	}

	c.Cookie(&fiber.Cookie{
		Name:     "sessn-jwt",
		Value:    tokens.Access,
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteLaxMode,
	})

	return c.Status(fiber.StatusOK).Redirect("/")
}

type loginQueryParams struct {
	provider string `form:"provider"`
	username string `form:"username"`
	password string `form:"password"`
}

// Creates session by username and password or oauth
func CreateSession(c *ctx.WebCtx) error {
	qp := new(loginQueryParams)
	if err := c.QueryParser(qp); err != nil {
		log.Println("Error parsing login query: ", err)
		return err
	}

	if qp.provider != "" {
		return utils.RedirectToProvider(c.Ctx, qp.provider)
	} else {
		return nil //TODO: Log in with username and password
	}
}
