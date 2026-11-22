package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                string
	Env                 string
	DatabaseUrl         string
	JwtKey              string
	StorageAccountID    string
	StorageAccessKey    string
	StorageAccessSecret string
	StorageBucket       string
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

	storageAccountID := os.Getenv("STORAGE_AccountID")

	if storageAccountID == "" {
		panic("Object Storage AccountID for Cloudflare R2 is required")
	}

	storageAccessKey := os.Getenv("STORAGE_AccessKey")

	if storageAccessKey == "" {
		panic("Object Storage AccessKey is required")
	}

	storageAccessSecret := os.Getenv("STORAGE_AccessSecret")

	if storageAccessSecret == "" {
		panic("Object Storage AccessSecret is required")
	}

	storageBucket := os.Getenv("STORAGE_BUCKET")

	if storageBucket == "" {
		panic("Object Storage Bucket is required")
	}

	return Config{
		Port:                port,
		Env:                 env,
		DatabaseUrl:         dbUrl,
		JwtKey:              JwtKey,
		StorageAccountID:    storageAccountID,
		StorageAccessKey:    storageAccessKey,
		StorageAccessSecret: storageAccessSecret,
		StorageBucket:       storageBucket,
	}
}
