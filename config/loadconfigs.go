package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const INIT string = "dev" // change to "production" to start with production config

func LoadConfigs(deploy string) (*Config, error) {

	var config Config

	if deploy == "production" {
		err := godotenv.Load("./config/prod.env")
		if err != nil {
			return nil, fmt.Errorf("error loading prod.env file: %v", err)
		}
		config.User = os.Getenv("USER")
		config.Password = os.Getenv("PASSWORD")
		config.Host = os.Getenv("HOST")
		config.Port = os.Getenv("PORT")
		config.Database = os.Getenv("DATABASE")
		config.SecretKey = os.Getenv("SECRETKEY")
		return &config, nil
	}

	err := godotenv.Load("./config/dev.env")
	if err != nil {
		return nil, fmt.Errorf("error loading dev.env file: %v", err)
	}
	config.User = os.Getenv("USER")
	config.Password = os.Getenv("PASSWORD")
	config.Host = os.Getenv("HOST")
	config.Port = os.Getenv("PORT")
	config.Database = os.Getenv("DATABASE")
	config.SecretKey = os.Getenv("SECRETKEY")

	return &config, nil
}
