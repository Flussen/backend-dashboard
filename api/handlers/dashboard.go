package handlers

import (
	"fiberproject/api/functions"

	"github.com/gofiber/fiber/v2"
)

func Dashboard() fiber.Handler {
	return func(c *fiber.Ctx) error {

		userRole, ok := c.Locals("userRole").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		uuid, ok := c.Locals("uuid").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}

		user, err := functions.GetInfoByUUID(uuid)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User information not found"})
		}

		switch userRole {
		case "admin":
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"role": "admin", "data": user})
		case "mod":
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"role": "mod", "data": user})
		case "user":
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"role": "user", "data": user})
		default:
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
		}
	}
}
