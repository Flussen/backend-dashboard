package handlers

import (
	"fiberproject/api/functions"

	"github.com/gofiber/fiber/v2"
)

func Dashboard() fiber.Handler {
	return func(c *fiber.Ctx) error {

		userRole, ok := c.Locals("userRole").(string)
		if !ok {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid or not provided user role")
		}

		uuid, ok := c.Locals("uuid").(string)
		if !ok {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid or not provided user UUID")
		}

		user, err := functions.GetInfoByUUID(uuid)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User information not found"})
		}

		switch userRole {
		case "admin":
			return c.JSON(fiber.Map{"role": "admin", "data": user})
		case "mod":
			return c.JSON(fiber.Map{"role": "mod", "data": user})
		case "user":
			return c.JSON(fiber.Map{"role": "user", "data": user})
		default:
			return c.JSON(fiber.Map{"role": "guest", "message": user})
		}
	}
}
