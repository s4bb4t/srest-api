package main

import (
	"fmt"

	"github.com/sabbatD/srest-api/internal/config"
)

func main() {
	// TODO: init config
	cfg := config.MustLoad()

	fmt.Println(cfg)
	// TODO: init logger

	// TODO: init storage

	// TODO: init router

	// TODO: run server
}
