package main

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/kiRiLL3311/Currency/news/docs"
	"github.com/kiRiLL3311/Currency/news/internal/config"
	"github.com/kiRiLL3311/Currency/news/internal/db"
	"github.com/kiRiLL3311/Currency/news/internal/handlers"
	"github.com/kiRiLL3311/Currency/news/internal/middleware"
	"github.com/kiRiLL3311/Currency/news/internal/repository"
	"github.com/kiRiLL3311/Currency/news/internal/services"
)

// @title FOREX News API
// @version 1.0
// @description Market news microservice (region + keyword pull).
// @host localhost:8082
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

	repo := repository.NewNewsRepository(database)
	service := services.NewNewsService(repo, config.Get("NEWS_API_KEY"))
	handler := &handlers.NewsHandler{Service: service}

	services.StartNewsSync(service)

	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("News service OK"))
	})

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Get("/news/sync", handler.Sync)

	r.Group(func(r chi.Router) {
		r.Use(middleware.JWT)
		r.Get("/news", handler.List)
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
