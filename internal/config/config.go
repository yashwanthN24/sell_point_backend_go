package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	Env         string
	DatabaseUrl string
	JwtKey      string
}

func MustLoad() Config {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		panic("Port is required")
	}

	env := os.Getenv("ENV")
	if env == "" {
		panic("ENV is required")
	}

	dbUrl := os.Getenv("DATABASE_URL")

	if dbUrl == "" {
		panic("Database Url is required")
	}

	JwtKey := os.Getenv("JWT_KEY")

	if JwtKey == "" {
		panic("JwtKey is required")
	}

	return Config{
		Port:        port,
		Env:         env,
		DatabaseUrl: dbUrl,
		JwtKey:      JwtKey,
	}
}
