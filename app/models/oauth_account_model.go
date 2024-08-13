package models

import (
	"time"

	"github.com/google/uuid"
)

type OAuthAccount struct {
	ID             int       `db:"id" json:"id" validate:"required,int,unique"`
	UserID         uuid.UUID `db:"user_id" json:"user_id" validate:"required,uuid,unique"`
	Provider       string    `db:"provider" json:"provider" validate:"required"`
	ProviderUserID string    `db:"provider_user_id" json:"provider_user_id" validate:"required"`
	AccessToken    string    `db:"refresh_token" json:"refresh_token"`
	RefreshToken   string    `db:"access_token" json:"access_token"`
	ExpiresAt      time.Time `db:"expires_at" json:"expires_at"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}
