package main

import (
	"net/http"
	"os"

	"log/slog"

	_ "github.com/sabbatD/srest-api/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sabbatD/srest-api/internal/config"
	sdb "github.com/sabbatD/srest-api/internal/database"
	"github.com/sabbatD/srest-api/internal/http-server/handlers/admin"
	"github.com/sabbatD/srest-api/internal/http-server/handlers/todo"
	"github.com/sabbatD/srest-api/internal/http-server/handlers/user"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/sabbatD/srest-api/internal/lib/api/access"
	"github.com/sabbatD/srest-api/internal/lib/logger/sl"
)

// @title           sAPI
// @version         v0.3.2
// @description     This is a RESTful API service for EasyDev. It provides various user management functionalities such as user registration, authentication, profile updates, and admin operations.

// @contact.name   s4bb4t
// @contact.email  s4bb4t@yandex.ru

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      easydev.club

// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT Bearer token required for accessing protected routes. Format: "Bearer <token>"

// @schemes https
func main() {
	cfg := config.MustLoad()

	log := sl.SetupLogger(cfg.Env)
	log.Info("Starting sAPI server")
	log.Debug("Debug mode enabled")

	storage, err := sdb.SetupDataBase(cfg.DbString, cfg.Env)
	if err != nil {
		log.Error("Failed to setup database", sl.Err(err))
		os.Exit(1)
	}

	route := chi.NewRouter()
	route.Route("/api/v1", func(router chi.Router) {

		router.Use(middleware.RequestID)
		router.Use(middleware.Logger)
		router.Use(middleware.Recoverer)
		router.Use(middleware.URLFormat)
		router.Use(CORSMiddleware)

		// swagger endpoint
		if cfg.Env != "prod" {
			router.Get("/swagger/*", httpSwagger.Handler(
				httpSwagger.URL("http://51.250.113.72:8082/api/v1/swagger/doc.json"),
			))
		} else {
			router.Get("/swagger/*", httpSwagger.Handler(
				httpSwagger.URL("https://easydev.club/api/v1/swagger/doc.json"),
			))
		}

		// Unknown users handlers
		router.Route("/auth", func(u chi.Router) {

			u.Post("/signup", user.Register(log, storage))
			u.Post("/signin", user.Auth(log, storage))
			u.Post("/refresh", user.Refresh(log, storage))
		})

		// Authenticated user handlers
		// JWTAuthMiddleware used for authenticating users with jwt token from header with prefix "Bearer "
		router.Route("/user", func(u chi.Router) {
			u.Use(access.JWTAuthMiddleware)

			u.Post("/logout", user.Logout(log, storage))

			u.Get("/profile", user.Profile(log, storage))
			u.Put("/profile", user.UpdateUser(log, storage))
			u.Put("/profile/reset-password", user.ChangePassword(log, storage))
		})

		// Authenticated admin handlers
		// JWTAuthMiddleware used for authenticating users with jwt token from header with prefix "Bearer "
		// All of the handlers use AdmCheck.
		router.Route("/admin", func(r chi.Router) {
			r.Use(access.JWTAuthMiddleware)

			r.Get("/users", admin.All(log, storage))

			r.Get("/users/{id}", admin.Profile(log, storage))
			r.Put("/users/{id}", admin.UpdateUser(log, storage))
			r.Delete("/users/{id}", admin.Remove(log, storage))

			r.Post("/users/{id}/block", admin.Block(log, storage))
			r.Post("/users/{id}/unblock", admin.Unblock(log, storage))
			r.Post("/users/{id}/rights", admin.Update(log, storage))
		})

		router.Route("/todos", func(t chi.Router) {
			t.Get("/{id}", todo.Get(log, storage))
			t.Put("/{id}", todo.Update(log, storage))
			t.Delete("/{id}", todo.Delete(log, storage))
		})

		// Todo handlers
		router.Post("/todos", todo.Create(log, storage))
		router.Get("/todos", todo.GetAll(log, storage))
	})

	log.Info("starting server", slog.String("address", cfg.Address))
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      route,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
