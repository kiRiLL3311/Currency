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

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// this lets u kno what usermade req
	// userID := r.Context().Value(middleware.UserIDKey).(int)
	// log.Println("Authenticated user:", userID)

	var req models.LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	token, err := h.Service.Login(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

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
