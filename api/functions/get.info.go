package functions

import (
	"context"
	"fiberproject/api/models"
	"fiberproject/db"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func GetInfoByUUID(uuid string) (*models.User, error) {
	conn, err := db.ConnectDB()
	if err != nil {
		return nil, err
	}
	defer conn.Close(context.Background())

	var user models.User
	err = conn.QueryRow(context.Background(),
		"SELECT uuid, username, email, role FROM users WHERE uuid = $1", uuid).
		Scan(&user.UUID, &user.Username, &user.Email, &user.Role)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}
