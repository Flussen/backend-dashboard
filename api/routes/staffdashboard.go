package routes

import (
	"fiberproject/api/middlewares"

	"github.com/gofiber/fiber/v2"
)

func staffdashboard(api fiber.Router) {
	api.Group("/admin", middlewares.ProtectedByRole("admin"))
}
