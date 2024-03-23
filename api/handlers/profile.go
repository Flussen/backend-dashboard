package handlers

import (
	"context"
	"fiberproject/api/models"
	"fiberproject/db"
	"fiberproject/pkg/util/hashing"
	"fiberproject/pkg/util/httputility"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

// GetProfile retrieves user profile information.
//
// @Summary Retrieve user profile.
// @Description Retrieves the user profile information for the specified user UUID.
// @Tags profile
// @Accept json
// @Produce json
// @Param id path string true "User UUID"
// @Success 200 {object} models.User "User profile retrieved successfully"
// @Failure 404 {object} httputility.HTTPError "User not found"
// @Failure 500
// @Router /v1/dashboard/profile/{id} [get]
func GetProfile() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uuid := c.Params("id")

		conn, err := db.ConnectDB()
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		defer conn.Close(context.Background())

		var user models.User

		err = conn.QueryRow(context.Background(),
			"SELECT uuid, username, email, role FROM users WHERE uuid = $1", uuid).
			Scan(&user.UUID, &user.Username, &user.Email, &user.Role)
		if err != nil {
			if err == pgx.ErrNoRows {
				return c.Status(fiber.StatusNotFound).JSON(httputility.HTTPError{
					Message: "user not found",
				})
			}
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.Status(fiber.StatusOK).JSON(&user)
	}
}

// UpdateProfile updates user profile information.
//
// @Summary Update user profile.
// @Description Updates the user profile information with the provided username and email.
// @Tags profile
// @Accept json
// @Produce json
// @Param id path string true "User UUID"
// @Param user body httputility.UpdateRequest true "User object containing username and email"
// @Success 200 {object} httputility.UpdateResponse "User profile updated successfully"
// @Failure 400 {object} httputility.HTTPError "Invalid request body"
// @Failure 500
// @Router /v1/dashboard/profile [put]
func UpdateProfile() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uuid := c.Params("id")

		var request httputility.UpdateRequest

		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(httputility.HTTPError{
				Message: "invalid request body",
			})
		}

		var (
			username string
			email    string
		)

		conn, err := db.ConnectDB()
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		defer conn.Close(context.Background())

		err = conn.QueryRow(context.Background(),
			"SELECT username, email FROM users WHERE uuid = $1", uuid).
			Scan(&username, &email)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if request.Username == "" {
			request.Username = username
		}

		if request.Email == "" {
			request.Email = email
		}

		_, err = conn.Exec(context.Background(),
			"UPDATE users SET username = COALESCE(NULLIF($1, ''), username), email = COALESCE(NULLIF($2, ''), email) WHERE uuid = $3 RETURNING username, email",
			request.Username, request.Email, uuid)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error: Unable to update user")
		}

		return c.Status(fiber.StatusOK).JSON(httputility.UpdateResponse{
			UUID:     uuid,
			Username: request.Username,
			Email:    request.Email,
		})
	}
}

// ChangePassword changes the password for the specified user.
//
// @Summary Change user password.
// @Description Changes the password for the specified user with the provided old and new passwords.
// @Tags profile
// @Accept json
// @Produce json
// @Param id path string true "User UUID"
// @Param data body httputility.ChangePasswordRequest true "Old and new password data"
// @Success 200
// @Failure 400 {object} httputility.HTTPError "Invalid request body"
// @Failure 401 {object} httputility.HTTPError "Unauthorized or invalid credentials"
// @Failure 404 {object} httputility.HTTPError "User not found or no changes made"
// @Failure 500
// @Router /v1/dashboard/profile/{id}/pwd [put]
func ChangePassword() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uuid := c.Params("id")

		var request httputility.ChangePasswordRequest
		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(httputility.HTTPError{
				Message: "Invalid request body",
			})
		}

		if request.Newpassword == "" || request.Oldpassword == "" {
			return c.Status(fiber.StatusBadRequest).JSON(httputility.HTTPError{
				Message: "Invalid request body",
			})
		}

		conn, err := db.ConnectDB()
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		defer conn.Close(context.Background())

		var passwordSaved string

		err = conn.QueryRow(context.Background(),
			"SELECT password FROM users WHERE uuid = $1", uuid).
			Scan(&passwordSaved)

		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		err = bcrypt.CompareHashAndPassword([]byte(passwordSaved), []byte(request.Oldpassword))
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(httputility.HTTPError{
				Message: "invalid credentials",
			})

		}

		hashed, err := hashing.HashPassword(request.Newpassword)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		commandTag, err := conn.Exec(context.Background(),
			"UPDATE users SET password = $1 WHERE uuid = $2", hashed, uuid)
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

// DeleteAccount deletes the user account associated with the specified UUID.
//
// @Summary Delete user account.
// @Description Deletes the user account associated with the specified UUID.
// @Tags profile
// @Accept json
// @Produce json
// @Param id path string true "User UUID"
// @Success 200 {object} httputility.SendMessage "Deleted user with UUID id"
// @Failure 404 {object} httputility.HTTPError "User not found or no changes made"
// @Failure 500
// @Router /v1/dashboard/profile/{id} [delete]
func DeleteAccount() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uuid := c.Params("id")

		conn, err := db.ConnectDB()
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		defer conn.Close(context.Background())

		commandTag, err := conn.Exec(context.Background(),
			"DELETE FROM users WHERE uuid = $1", uuid)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if commandTag.RowsAffected() == 0 {
			return c.Status(fiber.StatusNotFound).JSON(httputility.HTTPError{
				Message: "User not found or no changes made",
			})

		}

		rtrn := fmt.Sprintf("Deleted user with UUID %s", uuid)
		return c.Status(fiber.StatusOK).JSON(httputility.SendMessage{
			Message: rtrn,
		})
	}
}
