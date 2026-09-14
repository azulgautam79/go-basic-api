package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port  string
	Env   string
	DBUrl string
}

func MustLoad() Config {
	godotenv.Load()

	//! Port
	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT is required")
	}

	//! Env
	env := os.Getenv("ENV")
	if env == "" {
		panic("ENV is required")
	}

	//! DbUrl
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		panic("DATABASE_URL is required")
	}

	return Config{
		Port:  port,
		Env:   env,
		DBUrl: dbUrl,
	}
}
