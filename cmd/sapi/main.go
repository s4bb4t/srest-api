package main

import (
	"net/http"
	"os"

	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sabbatD/srest-api/internal/config"
	sdb "github.com/sabbatD/srest-api/internal/database"
	"github.com/sabbatD/srest-api/internal/http-server/handlers/admin"
	"github.com/sabbatD/srest-api/internal/http-server/handlers/user"

	"github.com/sabbatD/srest-api/internal/lib/api/access"
	"github.com/sabbatD/srest-api/internal/lib/logger/sl"
)

func main() {
	cfg := config.MustLoad()

	log := sl.SetupLogger(cfg.Env)
	log.Info("Starting sAPI server")
	log.Debug("Debug mode enabled")

	storage, err := sdb.SetupDataBase(cfg.DbString)
	if err != nil {
		log.Error("Failed to setup database", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/signup", user.Register(log, storage))
	router.Post("/signin", user.Auth(log, storage))

	router.Route("/user", func(u chi.Router) {
		u.Use(access.JWTAuthMiddleware)

		u.Put("/profile", user.UpdateUser(log, storage))

		u.Get("/profile", user.Profile(log, storage))
	})

	router.Route("/admin", func(r chi.Router) {
		r.Use(access.JWTAuthMiddleware)
		// r.Use(access.AdminAuthMiddleware)

		// update user's rights admin & blocked
		r.Post("/users/user={id}/rights", admin.Update(log, storage))

		// update user's rights admin & blocked
		r.Post("/users/user={id}/block", admin.Block(log, storage))
		r.Post("/users/user={id}/unblock", admin.Unblock(log, storage))

		// create a new user
		r.Post("/users/registrate/new", user.Register(log, storage)) //

		// update all user's fields
		r.Put("/users/profile/user={id}", admin.UpdateUser(log, storage)) //

		// get whole user's information
		r.Get("/users/profile/user={id}", admin.Profile(log, storage)) //

		// get array of all users with whole information
		r.Get("/users", admin.GetAll(log, storage))

		// delete user with following username
		r.Delete("/users/user={id}", admin.Remove(log, storage))
	})

	log.Info("starting server", slog.String("address", cfg.Address))
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}
