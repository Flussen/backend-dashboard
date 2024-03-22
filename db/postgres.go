package db

import (
	"context"
	"fiberproject/config"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func ConnectDB() (*pgx.Conn, error) {

	dsn := fmt.Sprintf("postgresql://%v:%v@%v:%v/%v", config.AppConfig.User, config.AppConfig.Password,
		config.AppConfig.Host, config.AppConfig.Port, config.AppConfig.Database)

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}
	return conn, nil
}
