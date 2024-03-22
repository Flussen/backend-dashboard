package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func ConfigInit() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading main .env file: %v", err)
	}

	environment := os.Getenv("APP_ENV")
	if environment == "" {
		environment = "dev"
	}

	var envFile string
	if environment == "production" {
		envFile = "./config/prod.env"
	} else {
		envFile = "./config/dev.env"
	}

	if err := godotenv.Load(envFile); err != nil {
		log.Fatalf("Error loading %s file: %v", envFile, err)
	}

	AppConfig = &Config{
		User:      os.Getenv("USER"),
		Password:  os.Getenv("PASSWORD"),
		Host:      os.Getenv("HOST"),
		Port:      os.Getenv("PORT"),
		Database:  os.Getenv("DATABASE"),
		SecretKey: os.Getenv("SECRETKEY"),
	}

	if AppConfig.SecretKey == "" {
		log.Fatal("SECRETKEY must not be empty")
	}
}
