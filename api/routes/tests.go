package routes

import (
	"fiberproject/api/handlers"

	"github.com/gofiber/fiber/v2"
)

func tests(api fiber.Router) {

	test := api.Group("/test")

	test.Get("/test", handlers.TestBody())
	test.Get("/testing", handlers.TestConnection())
	test.Static("/file", "./test.txt")
}
