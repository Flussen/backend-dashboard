package middleware

import (
	"fiberproject/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

func ExtractName() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := jwt.GetUserRoleFromRequest(c)

		c.Locals("userRole", userRole)

		return c.Next()
	}
}
