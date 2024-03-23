package handlers

import (
	"fiberproject/api/functions"
	"fiberproject/pkg/util/httputility"

	"github.com/gofiber/fiber/v2"
)

// Dashboard 	 Shows the main user dashboard by role
// @Summary      Show the Dashbord by the role
// @Description  will return role if authenticated, if not authenticated it will send an error
// @Tags         dashboard
// @Produce      json
// @Security BearerAuth
// @Success      200  {object}  httputility.DashboardResponse
// @Failure      401  {object}  httputility.HTTPError
// @Failure      404  {object}  httputility.HTTPError
// @Router       /v1/dashboard [get]
func Dashboard() fiber.Handler {
	return func(c *fiber.Ctx) error {

		uuid, ok := c.Locals("uuid").(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(httputility.HTTPError{
				Message: "Unauthorized",
			})
		}

		user, err := functions.GetInfoByUUID(uuid)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(httputility.HTTPError{
				Message: "user not found",
			})
		}

		switch user.Role {
		case "user", "admin", "mod":
			return c.Status(fiber.StatusOK).JSON(httputility.DashboardResponse{
				UUID:     user.UUID,
				Username: user.Username,
				Role:     user.Role,
			})
		default:
			return c.Status(fiber.StatusUnauthorized).JSON(httputility.HTTPError{
				Message: "Unauthorized",
			})
		}
	}
}
