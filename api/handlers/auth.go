package handlers

import (
	"context"
	"fiberproject/api/models"
	"fiberproject/db"
	"fiberproject/pkg/jwt"
	"fiberproject/pkg/util/hashing"
	"fiberproject/pkg/util/httputility"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

// Register registers a new user.
//
// @Summary Register a new user.
// @Description Registers a new user with the provided username, email, and password.
// @Tags auth
// @Accept json
// @Produce json
// @Param user body httputility.RegisterRequest true "User object containing username, email, and password"
// @Success 201 {object} httputility.IDResponse "User registered successfully"
// @Failure 400 {object} httputility.HTTPError "Invalid request body or fields cannot be empty"
// @Failure 500
// @Router /v1/auth/register [post]
func Register() fiber.Handler {
	return func(c *fiber.Ctx) error {

		var request httputility.RegisterRequest

		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(httputility.HTTPError{
				Message: "invalid request body",
			})
		}

		if request.Username == "" || request.Email == "" || request.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(httputility.HTTPError{
				Message: "fields cannot be empty",
			})
		}

		conn, err := db.ConnectDB()
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		defer conn.Close(context.Background())

		hashedPassword, err := hashing.HashPassword(request.Password)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		var uuid string

		err = conn.QueryRow(context.Background(),
			"INSERT INTO users (username, password, email) VALUES ($1, $2, $3)",
			request.Username, hashedPassword, request.Email).Scan(&uuid)
		if err != nil {

			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.Status(fiber.StatusCreated).JSON(httputility.IDResponse{UUID: uuid})
	}
}

// Login logs in a user.
//
// @Summary Log in a user.
// @Description Logs in a user with the provided username and password.
// @Tags auth
// @Accept json
// @Produce json
// @Param user body httputility.LoginRequest true "User object containing username and password"
// @Success 200 {object} httputility.TokenResponse "User logged in successfully"
// @Failure 400 {object} httputility.HTTPError "Invalid request body or fields cannot be empty"
// @Failure 401 {object} httputility.HTTPError "Unauthorized or invalid credentials"
// @Failure 404 {object} httputility.HTTPError "User not found"
// @Failure 500
// @Router /v1/auth/login [post]
func Login() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request httputility.LoginRequest

		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(httputility.HTTPError{
				Message: "invalid request body",
			})
		}

		if request.Username == "" || request.Password == "" {
			return c.Status(fiber.StatusBadRequest).JSON(httputility.HTTPError{
				Message: "fields cannot be empty",
			})
		}

		conn, err := db.ConnectDB()
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		defer conn.Close(context.Background())

		var user models.User

		err = conn.QueryRow(context.Background(),
			"SELECT uuid, username, password, role FROM users WHERE username = $1",
			request.Username).Scan(&user.UUID, &user.Username, &user.Password, &user.Role)
		if err != nil {
			if err == pgx.ErrNoRows {
				return c.Status(fiber.StatusNotFound).JSON(httputility.HTTPError{
					Message: "user not found",
				})
			}
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(httputility.HTTPError{
				Message: "invalid credentials",
			})

		}

		token, err := jwt.CreateNewToken(user.UUID, user.Username, user.Role)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.Status(fiber.StatusOK).JSON(httputility.TokenResponse{Token: token})
	}
}
