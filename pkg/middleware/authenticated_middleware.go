package middleware

import (
	"my_project/app/state"
	ctx "my_project/pkg/context"

	"github.com/gofiber/fiber/v2"
)

func AuthenticatedUser(c *ctx.WebCtx, _ *state.AppState) error {
	if c.Doer != nil {
		// User is known
		return c.Next()
	}

	// User not authenticated
	return c.Status(fiber.StatusUnauthorized).Redirect("/login")
}
