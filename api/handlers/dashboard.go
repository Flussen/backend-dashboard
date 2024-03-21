package handlers

import "github.com/gofiber/fiber/v2"

func Dashboard() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("userRole").(string)

		switch userRole {
		case "admin":
			return c.JSON(fiber.Map{"message": "Bienvenido al Dashboard de Admin"})
		case "mod":
			return c.JSON(fiber.Map{"message": "Bienvenido al Dashboard de Moderador"})
		case "user":
			return c.JSON(fiber.Map{"message": "Bienvenido al Dashboard"})
		default:
			return c.JSON(fiber.Map{"error": "No has iniciado sesion!"})
		}

	}
}
