package context

import (
	"fmt"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
)

type WebCtx struct {
	*Base
}

func NewWebContext(base *Base) *WebCtx {
	return &WebCtx{
		Base: base,
	}
}

func (self *WebCtx) Render(component templ.Component) error {
	self.Set("Content-Type", "text/html")
	return component.Render(self.UserContext(), self.Response().BodyWriter())
}

func (c *WebCtx) ReissueJWT() error {
	t := TokenMetadata{
		User:            c.Doer,
		OAuthAccount:    c.OAuthAccount,
		IsOAuthAccount:  c.OAuthAccount != nil,
	}

	tokens, err := t.GenerateNewTokens()
	if err != nil {
		return fmt.Errorf("Couldn't generate token: %s", err)
	}

	err = c.DB.UpdateUserRefreshToken(t.User.ID, tokens.Refresh)
	if err != nil {
		return fmt.Errorf("Couldn't update refresh token: %s", err)
	}

	c.ClearCookie("sessn-jwt")
	c.Cookie(&fiber.Cookie{
		Name:     "sessn-jwt",
		Value:    tokens.Access,
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteLaxMode,
	})

	return nil
}
