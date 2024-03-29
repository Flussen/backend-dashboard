package jwt

import (
	"fiberproject/config"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func ProtectedByToken() func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing auth token"})
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid auth token format"})
		}

		tokenStr := headerParts[1]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "unexpected signing method")
			}

			return []byte(config.AppConfig.SecretKey), nil
		})

		if err != nil || !token.Valid {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid auth token"})
		}

		return ctx.Next()
	}
}
