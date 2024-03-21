package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type Config struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
}

func ConnectToDB() (*pgx.Conn, error) {
	dsn, err := setDSN()
	if err != nil {
		return nil, err
	}

	fmt.Println(dsn)

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}
	return conn, nil
}

func setDSN() (string, error) {
	config, err := envFile()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("postgres://%v:%v@%v:%v/%v", config.User, config.Password, config.Host, config.Port, config.Database), nil
}

func envFile() (*Config, error) {

	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}
	config := Config{
		User:     os.Getenv("USER"),
		Password: os.Getenv("PASSWORD"),
		Host:     os.Getenv("HOST"),
		Port:     os.Getenv("PORT"),
		Database: os.Getenv("DATABASE"),
	}
	return &config, nil
}
