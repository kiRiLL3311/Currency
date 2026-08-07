package main

import (
	"net/http"
	"strings"

	"github.com/kiRiLL3311/Currency/backend/docs"
	"github.com/kiRiLL3311/Currency/backend/internal/config"
	"github.com/kiRiLL3311/Currency/backend/internal/db"
	"github.com/kiRiLL3311/Currency/backend/internal/handlers"
	"github.com/kiRiLL3311/Currency/backend/internal/middleware"
	"github.com/kiRiLL3311/Currency/backend/internal/repository"
	"github.com/kiRiLL3311/Currency/backend/internal/services"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title FOREX Rates API
// @version 1.0
// @description Currency rates and converter microservice.
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT access token. Example: Bearer <token>
func main() {
	config.LoadEnv()
	configureSwagger()

	database := db.Connect()
	defer database.Close()

	repo := repository.NewRateRepository(database)
	service := services.NewRateService(repo)
	handler := &handlers.RateHandler{
		Service: service,
	}

	services.StartRateSync(service)

	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := database.Ping(); err != nil {
			http.Error(w, "Database is not connected", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Get("/rates/sync", handler.SyncRates)

	r.Group(func(r chi.Router) {
		r.Use(middleware.JWT)

		r.Get("/rates", handler.GetRates)
		r.Get("/convert", handler.Convert)
	})

	http.ListenAndServe(":8080", r)
}

func configureSwagger() {
	docs.SwaggerInfo.Host = config.Get("SWAGGER_HOST")
	if base := strings.TrimSpace(config.Get("SWAGGER_BASE_PATH")); base != "" {
		docs.SwaggerInfo.BasePath = base
	}
	if schemes := strings.TrimSpace(config.Get("SWAGGER_SCHEMES")); schemes != "" {
		docs.SwaggerInfo.Schemes = strings.Split(schemes, ",")
	} else if docs.SwaggerInfo.Host == "" {
		docs.SwaggerInfo.Schemes = []string{"https", "http"}
	}
}
