package routes

import (
	"fiberproject/api/handlers"

	"github.com/gofiber/fiber/v2"
)

func auth(api fiber.Router) {

	user := api.Group("/auth")

	user.Post("", handlers.Register())
	user.Get("", handlers.Login())
}
