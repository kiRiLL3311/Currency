package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kiRiLL3311/Currency/auth/internal/middleware"
	"github.com/kiRiLL3311/Currency/auth/internal/models"
	"github.com/kiRiLL3311/Currency/auth/internal/services"
)

type AuthHandler struct {
	Service *services.AuthService
}

// Register godoc
// @Summary Register a user
// @Description Creates a new account. Does not return tokens — sign in afterwards.
// @Tags Auth
// @Accept json
// @Produce plain
// @Param body body models.RegisterRequest true "Registration payload"
// @Success 201 {string} string "User created"
// @Failure 400 {string} string "Invalid JSON"
// @Failure 409 {string} string "Username or email already exists"
// @Failure 500 {string} string "Internal Server Error"
// @Router /register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {

	var req models.RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = h.Service.Register(req)
	if err != nil {

		switch err.Error() {

		case "username already exists":
			http.Error(w, err.Error(), http.StatusConflict)
			return

		case "email already exists":
			http.Error(w, err.Error(), http.StatusConflict)
			return

		default:
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created"))
}

// Login godoc
// @Summary Login
// @Description Returns access and refresh tokens for an existing user.
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body models.LoginRequest true "Login payload"
// @Success 200 {object} models.AuthResponse
// @Failure 400 {string} string "Invalid JSON"
// @Failure 401 {string} string "Unauthorized"
// @Router /login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	tokens, err := h.Service.Login(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(tokens)
}

// Refresh godoc
// @Summary Refresh tokens
// @Description Rotates refresh token and returns a new access + refresh pair.
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body models.RefreshRequest true "Refresh payload"
// @Success 200 {object} models.AuthResponse
// @Failure 400 {string} string "Invalid JSON"
// @Failure 401 {string} string "Unauthorized"
// @Router /refresh [post]
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	tokens, err := h.Service.Refresh(req.RefreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokens)
}

// Me godoc
// @Summary Current user
// @Description Returns the authenticated user profile.
// @Tags Profile
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /me [get]
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {

	userID, ok := r.Context().Value(middleware.UserIDKey).(int)

	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.Service.Me(userID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.PasswordHash = ""

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(user)
}

// UpdateProfile godoc
// @Summary Update profile region
// @Description Updates the user's news region (US, EU, GB, JP, AU).
// @Tags Profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body models.UpdateProfileRequest true "Profile update"
// @Success 200 {object} models.User
// @Failure 400 {string} string "Invalid JSON or region"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /me [patch]
func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	user, err := h.Service.UpdateProfile(userID, req)
	if err != nil {
		if err.Error() == "invalid region" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.PasswordHash = ""
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
