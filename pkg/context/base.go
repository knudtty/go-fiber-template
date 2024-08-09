package context

import (
	"context"

	"github.com/gofiber/fiber/v3"
)

type Base struct {
	fiber.Ctx
}

func NewBaseContext(ctx fiber.Ctx) *Base {
	return &Base{
		Ctx: ctx,
	}
}

func (b *Base) SetContext(key string, val any) {
	b.SetUserContext(context.WithValue(b.UserContext(), key, val))
}

func (b *Base) GetContext(key string) any {
	return b.UserContext().Value(key)
}
