package main

import (
	"fmt"
	"os"

	// "github.com/go-chi/chi/v5"
	"github.com/sabbatD/srest-api/internal/config"
	sdb "github.com/sabbatD/srest-api/internal/database"
	"github.com/sabbatD/srest-api/internal/logger/sl"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)

	log := sl.SetupLogger(cfg.Env)
	log.Info("Starting sAPI server")
	log.Debug("Debug mode enabled")

	Storage, err := sdb.SetupDataBase(cfg.DbString)
	if err != nil {
		log.Error("Failed to setup database", sl.Err(err))
		os.Exit(1)
	}
	_ = Storage

	// TODO: init router
	// r := chi.NewRouter()

	// TODO: run server
}