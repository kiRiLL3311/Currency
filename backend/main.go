// package main

// import (
// 	"net/http"

// 	"github.com/kiRiLL3311/Currency/backend/internal/config"
// 	"github.com/kiRiLL3311/Currency/backend/internal/db"
// 	"github.com/kiRiLL3311/Currency/backend/internal/handlers"
// 	"github.com/kiRiLL3311/Currency/backend/internal/middleware"
// 	"github.com/kiRiLL3311/Currency/backend/internal/repository"
// 	"github.com/kiRiLL3311/Currency/backend/internal/services"

// 	_ "github.com/kiRiLL3311/Currency/backend/docs"

// 	chi "github.com/go-chi/chi/v5"
// 	httpSwagger "github.com/swaggo/http-swagger"
// 	// replace "backend" with your module name
// )

// // @title Currency API
// // @version 1.0
// // @description Currency conversion service.
// // @host localhost:8080
// // @BasePath /
// func main() {
// 	//loading env
// 	config.LoadEnv()

// 	database := db.Connect()
// 	defer database.Close()
// 	//rates route
// 	repo := repository.NewRateRepository(database)
// 	service := services.NewRateService(repo)
// 	handler := &handlers.RateHandler{
// 		Service: service,
// 	}
// 	//start sheduler
// 	services.StartRateSync(service)
// 	r := chi.NewRouter()

// 	// Public routes
// 	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
// 		if err := database.Ping(); err != nil {
// 			http.Error(w, "Database is not connected", http.StatusInternalServerError)
// 			return
// 		}

// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte("OK"))
// 	})

// 	r.Get("/swagger/*", httpSwagger.WrapHandler)

// 	// Protected routes
// 	r.Group(func(r chi.Router) {
// 		r.Use(middleware.JWT)

// 		//r.Get("/rates", handler.GetRates)
// 		//r.Get("/convert", handler.Convert)
// 		r.Get("/rates/sync", handler.SyncRates)
// 	})

// 	r.Group(func(r chi.Router) {

// 		r.Use(middleware.JWT)

//			r.Get("/rates", handler.GetRates)
//			r.Get("/convert", handler.Convert)
//		})
//		http.ListenAndServe(":8080", r)
//	}
package main

import (
	"net/http"

	"github.com/kiRiLL3311/Currency/backend/internal/config"
	"github.com/kiRiLL3311/Currency/backend/internal/db"
	"github.com/kiRiLL3311/Currency/backend/internal/handlers"
	"github.com/kiRiLL3311/Currency/backend/internal/middleware"
	"github.com/kiRiLL3311/Currency/backend/internal/repository"
	"github.com/kiRiLL3311/Currency/backend/internal/services"

	_ "github.com/kiRiLL3311/Currency/backend/docs"

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
	// Load environment variables
	config.LoadEnv()

	// Connect database
	database := db.Connect()
	defer database.Close()

	// Dependency Injection
	repo := repository.NewRateRepository(database)
	service := services.NewRateService(repo)
	handler := &handlers.RateHandler{
		Service: service,
	}

	// Start scheduler
	services.StartRateSync(service)

	// Router
	r := chi.NewRouter()

	// -------------------------
	// Public routes
	// -------------------------

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if err := database.Ping(); err != nil {
			http.Error(w, "Database is not connected", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	// Scheduler endpoint (public)
	r.Get("/rates/sync", handler.SyncRates)

	// -------------------------
	// Protected routes
	// -------------------------

	r.Group(func(r chi.Router) {
		r.Use(middleware.JWT)

		r.Get("/rates", handler.GetRates)
		r.Get("/convert", handler.Convert)
	})

	http.ListenAndServe(":8080", r)
}
