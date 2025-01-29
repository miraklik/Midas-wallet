package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	CONTRACT_ADDRESS string `mapstructure:"CONTRACT_ADDRESS"`
	RPC_URL          string `mapstructure:"RPC_URL"`
	SK               string `mapstructure:"SK"`

	PORT string `mapstructure:"PORT"`

	DB_USER string `mapstructure:"DB_USER"`
	DB_PASS string `mapstructure:"DB_PASS"`
	DB_NAME string `mapstructure:"DB_NAME"`
	DB_HOST string `mapstructure:"DB_HOST"`
	DB_PORT string `mapstructure:"DB_PORT"`
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file")
		return nil, err
	}

	return &Config{
		CONTRACT_ADDRESS: os.Getenv("CONTRACT_ADDRESS"),
		RPC_URL:          os.Getenv("RPC_URL"),
		SK:               os.Getenv("SK"),
		PORT:             os.Getenv("PORT"),
		DB_USER:          os.Getenv("DB_USER"),
		DB_PASS:          os.Getenv("DB_PASS"),
		DB_NAME:          os.Getenv("DB_NAME"),
		DB_HOST:          os.Getenv("DB_HOST"),
		DB_PORT:          os.Getenv("DB_PORT"),
	}, nil
}
