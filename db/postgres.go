package db

import (
	"context"
	"fiberproject/config"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func ConnectDB() (*pgx.Conn, error) {

	config, err := config.LoadConfigs("dev")
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("postgresql://%v:%v@%v:%v/%v", config.User, config.Password, config.Host, config.Port, config.Database)

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}
	return conn, nil
}
