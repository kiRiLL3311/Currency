package services

import (
	"errors"
	"strings"

	"github.com/kiRiLL3311/Currency/auth/internal/models"
	"github.com/kiRiLL3311/Currency/auth/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	Repo        *repository.UserRepository
	RefreshRepo *repository.RefreshTokenRepository
}

//	func NewAuthService(repo *repository.UserRepository) *AuthService {
//		return &AuthService{
//			Repo: repo,
//		}
//	}
func NewAuthService(
	repo *repository.UserRepository,
	refreshRepo *repository.RefreshTokenRepository,
) *AuthService {

	return &AuthService{
		Repo:        repo,
		RefreshRepo: refreshRepo,
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
// func (s *AuthService) Login(req models.LoginRequest) (*models.AuthResponse, error) {

// 	user, err := s.Repo.GetUserByEmail(req.Email)
// 	if err != nil {
// 		return "", errors.New("invalid email or password")
// 	}

// 	err = bcrypt.CompareHashAndPassword(
// 		[]byte(user.PasswordHash),
// 		[]byte(req.Password),
// 	)

// 	if err != nil {
// 		return "", errors.New("invalid email or password")
// 	}

// 	token, err := GenerateJWT(user.ID, user.Email)
// 	if err != nil {
// 		return "", err
// 	}

//		return token, nil
//	}
func (s *AuthService) Login(req models.LoginRequest) (*models.AuthResponse, error) {

	user, err := s.Repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	)

	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate short-lived access token
	accessToken, err := GenerateJWT(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	// Store refresh token in database
	err = s.RefreshRepo.Create(
		user.ID,
		refreshToken,
		RefreshExpiry(),
	)

	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// Refresh exchanges a valid refresh token for a new access + refresh token pair.
// The old refresh token is deleted (rotation) so it cannot be reused.
func (s *AuthService) Refresh(refreshToken string) (*models.AuthResponse, error) {
	if strings.TrimSpace(refreshToken) == "" {
		return nil, errors.New("refresh token required")
	}

	userID, err := s.RefreshRepo.GetByToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	user, err := s.Repo.GetByID(userID)
	if err != nil {
		return nil, errors.New("invalid or expired refresh token")
	}

	// Rotate: invalidate the presented refresh token before issuing a new one.
	if err := s.RefreshRepo.Delete(refreshToken); err != nil {
		return nil, err
	}

	accessToken, err := GenerateJWT(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	if err := s.RefreshRepo.Create(user.ID, newRefreshToken, RefreshExpiry()); err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *AuthService) Me(userID int) (*models.User, error) {
	return s.Repo.GetByID(userID)
}

var allowedRegions = map[string]bool{
	"US": true,
	"EU": true,
	"GB": true,
	"JP": true,
	"AU": true,
}

func (s *AuthService) UpdateProfile(userID int, req models.UpdateProfileRequest) (*models.User, error) {
	region := strings.TrimSpace(strings.ToUpper(req.Region))
	if region == "" || !allowedRegions[region] {
		return nil, errors.New("invalid region")
	}

	if err := s.Repo.UpdateRegion(userID, region); err != nil {
		return nil, err
	}

	return s.Repo.GetByID(userID)
}
