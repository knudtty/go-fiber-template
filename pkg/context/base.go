package context

import (
	"context"
	"my_project/app/models"
	"my_project/app/state"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Base struct {
	*fiber.Ctx
	Doer         *models.User
	OAuthAccount *models.OAuthAccount
	*state.AppState
}

func NewBaseContext(c *fiber.Ctx) *Base {
	b := &Base{
		Ctx: c,
	}
	if tokenData, ok := b.GetContext("token_data").(*TokenMetadata); ok {
		b.Doer = tokenData.User
		b.OAuthAccount = tokenData.OAuthAccount
	}

	return b
}

func (b *Base) SetContext(key string, val any) {
	b.SetUserContext(context.WithValue(b.UserContext(), key, val))
}

func (b *Base) GetContext(key string) any {
	return b.UserContext().Value(key)
}

func (b *Base) ClearCookie(key string) {
	b.Cookie(&fiber.Cookie{
		Name: key,
		// Set expiry date to the past to clear
		Expires: time.Now().Add(-(time.Hour * 2)),
	})
}
