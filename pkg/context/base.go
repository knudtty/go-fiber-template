package context

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

type User struct {
	Email string
	Role  string
}

type Base struct {
	*fiber.Ctx
	Doer *User
}

func NewBaseContext(c *fiber.Ctx) *Base {
	return &Base{
		Ctx: c,
	}
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
