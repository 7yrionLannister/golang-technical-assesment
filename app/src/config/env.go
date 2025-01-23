package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// [Environment] is a struct that defines the environment variables that the application needs.
type Environment struct {
	DataBaseHost     string `env:"DB_HOST,required"`
	DataBaseEngine   string `env:"DB_ENGINE,required"`
	DataBaseUser     string `env:"DB_USER,required"`
	DataBasePassword string `env:"DB_PASSWORD,required"`
	DataBaseName     string `env:"DB_NAME,required"`
	DataBaseUrl      string
}

// global variable that holds the environment variables that the application needs.
// Use only after calling [LoadEnv].
var Env Environment = Environment{}

// loads the environment variables into [Env].
func LoadEnv() error {
	err := godotenv.Load(".env")
	if err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}
	err = env.Parse(&Env)
	if err != nil {
		return fmt.Errorf("error parsing environment variables: %w", err)
	}
	baseURL := fmt.Sprintf("%s://%s:%s@%s/%s", Env.DataBaseEngine, Env.DataBaseUser, Env.DataBasePassword, Env.DataBaseHost, Env.DataBaseName)
	Env.DataBaseUrl = baseURL
	return nil
}
