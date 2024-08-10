package controllers

import (
	"fmt"
	"log"

	ctx "my_project/pkg/context"
	"my_project/pkg/utils"
)

func UserSignUp(c *ctx.ApiCtx) error {
	return nil
}

func AuthRedirectFromProvider(c *ctx.WebCtx) error {
	var providerUrl string
	oauth2Token, err := utils.GetOAuthToken(c.Base, providerUrl)
	if err != nil {
		c.Status(400).Redirect("/login")
	}

	// TODO: Store into JWT
	fmt.Println(oauth2Token)
	return nil
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
