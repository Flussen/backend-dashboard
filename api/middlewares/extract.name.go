package middlewares

import (
	"fiberproject/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

func ExtractName() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, err := jwt.GetUserNameFromRequest(c)
		if err != nil {
			panic("ERROR in extract NAME!")
		}

		c.Locals("userRole", userRole)

		return c.Next()
	}
}
