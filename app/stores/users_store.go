package stores

import (
	"context"
	"my_project/app/models"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UsersStore struct {
	DB *sqlx.DB
}

type UsersWithOAuth struct {
	*models.User
	*models.OAuthAccount
}

func (us *UsersStore) GetUserByProviderId(providerId string) (*models.User, *models.OAuthAccount, error) {
	type userReturn struct {
		*models.User
		error
	}
	type accountReturn struct {
		*models.OAuthAccount
		error
	}
	var wg sync.WaitGroup
	var userChan = make(chan userReturn)
	var accountChan = make(chan accountReturn)

	wg.Add(2)
	go func() {
		defer wg.Done()
		user, err := us.getUserByProviderId(providerId)
		userChan <- userReturn{user, err}
	}()
	go func() {
		defer wg.Done()
		account, err := us.getOAuthAccountByProviderId(providerId)
		accountChan <- accountReturn{account, err}
	}()
	go func() {
		wg.Wait()
		close(userChan)
		close(accountChan)
	}()

	user := <-userChan
	account := <-accountChan

	if user.error != nil {
		return nil, nil, user.error
	}

	if account.error != nil {
		return nil, nil, account.error
	}

	return user.User, account.OAuthAccount, nil
}

func (us *UsersStore) getUserByProviderId(providerId string) (*models.User, error) {
	var user models.User

	err := us.DB.Get(&user, `
        SELECT users.* FROM users
            JOIN oauth_accounts oa ON oa.user_id = users.id
            WHERE oa.provider_user_id = $1;
    `, providerId)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (us *UsersStore) getOAuthAccountByProviderId(providerId string) (*models.OAuthAccount, error) {
	var oauthAccount models.OAuthAccount

	err := us.DB.Get(&oauthAccount, `
        SELECT * FROM oauth_accounts oa
            WHERE oa.provider_user_id = $1;
    `, providerId)
	if err != nil {
		return nil, err
	}

	return &oauthAccount, nil
}

func (us *UsersStore) CreateOAuthUser(providerId, provider, email, role string) (*models.User, *models.OAuthAccount, error) {
	tx, err := us.DB.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, nil, err
	}

	_, err = tx.Exec("INSERT INTO users (email, user_status, user_role) VALUES ($1, 1, $2)", email, role)
	if err != nil {
		return nil, nil, err
	}

	_, err = tx.Exec("INSERT INTO oauth_accounts (user_id, provider, provider_user_id) VALUES ((SELECT id FROM users WHERE email = $1), $2, $3)", email, provider, providerId)
	if err != nil {
		return nil, nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, nil, err
	}

	return us.GetUserByProviderId(providerId)
}

func (us *UsersStore) SetUserOAuthTokens(accessToken, refreshToken string, expiration time.Time) error {
	_, err := us.DB.Exec("UPDATE oauth_accounts SET access_token = $1, refresh_token = $2, expires_at = $3, updated_at = NOW()", accessToken, refreshToken, expiration)
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
