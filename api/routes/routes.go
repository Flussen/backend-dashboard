package routes

import (
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	v1 := app.Group("/v1")

	//
	dashboard(v1)
	//
	staffdashboard(v1)
	//
	auth(v1)
	//
	tests(v1)

}
