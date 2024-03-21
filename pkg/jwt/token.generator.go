package jwt

import (
	"fiberproject/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateNewToken(uuid, username, role string) (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uuid"] = uuid
	claims["username"] = username
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Hour * 720).Unix()

	config, err := config.LoadConfigs("dev")
	if err != nil {
		return "", err
	}

	tokenString, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
