package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/kiRiLL3311/Currency/auth/internal/config"
	"github.com/kiRiLL3311/Currency/auth/internal/db"
	"github.com/kiRiLL3311/Currency/auth/internal/handlers"
	"github.com/kiRiLL3311/Currency/auth/internal/repository"
	"github.com/kiRiLL3311/Currency/auth/internal/services"
)

func main() {

	config.LoadEnv()

	database := db.Connect()
	defer database.Close()

	repo := repository.NewUserRepository(database)

	service := services.NewAuthService(repo)

	handler := &handlers.AuthHandler{
		Service: service,
	}

	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Auth service OK"))
	})

	// We'll implement these next
	r.Post("/register", handler.Register)
	r.Post("/login", handler.Login)

	http.ListenAndServe(":8080", r)
}
