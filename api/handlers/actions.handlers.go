package handlers

import (
	"context"
	"fiberproject/api/models"
	"fiberproject/db"
	"fiberproject/pkg/util/hashing"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

func GetUsers() fiber.Handler {
	return func(c *fiber.Ctx) error {

		conn, err := db.ConnectDB()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}
		defer conn.Close(context.Background())

		var users []models.User

		rows, err := conn.Query(context.Background(), "SELECT uuid, username, password, email FROM users")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}

		defer rows.Close()

		for rows.Next() {
			var user models.User
			if err := rows.Scan(&user.UUID, &user.Username, &user.Password, &user.Email); err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("Error scanning user")
			}
			users = append(users, user)
		}

		if err := rows.Err(); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error after iterating over rows")
		}

		return c.Status(fiber.StatusOK).JSON(&users)
	}
}

func GetUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uuid := c.Params("id")

		conn, err := db.ConnectDB()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}
		defer conn.Close(context.Background())

		var user models.User

		err = conn.QueryRow(context.Background(),
			"SELECT uuid, username, password, email FROM users WHERE uuid = $1", uuid).
			Scan(&user.UUID, &user.Username, &user.Password, &user.Email)
		if err != nil {
			if err == pgx.ErrNoRows {
				return c.Status(fiber.StatusNotFound).SendString("User not found")
			}
			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Error querying the database : %v", err))
		}

		return c.Status(fiber.StatusOK).JSON(&user)
	}
}

func CreateUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var user models.User

		if err := c.BodyParser(&user); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
		}

		if user.Username == "" || user.Email == "" || user.Password == "" {
			return c.Status(fiber.StatusBadRequest).SendString("invalid request body")
		}

		conn, err := db.ConnectDB()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}
		defer conn.Close(context.Background())

		hashedPassword, err := hashing.HashPassword(user.Password)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}

		user.Password = hashedPassword

		err = conn.QueryRow(context.Background(),
			"INSERT INTO users (username, password, email) VALUES ($1, $2, $3) RETURNING uuid",
			user.Username, user.Password, user.Email).Scan(&user.UUID)
		if err != nil {

			return c.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Internal Server Error: Unable to create user : %v", err))
		}

		return c.Status(fiber.StatusCreated).JSON(&user)
	}
}

func UpdateUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uuid := c.Params("id")

		var newUser models.User

		if err := c.BodyParser(&newUser); err != nil {
			invalidRequestBody := fmt.Sprintf("Invalid request body, %v", err)
			return c.Status(fiber.StatusBadRequest).SendString(invalidRequestBody)
		}

		conn, err := db.ConnectDB()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}
		defer conn.Close(context.Background())

		var (
			username string
			password string
			email    string
		)

		err = conn.QueryRow(context.Background(),
			"SELECT username, password, email FROM users WHERE uuid = $1", uuid).
			Scan(&username, &password, &email)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error: Unable to update user")
		}

		var hashed string

		if newUser.Username == "" {
			newUser.Username = username
		}

		if newUser.Password == "" {
			newUser.Password = password
		} else {
			hash, err := hashing.HashPassword(newUser.Password)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error: Unable to update user")
			}
			hashed = hash
		}

		if newUser.Email == "" {
			newUser.Email = email
		}

		if hashed != "" {
			newUser.Password = hashed
		}

		_, err = conn.Exec(context.Background(),
			"UPDATE users SET username = $1, password = $2, email = $3 WHERE uuid = $4",
			newUser.Username, newUser.Password, newUser.Email, uuid)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error: Unable to update user")
		}
		newUser.UUID = uuid
		return c.Status(fiber.StatusOK).JSON(&newUser)
	}
}

func DeleteUser() fiber.Handler {
	return func(c *fiber.Ctx) error {
		uuid := c.Params("id")

		conn, err := db.ConnectDB()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
		}
		defer conn.Close(context.Background())

		_, err = conn.Exec(context.Background(),
			"DELETE FROM users WHERE uuid = $1", uuid)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error: Unable to delete user")
		}
		returnString := fmt.Sprintf("Deleted user with UUID %s", uuid)
		return c.Status(fiber.StatusOK).SendString(returnString)
	}
}

func TestBody() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var user models.User

		if err := c.BodyParser(&user); err != nil {
			invalidRequestBody := fmt.Sprintf("Invalid request body, %v", err)
			return c.Status(fiber.StatusBadRequest).SendString(invalidRequestBody)
		}

		return c.Status(200).JSON(&user)
	}
}

func TestConnection() fiber.Handler {
	return func(c *fiber.Ctx) error {

		conn, err := db.ConnectDB()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error connecting to database")
		}
		defer conn.Close(context.Background())

		// Get the current database user
		var currentUser string
		row := conn.QueryRow(context.Background(), "SELECT current_user;")
		err = row.Scan(&currentUser)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error fetching current user")
		}

		// Get the list of tables
		rows, err := conn.Query(context.Background(), "SELECT table_schema, table_name FROM information_schema.tables WHERE table_schema NOT IN ('information_schema', 'pg_catalog') AND table_type = 'BASE TABLE' ORDER BY table_schema, table_name;")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error fetching tables")
		}
		defer rows.Close()

		var tables []string
		for rows.Next() {
			var schema, name string
			if err := rows.Scan(&schema, &name); err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("Error scanning tables")
			}
			tables = append(tables, schema+"."+name)
		}

		if err := rows.Err(); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Error iterating over tables")
		}

		// Here you could return the currentUser and tables information in your response
		// For example, you might return them as JSON

		return c.Status(fiber.StatusOK).SendString("Current user: " + currentUser + ", Tables: " + strings.Join(tables, ", "))
	}
}
