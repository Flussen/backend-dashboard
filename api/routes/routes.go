package routes

import (
	"fiberproject/api/handlers"
	"fiberproject/api/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	// userGroup := app.Group("/users")
	// userGroup.Get("", handlers.GetUsers())
	// userGroup.Get("/:id", handlers.GetUser())
	// userGroup.Post("", handlers.CreateUser())
	// userGroup.Put("/:id", handlers.UpdateUser())
	// userGroup.Delete("/:id", handlers.DeleteUser())

	// Main API
	v1 := app.Group("/v1")

	// Main Handlers
	dashboard := v1.Group("/dashboard", middlewares.ExtractRole(), middlewares.ExtractUUID())
	dashboard.Get("", handlers.Dashboard())

	// dmin Handlers
	admin := v1.Group("/admin", middlewares.ProtectedByRole("admin"))
	admin.Get("/dashboard", handlers.Dashboard())

	// Mod Handlers
	mod := v1.Group("/mod", middlewares.ProtectedByRole("mod"))
	mod.Get("/dashboard", handlers.Dashboard())

	// User Handlers
	user := v1.Group("/auth")
	user.Post("", handlers.Register())
	user.Get("", handlers.Login())

	// Tests Handlers
	test := app.Group("/test")
	test.Get("/test", handlers.TestBody())
	test.Get("/testing", handlers.TestConnection())
	test.Static("/file", "./test.txt")

}
