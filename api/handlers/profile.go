package handlers

import (
	"context"
	"fiberproject/api/models"
	"fiberproject/db"
	"fiberproject/pkg/util/hashing"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func GetProfile() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uuid := c.Params("id")

		conn, err := db.ConnectDB()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}
		defer conn.Close(context.Background())

		var user models.User

		err = conn.QueryRow(context.Background(),
			"SELECT uuid, username, email, role FROM users WHERE uuid = $1", uuid).
			Scan(&user.UUID, &user.Username, &user.Email, &user.Role)
		if err != nil {
			if err == pgx.ErrNoRows {
				return c.Status(fiber.StatusNotFound).SendString("User not found")
			}
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Error querying the database : %v", err))
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"uuid":     uuid,
			"username": &user.Username,
			"email":    &user.Email,
			"role":     &user.Role,
		})
	}
}

func UpdateProfile() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uuid := c.Params("id")
		var newUser models.User
		if err := c.BodyParser(&newUser); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
		}

		conn, err := db.ConnectDB()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}
		defer conn.Close(context.Background())

		var (
			username string
			email    string
		)

		err = conn.QueryRow(context.Background(),
			"SELECT username, email FROM users WHERE uuid = $1", uuid).
			Scan(&username, &email)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error: Unable to get the user")
		}

		if newUser.Username == "" {
			newUser.Username = username
		}

		if newUser.Email == "" {
			newUser.Email = email
		}

		commandTag, err := conn.Exec(context.Background(),
			"UPDATE users SET username = $1, email = $2 WHERE uuid = $3", newUser.Username, newUser.Email, uuid)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error: Unable to update user")
		}

		if commandTag.RowsAffected() == 0 {
			return c.Status(fiber.StatusNotFound).SendString("User not found or no changes made")

		}

		newUser.UUID = uuid
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"uuid":     newUser.UUID,
			"username": newUser.Username,
			"email":    newUser.Email,
		})
	}
}

func ChangePassword() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uuid := c.Params("id")

		type Receive struct {
			Oldpassword string `json:"oldpassword"`
			Newpassword string `json:"newpassword"`
		}

		var data Receive
		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
		}

		if data.Newpassword == "" || data.Oldpassword == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
		}

		conn, err := db.ConnectDB()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}
		defer conn.Close(context.Background())

		var passwordSaved string

		err = conn.QueryRow(context.Background(),
			"SELECT password FROM users WHERE uuid = $1", uuid).
			Scan(&passwordSaved)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error: Unable to get the user")
		}

		err = bcrypt.CompareHashAndPassword([]byte(passwordSaved), []byte(data.Oldpassword))
		if err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		hashed, err := hashing.HashPassword(data.Newpassword)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}

		commandTag, err := conn.Exec(context.Background(),
			"UPDATE users SET password = $1 WHERE uuid = $2", hashed, uuid)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error: Unable to update the user")
		}

		if commandTag.RowsAffected() == 0 {
			return c.Status(fiber.StatusNotFound).SendString("User not found or no changes made")

		}

		return c.SendStatus(fiber.StatusOK)
	}
}

func DeleteAccount() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uuid := c.Params("id")

		conn, err := db.ConnectDB()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}
		defer conn.Close(context.Background())

		commandTag, err := conn.Exec(context.Background(),
			"DELETE FROM users WHERE uuid = $1", uuid)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error: Unable to delete the user")
		}

		if commandTag.RowsAffected() == 0 {
			return c.Status(fiber.StatusNotFound).SendString("User not found or no changes made")

		}

		rtrn := fmt.Sprintf("Deleted user with UUID %s", uuid)
		return c.Status(fiber.StatusOK).SendString(rtrn)
	}
}
