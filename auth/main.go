package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/kiRiLL3311/Currency/auth/docs"
	"github.com/kiRiLL3311/Currency/auth/internal/config"
	"github.com/kiRiLL3311/Currency/auth/internal/db"
	"github.com/kiRiLL3311/Currency/auth/internal/handlers"
	"github.com/kiRiLL3311/Currency/auth/internal/middleware"
	"github.com/kiRiLL3311/Currency/auth/internal/repository"
	"github.com/kiRiLL3311/Currency/auth/internal/services"
)

// @title FOREX Auth API
// @version 1.0
// @description Authentication and user profile microservice.
// @host localhost:8081
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT access token. Example: Bearer <token>
func main() {

	config.LoadEnv()

	database := db.Connect()
	defer database.Close()

	userRepo := repository.NewUserRepository(database)

	refreshRepo := repository.NewRefreshTokenRepository(database)
	service := services.NewAuthService(
		userRepo,
		refreshRepo,
	)

	handler := &handlers.AuthHandler{
		Service: service,
	}

	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Auth service OK"))
	})

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Post("/register", handler.Register)
	r.Post("/login", handler.Login)
	r.Post("/refresh", handler.Refresh)

	r.Group(func(r chi.Router) {

		r.Use(middleware.JWT)

		r.Get("/me", handler.Me)
		r.Patch("/me", handler.UpdateProfile)

	})

	http.ListenAndServe(":8080", r)
}
