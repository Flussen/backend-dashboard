package staffactions

import (
	"context"
	"fiberproject/db"
	"fiberproject/pkg/util/httputility"
	httputility_staff "fiberproject/pkg/util/httputility/staff"

	"github.com/gofiber/fiber/v2"
)

// SetRole assigns a new role to the specified user.
//
// @Summary Assign user role.
// @Description Assigns a new role to the specified user identified by the provided UserID.
// @Tags staff
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body httputility_staff.SetRoleRequest true "Request body containing Role and UserID"
// @Success 200
// @Failure 400 {object} httputility.HTTPError "Invalid request body or missing parameters"
// @Failure 404 {object} httputility.HTTPError "User not found or no changes made"
// @Failure 500
// @Router /v1/admin/actions/roles [put]
func SetRole() fiber.Handler {
	return func(c *fiber.Ctx) error {

		var request httputility_staff.SetRoleRequest

		if err := c.BodyParser(&request); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		if request.Role == "" || request.UserID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(httputility.HTTPError{
				Message: httputility.ErrNoEmpty,
			})
		}

		conn, err := db.ConnectDB()
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		defer conn.Close(context.Background())

		commandTag, err := conn.Exec(context.Background(),
			"UPDATE users SET role = $1 WHERE uuid = $2", &request.Role, &request.UserID)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if commandTag.RowsAffected() == 0 {
			return c.Status(fiber.StatusNotFound).JSON(httputility.HTTPError{
				Message: "User not found or no changes made",
			})

		}

		return c.SendStatus(fiber.StatusOK)
	}
}
