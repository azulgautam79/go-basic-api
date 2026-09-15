package main

import (
	"log"
	"os"
	"strconv"

	"github.com/azulgautam79/go-sqlite-tasks/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: migrate <up | down | force [version]>")
	}

	cfg := config.MustLoad()

	m, err := migrate.New(
		"file://migrations",
		cfg.DBUrl,
	)
	if err != nil {
		log.Fatalf("migration.new: %v", err)
	}
	defer m.Close()

	switch os.Args[1] {
	case "up":
		if err := m.Up(); err != nil {
			log.Fatal(err)
		}
	case "down":
		if err := m.Steps(-1); err != nil {
			log.Fatal(err)
		}
	case "force":
		if len(os.Args) < 3 {
			log.Fatal("usage: migrate force <version>")
		}
		version, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("invalid version number: %v", err)
		}
		if err := m.Force(version); err != nil {
			log.Fatal(err)
		}
		log.Printf("Forced database migration version to %d\n", version)

	default:
		log.Fatalf("unknown command: %s", os.Args[1])
	}
}
