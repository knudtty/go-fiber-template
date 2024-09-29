package controllers

import (
	"log"

	ctx "my_project/pkg/context"
	"my_project/pkg/utils"

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

	if err := utils.GetOrCreateOAuthUser(c, oauth2Token, provider); err != nil {
		log.Println("Couldn't get user info: ", err)
		return c.Status(400).Redirect("/login")
	}

	if err := c.ReissueJWT(); err != nil {
		return err
	}

	err = c.DB.SetUserOAuthTokens(oauth2Token.AccessToken, oauth2Token.RefreshToken, oauth2Token.Expiry)
	if err != nil {
		log.Println("Couldn't set refresh and access tokens: ", err)
		return c.Status(400).Redirect("/login")
	}

	return c.Status(fiber.StatusOK).Redirect("/")
}

type loginQueryParams struct {
	Provider string `form:"provider"`
	Username string `form:"username"`
	Password string `form:"password"`
}

// Creates session by username and password or oauth
func CreateSession(c *ctx.WebCtx) error {
	qp := new(loginQueryParams)
	if err := c.QueryParser(qp); err != nil {
		log.Println("Error parsing login query: ", err)
		return err
	}

	if qp.Provider != "" {
		return utils.RedirectToProvider(c.Ctx, qp.Provider)
	} else {
		return nil //TODO: Log in with username and password
	}
}
