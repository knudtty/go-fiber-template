package stores

import (
	"my_project/app/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UsersStore struct {
	DB *sqlx.DB
}

func (us *UsersStore) GetUserByProviderId(providerId string) (*models.User, error) {
	user := models.User{}
	err := us.DB.Get(&user, `
        SELECT * FROM users
            JOIN oauth_accounts oa ON oa.user_id = users.id
            WHERE oa.provider_user_id = $1
    `, providerId)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (us *UsersStore) CreateOAuthUser(providerId, provider, email, role string) (*models.User, error) {
	user := models.User{}
	err := us.DB.Get(&user, `
        WITH new_user AS (
            INSERT INTO users (email, user_status, user_role)
            VALUES ($1, 1, $2)
            RETURNING id
        )
        INSERT INTO oauth_accounts (user_id, provider, provider_user_id)
        SELECT id, $3, $4 FROM new_user
        RETURNING *;
    `, email, role, provider, providerId)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (us *UsersStore) SetUserOAuthTokens(accessToken, refreshToken string) error {
	_, err := us.DB.Exec("UPDATE oauth_accounts SET access_token = $1, refresh_token = $2", accessToken, refreshToken)
	return err
}

func (us *UsersStore) UpdateUserRefreshToken(userId uuid.UUID, refreshToken string) error {
	_, err := us.DB.Exec("UPDATE users SET refresh_token = $1 WHERE id = $2", refreshToken, userId)
	return err
}

func (us *UsersStore) UpdateUserEmail(userId uuid.UUID, email string) error {
	_, err := us.DB.Exec("UPDATE users SET email = $1 WHERE id = $2", email, userId)
	return err
}
