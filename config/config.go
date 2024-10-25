package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	APIPort    string

	CACHEHost     string
	CACHEPort     string
	CACHEPassword string
	CACHEDB       func() int
}

func LoadConfig() (*Config, error) {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file")
	}
	config := &Config{
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		APIPort:       os.Getenv("API_PORT"),
		CACHEHost:     os.Getenv("CACHE_HOST"),
		CACHEPort:     os.Getenv("CACHE_PORT"),
		CACHEPassword: os.Getenv("CACHE_PASSWORD"),
		CACHEDB: func() int {
			db, _ := strconv.Atoi(os.Getenv("CACHE_DB"))
			return db
		},
	}
	// Check if the required environment variables are set
	if config.DBHost == "" || config.DBPort == "" || config.DBUser == "" || config.DBPassword == "" || config.DBName == "" || config.APIPort == "" || config.CACHEHost == "" || config.CACHEPort == "" || config.CACHEPassword == "" {
		return nil, fmt.Errorf("required environment variables are not set")
	}

	return config, nil
}
