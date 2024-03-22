package main

import (
	"fiberproject/api/routes"
	"fiberproject/config"

	"github.com/gofiber/fiber/v2"
)

func main() {

	// Load configs
	config.ConfigInit()

	//Apì init
	app := fiber.New()

	// Routes
	routes.Setup(app)

	//Port
	app.Listen(":8080")
}
