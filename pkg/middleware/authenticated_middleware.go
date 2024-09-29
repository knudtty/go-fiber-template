package middleware

import (
	"errors"
	ctx "my_project/pkg/context"

	"github.com/gofiber/fiber/v2"
)

func AuthenticatedUser(c *ctx.WebCtx) error {
	if c.Doer != nil {
		// User is known
		err := c.Next()
		var ferr *fiber.Error
		if errors.As(err, &ferr) && ferr.Code == 404 {
            // TODO: Send to 404 page
			return c.Status(fiber.StatusNotFound).SendString(err.Error())
		}

		return err
	}

	// User not authenticated
	return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized")
}
