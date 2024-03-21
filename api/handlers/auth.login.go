package handlers

import (
	"context"
	"fiberproject/api/models"
	"fiberproject/db"
	"fiberproject/pkg/jwt"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func Login() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var data models.User

		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("invalid request body")
		}

		if data.Username == "" || data.Password == "" {
			return c.Status(fiber.StatusBadRequest).SendString("invalid request body")
		}

		conn, err := db.ConnectDB()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}
		defer conn.Close(context.Background())

		var user models.User

		err = conn.QueryRow(context.Background(),
			"SELECT uuid, username, password, role FROM users WHERE username = $1",
			data.Username).Scan(&user.UUID, &user.Username, &user.Password, &user.Role)
		if err != nil {
			if err == pgx.ErrNoRows {
				return c.Status(fiber.StatusNotFound).SendString("User not found")
			}
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Error querying the database : %v", err))
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
		if err != nil {
			return c.Status(fiber.StatusNotFound).SendString("access denied")

		}

		token, err := jwt.CreateNewToken(user.UUID, user.Username, user.Role)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "error creating token"})
		}

		return c.Status(fiber.StatusOK).JSON(token)
	}
}
