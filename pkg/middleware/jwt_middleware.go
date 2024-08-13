package middleware

import (
	"context"
	"log"
	"my_project/pkg/utils"
	"os"

	jwtMiddleware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTParser() func(*fiber.Ctx) error {
	// Create config for JWT authentication middleware.
	config := jwtMiddleware.Config{
		SigningKey:     jwtMiddleware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET_KEY"))},
		SuccessHandler: jwtSuccess,
		ErrorHandler:   jwtError,
		TokenLookup:    "cookie:sessn-jwt",
		ContextKey:     "jwt", // used in private routes
	}

	return jwtMiddleware.New(config)
}

func jwtSuccess(c *fiber.Ctx) error {
	token := c.Locals("jwt").(*jwt.Token)
	tokenMetadata, err := utils.ExtractVerifiedTokenMetadata(token)
	if err != nil {
		// We can't parse the jwt,
		log.Println("Unexpected jwt received: ", token, err)
		return c.Next()
	}
	c.SetUserContext(context.WithValue(c.UserContext(), "user_data", tokenMetadata.User))
	if tokenMetadata.IsOAuthAccount {
		c.SetUserContext(context.WithValue(c.UserContext(), "oauth_account_data", tokenMetadata.OAuthAccount))
	}
	return c.Next()
}

func jwtError(c *fiber.Ctx, err error) error {
	// We know nothing about this user
	return c.Next()
}
