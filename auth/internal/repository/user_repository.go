package repository

import (
	"database/sql"

	"github.com/kiRiLL3311/Currency/auth/internal/models"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	query := `
	INSERT INTO users(username,email,password_hash)
	VALUES($1,$2,$3)
	`

	_, err := r.DB.Exec(
		query,
		user.Username,
		user.Email,
		user.PasswordHash,
	)

	return err
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, created_at
		FROM users
		WHERE email = $1
	`

	var user models.User

	err := r.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// return user info
func (r *UserRepository) GetByID(id int) (*models.User, error) {
	query := `
	SELECT id, username, email, password_hash, created_at
	FROM users
	WHERE id = $1
	`

	user := &models.User{}

	err := r.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
