package routes

import (
	"fiberproject/api/handlers"

	"github.com/gofiber/fiber/v2"
)

func auth(api fiber.Router) {

	user := api.Group("/auth")

	user.Post("/register", handlers.Register())
	user.Post("/login", handlers.Login())
}
