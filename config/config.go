package config

import (
	l "geotrack_api/config/logger"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func LoadConfig() Config {

	rootPath, err := filepath.Abs("D:/Estudos Go/API/geotrack_api")
	if err != nil {
		l.Logger.Fatal("Error determining root path: %v", zap.Error(err))
	}

	envPath := filepath.Join(rootPath, ".env")

	err = godotenv.Load(envPath)
	if err != nil {
		l.Logger.Fatal("Error loading .env file: %v", zap.Error(err))
	}

	return Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
	}
}
