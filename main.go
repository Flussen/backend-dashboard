package main

import (
	"fiberproject/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func main() {
	app := fiber.New()

	// User handlers
	app.Use(requestid.New())

	// User Routes
	userGroup := app.Group("/users")
	userGroup.Get("", handlers.GetUsers())
	userGroup.Get("/:id", handlers.GetUser())
	userGroup.Post("", handlers.CreateUser())
	userGroup.Put("/:id", handlers.UpdateUser())
	userGroup.Delete("/:id", handlers.DeleteUser())

	// test
	app.Get("/test", handlers.TestBody())
	app.Get("/testing", handlers.TestConnection())
	app.Static("/file", "./test.txt")

	app.Listen(":8080")
}
