package routes

import (
	"fiberproject/api/handlers"
	"fiberproject/api/middlewares"

	"github.com/gofiber/fiber/v2"
)

func dashboard(api fiber.Router) {

	dashboard := api.Group("/dashboard", middlewares.ExtractRole(), middlewares.ExtractUUID())

	dashboard.Get("", handlers.Dashboard())

	profile := dashboard.Group("/profile")
	profile.Get("/:id", handlers.GetProfile())
	profile.Put("/:id", handlers.UpdateProfile())
	profile.Put("/:id/pwd", handlers.ChangePassword())
	profile.Delete("/:id", handlers.DeleteAccount())
}
