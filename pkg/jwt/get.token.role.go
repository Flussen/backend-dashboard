package jwt

import (
	"fiberproject/config"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func GetUserRoleFromRequest(c *fiber.Ctx) string {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return "guest"
	}

	headerParts := strings.Split(authHeader, " ")

	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "guest"
	}

	tokenStr := headerParts[1]

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		config, err := config.LoadConfigs(config.INIT)
		if err != nil {
			return nil, err
		}
		return []byte(config.SecretKey), nil
	})

	if err != nil {
		return "guest"
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if role, ok := claims["role"].(string); ok {
			return role
		}
	}

	return "guest"
}
