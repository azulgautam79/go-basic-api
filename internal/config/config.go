package configuration

import (
	"os"
)

type Config struct {
	Port         string
	DatabasePath string
	Environment  string
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	databasePath := os.Getenv("DATABASE_PATH")
	if databasePath == "" {
		databasePath = "tasks.db"
	}

	environment := os.Getenv("ENV")
	if environment == "" {
		environment = "development"
	}

	return Config{
		Port:         port,
		DatabasePath: databasePath,
		Environment:  environment,
	}
}
