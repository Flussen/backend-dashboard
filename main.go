package main

import (
	"fiberproject/api/routes"
	"fiberproject/config"

	_ "fiberproject/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger" // swagger handler
)

// @title Styerr GO Api
// @version 1.0
// @description Styerr network internal api
// @contact.name flussen in discord
// @host localhost:8080
// @SecurityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @BasePath /
func main() {

	// Load configs
	config.ConfigInit()

	//Apì init
	app := fiber.New()
	app.Get("/swagger/*", swagger.HandlerDefault) // default

	// Routes
	routes.Setup(app)

	//Port
	app.Listen(":8080")
}
