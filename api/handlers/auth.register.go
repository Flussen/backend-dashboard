package handlers

import (
	"context"
	"fiberproject/api/models"
	"fiberproject/db"
	"fiberproject/pkg/util/hashing"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func Register() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var user models.User

		if err := c.BodyParser(&user); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
		}

		if user.Username == "" || user.Email == "" || user.Password == "" {
			return c.Status(fiber.StatusBadRequest).SendString("invalid request body")
		}

		conn, err := db.ConnectDB()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}
		defer conn.Close(context.Background())

		hashedPassword, err := hashing.HashPassword(user.Password)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}

		user.Password = hashedPassword

		err = conn.QueryRow(context.Background(),
			"INSERT INTO users (username, password, email) VALUES ($1, $2, $3) RETURNING uuid",
			user.Username, user.Password, user.Email).Scan(&user.UUID)
		if err != nil {

			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Internal Server Error: Unable to create user : %v", err))
		}

		return c.Status(fiber.StatusCreated).SendString("user created")
	}
}
