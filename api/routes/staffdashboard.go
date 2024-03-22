package routes

import (
	staffactions "fiberproject/api/handlers/actions/staff"
	"fiberproject/api/middlewares"

	"github.com/gofiber/fiber/v2"
)

func staffdashboard(api fiber.Router) {
	admin := api.Group("/admin", middlewares.ProtectedByRole("admin"))

	actions := admin.Group("/actions")
	actions.Post("/roles", staffactions.SetRole())
}
