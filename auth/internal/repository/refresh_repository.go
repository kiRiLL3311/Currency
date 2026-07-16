package repository

import (
	"database/sql"
	"time"
)

type RefreshTokenRepository struct {
	DB *sql.DB
}

func NewRefreshTokenRepository(db *sql.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		DB: db,
	}
}

func (r *RefreshTokenRepository) Create(
	userID int,
	token string,
	expiresAt time.Time,
) error {

	query := `
	INSERT INTO refresh_tokens(user_id, token, expires_at)
	VALUES($1,$2,$3)
	`

	_, err := r.DB.Exec(
		query,
		userID,
		token,
		expiresAt,
	)

	return err
}

func (r *RefreshTokenRepository) GetByToken(
	token string,
) (int, error) {

	var userID int

	query := `
	SELECT user_id
	FROM refresh_tokens
	WHERE token=$1
	AND expires_at > NOW()
	`

	err := r.DB.QueryRow(
		query,
		token,
	).Scan(&userID)

	return userID, err
}

func (r *RefreshTokenRepository) Delete(
	token string,
) error {

	query := `
	DELETE FROM refresh_tokens
	WHERE token=$1
	`

	_, err := r.DB.Exec(
		query,
		token,
	)

	return err
}
