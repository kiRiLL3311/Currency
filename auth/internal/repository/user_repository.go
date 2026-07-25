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
	region := user.Region
	if region == "" {
		region = "US"
	}

	query := `
	INSERT INTO users(username,email,password_hash,region)
	VALUES($1,$2,$3,$4)
	`

	_, err := r.DB.Exec(
		query,
		user.Username,
		user.Email,
		user.PasswordHash,
		region,
	)

	return err
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, region, created_at
		FROM users
		WHERE email = $1
	`

	var user models.User

	err := r.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Region,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByID(id int) (*models.User, error) {
	query := `
	SELECT id, username, email, password_hash, region, created_at
	FROM users
	WHERE id = $1
	`

	user := &models.User{}

	err := r.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Region,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) UpdateRegion(id int, region string) error {
	_, err := r.DB.Exec(`
		UPDATE users
		SET region = $1
		WHERE id = $2
	`, region, id)
	return err
}
