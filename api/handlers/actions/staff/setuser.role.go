package staffactions

import (
	"context"
	"fiberproject/db"

	"github.com/gofiber/fiber/v2"
)

func SetRole() fiber.Handler {
	return func(c *fiber.Ctx) error {
		type Receive struct {
			UserID string `json:"userid"`
			Role   string `json:"role"`
		}

		var data Receive

		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
		}

		if data.Role == "" || data.UserID == "" {
			return c.Status(fiber.StatusBadRequest).SendString("cannot be empty parameters")
		}

		conn, err := db.ConnectDB()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}
		defer conn.Close(context.Background())

		commandTag, err := conn.Exec(context.Background(),
			"UPDATE users SET role = $1 WHERE uuid = $2", &data.Role, &data.UserID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error: Unable to update the user")
		}

		if commandTag.RowsAffected() == 0 {
			return c.Status(fiber.StatusNotFound).SendString("User not found or no changes made")

		}

		return c.SendStatus(fiber.StatusOK)
	}
}
