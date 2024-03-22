package middlewares

import (
	"fiberproject/config"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func ProtectedByRole(requiredRole string) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing auth token"})
		}

		headerParts := strings.Split(authHeader, " ")

		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or missing auth token"})
		}

		tokenStr := headerParts[1]

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {

			return []byte(config.AppConfig.SecretKey), nil
		})

		if err != nil {
			return err
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if claims["role"] == requiredRole {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "insufficient permissions"})

	}
}
