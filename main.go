package main

import (
	"fiberproject/api/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {

	//Apì init
	app := fiber.New()

	// Routes
	routes.Setup(app)

	//Port
	app.Listen(":8080")
}
