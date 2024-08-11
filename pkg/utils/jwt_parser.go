package utils

import (
	"fmt"
	"my_project/app/models"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TokenMetadata struct to describe metadata in JWT.
type TokenMetadata struct {
	models.User
	Expires int64
}

// ExtractTokenMetadata func to extract metadata from JWT.
func ExtractTokenMetadata(c *fiber.Ctx) (*TokenMetadata, error) {
	token, err := verifyToken(c)
	if err != nil {
		return nil, err
	}

	// Setting and checking token and credentials.
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		// User ID.
		userID, err := uuid.Parse(claims["id"].(string))
		if err != nil {
			return nil, err
		}

		// Extract email
		email, ok := claims["email"].(string)
		if !ok {
			return nil, fmt.Errorf("Email not found in jwt claims")
		}

		// Extract role
		role, ok := claims["role"].(string)
		if !ok {
			return nil, fmt.Errorf("Role %v not found", role)
		}
		if err := verifyRole(role); err != nil {
			return nil, fmt.Errorf("Role %v not found", role)
		}

		// Expires time.
		expires := int64(claims["expires"].(float64))

		return &TokenMetadata{
			User: models.User{
				ID:       userID,
				Email:    email,
				UserRole: role,
			},
			Expires: expires,
		}, nil
	}

	return nil, err
}

func verifyToken(c *fiber.Ctx) (*jwt.Token, error) {
	tokenString := c.Cookies("sess-jwt")

	token, err := jwt.Parse(tokenString, jwtKeyFunc)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func jwtKeyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(os.Getenv("JWT_SECRET_KEY")), nil
}
