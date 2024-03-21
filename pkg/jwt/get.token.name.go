package jwt

import (
	"fiberproject/config"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func GetUserNameFromRequest(c *fiber.Ctx) (string, error) {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("error in header Authorization: is clean")
	}

	headerParts := strings.Split(authHeader, " ")

	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", fmt.Errorf("error in header Authorization: need to be 2 arguments")
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
		return "", fmt.Errorf("error in header Authorization: error parsing token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if username, ok := claims["username"].(string); ok {
			return username, nil
		}
	}

	return "", fmt.Errorf("not posible")
}
