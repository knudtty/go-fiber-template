package utils

import (
	"errors"
	"fmt"
	"my_project/app/models"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// TokenMetadata struct to describe metadata in JWT.
type TokenMetadata struct {
	User           models.User
	OAuthAccount   models.OAuthAccount
	IsOAuthAccount bool
	Expires        int64
}

// ExtractVerifiedTokenMetadata func to extract metadata from JWT.
func ExtractVerifiedTokenMetadata(token *jwt.Token) (*TokenMetadata, error) {
	// Setting and checking token and credentials.
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		// User ID.
		userID, err := uuid.Parse(claims["user_id"].(string))
		if err != nil {
			return nil, err
		}

		// Extract email
		email, ok := claims["user_email"].(string)
		if !ok {
			return nil, fmt.Errorf("Email not found in jwt claims")
		}

		// Extract role
		role, ok := claims["user_role"].(string)
		if !ok {
			return nil, fmt.Errorf("Role %v not found", role)
		} else if err := verifyRole(role); err != nil {
			return nil, fmt.Errorf("Role %v not found", role)
		}

		// JWT Expires time.
		expires := int64(claims["expires"].(float64))

		accountType, ok := claims["account_type"]
		if !ok {
			return nil, fmt.Errorf("Unknown account type")
		}
		var isOauthAccount bool
		switch accountType {
		case "oauth2":
			isOauthAccount = true
		case "password":
			isOauthAccount = false
		default:
			return nil, fmt.Errorf("Unknown account type %s", accountType)
		}

		var oauthAccount models.OAuthAccount
		if isOauthAccount {
			oauthAccount.UserID = userID

			// OAuthAccount id
			oauthID, ok := claims["oauth_id"].(float64)
			if !ok {
				return nil, fmt.Errorf("OAuth ID not found")
			}
			oauthAccount.ID = int(oauthID)

			oauthAccount.Provider, ok = claims["oauth_provider"].(string)
			if !ok || !verifyOAuthProvider(oauthAccount.Provider) {
				return nil, fmt.Errorf("Unknown OAuth provider: %s", oauthAccount.Provider)
			}

			oauthAccount.ProviderUserID, ok = claims["oauth_provider_user_id"].(string)
			if !ok {
				return nil, fmt.Errorf("Provider user id not available")
			}
		}

		return &TokenMetadata{
			User: models.User{
				ID:       userID,
				Email:    email,
				UserRole: role,
			},
			IsOAuthAccount: isOauthAccount,
			OAuthAccount:   oauthAccount,
			Expires:        expires,
		}, nil
	}

	return nil, errors.New("Cannot extract metadata from invalid jwt")
}

// Not currently using, here if needed
func verifyToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, jwtKeyFunc)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func jwtKeyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(os.Getenv("JWT_SECRET_KEY")), nil
}
