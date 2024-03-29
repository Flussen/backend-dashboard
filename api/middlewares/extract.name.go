package middlewares

import (
	"fiberproject/pkg/jwt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func ExtractName() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, err := jwt.GetUserNameFromRequest(c)
		if err != nil {
			log.Println(err)
		}

		c.Locals("username", userRole)

		return c.Next()
	}
}
