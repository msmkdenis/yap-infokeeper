package config

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Address       string
	DatabaseURI   string
	GRPCServer    string
	TokenName     string
	TokenSecret   string
	TokenExpHours int
}

func New() *Config {
	err := godotenv.Load("infokeeper_server.env")
	if err != nil {
		slog.Error("Error loading .env file, using default values")

		config := &Config{}
		config.Address = "127.0.0.1:7000"
		config.DatabaseURI = "postgres://postgres:postgres@localhost:5432/yap-infokeeper?sslmode=disable"
		config.GRPCServer = ":3300"
		config.TokenName = "token"
		config.TokenSecret = "secret"
		config.TokenExpHours = 24
		return config
	}

	config := &Config{}
	config.Address = os.Getenv("RUN_ADDRESS")
	config.DatabaseURI = os.Getenv("DATABASE_URI")
	config.GRPCServer = os.Getenv("GRPC_SERVER")
	config.TokenName = os.Getenv("TOKEN_NAME")
	expHours, err := strconv.Atoi(os.Getenv("TOKEN_EXP_HOURS"))
	if err != nil {
		slog.Error("Error loading .env file, using default values")
		config.TokenExpHours = 24
	} else {
		config.TokenExpHours = expHours
	}
	config.TokenSecret = os.Getenv("TOKEN_SECRET")
	return config
}
