package middlewares

import (
	"fiberproject/pkg/jwt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func ExtractUUID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uuid, err := jwt.GetUserUUIDFromRequest(c)
		if err != nil {
			log.Println(err)
		}

		c.Locals("uuid", uuid)

		return c.Next()
	}
}
