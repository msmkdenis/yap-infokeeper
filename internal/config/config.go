package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Address     string
	DatabaseURI string
	GRPCServer  string
}

func New() *Config {
	err := godotenv.Load("infokeeper.env")
	if err != nil {
		slog.Error("Error loading .env file, using default values")

		config := &Config{}
		config.Address = "127.0.0.1:7000"
		config.DatabaseURI = "postgres://postgres:postgres@localhost:5432/yap-infokeeper?sslmode=disable"
		config.GRPCServer = ":3300"
		return config
	}

	config := &Config{}
	config.Address = os.Getenv("RUN_ADDRESS")
	config.DatabaseURI = os.Getenv("DATABASE_URI")
	config.GRPCServer = os.Getenv("GRPC_SERVER")
	return config
}
