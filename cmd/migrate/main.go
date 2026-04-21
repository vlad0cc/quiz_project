package main

import (
	"fmt"
	"log"
	"os"

	"profil-math/internal/bootstrap"
	"profil-math/internal/config"
)

func main() {
	cfg := config.Load()
	command := "up"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	switch command {
	case "up", "seed":
		if err := bootstrap.RunMigrations(cfg.DatabaseURL(), cfg.MigrationsPath); err != nil {
			log.Fatal(err)
		}
	case "down":
		if err := bootstrap.RollbackMigrations(cfg.DatabaseURL(), cfg.MigrationsPath); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal(fmt.Errorf("unsupported migration command: %s", command))
	}
}
