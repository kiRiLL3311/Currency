package services

import (
	"errors"
	"strings"

	"github.com/kiRiLL3311/Currency/auth/internal/models"
	"github.com/kiRiLL3311/Currency/auth/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Repo *repository.UserRepository
}

func NewAuthService(repo *repository.UserRepository) *AuthService {
	return &AuthService{
		Repo: repo,
	}
}

func (s *AuthService) Register(req models.RegisterRequest) error {

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hash),
	}

	err = s.Repo.CreateUser(user)
	if err != nil {

		if strings.Contains(err.Error(), "users_username_key") {
			return errors.New("username already exists")
		}

		if strings.Contains(err.Error(), "users_email_key") {
			return errors.New("email already exists")
		}

		return err
	}

	return nil
}

// Looks up the user by email.

// Reads the stored bcrypt hash from password_hash.

// Compares it with the password the user entered.

// Returns the user if they match.

// Returns the same error for both "email not found" and "wrong password", which is a common security practice because it doesn't reveal which part was incorrect.
func (s *AuthService) Login(req models.LoginRequest) (string, error) {

	user, err := s.Repo.GetUserByEmail(req.Email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	)

	if err != nil {
		return "", errors.New("invalid email or password")
	}

	token, err := GenerateJWT(user.ID, user.Email)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) Me(userID int) (*models.User, error) {
	return s.Repo.GetByID(userID)
}
