package handlers

import (
	"fmt"
	"strconv"
	"time"

	"bats.com/local-server/api/models"
	"bats.com/local-server/io/db"
	"github.com/gofiber/fiber/v2"
)

type UsersHandler struct {
	db *db.UserDB
}

func SetUpUserRoutes(v1 fiber.Router) {
	userDB := db.GetUserDB()
	uh := UsersHandler{db: userDB}
	userGrp := v1.Group("/users")
	userGrp.Get("/", uh.getUsers)
	userGrp.Get("/:id", uh.getUserById)
	userGrp.Put("/:id", uh.updateUser)
	userGrp.Post("/", uh.createUser)
	userGrp.Delete("/:id", uh.deleteUser)
	userGrp.Get("/search/:username", uh.searchUser)
}

func (u UsersHandler) getUsers(c *fiber.Ctx) error {
	users, err := u.db.GetUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}
	return c.JSON(users)
}

func (u UsersHandler) getUserById(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}
	user, err := u.db.GetUserByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}
	return c.JSON(user)
}

func (u UsersHandler) updateUser(c *fiber.Ctx) error {
	// Parse ID param
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Check if user exists
	existingUser, err := u.db.GetUserByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Parse request body
	var updateData models.User
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update fields
	existingUser.Name = updateData.Name
	existingUser.Username = updateData.Username
	existingUser.Email = updateData.Email
	existingUser.Password = updateData.Password
	existingUser.Age = updateData.Age
	existingUser.Gender = updateData.Gender
	existingUser.Phone = updateData.Phone
	existingUser.City = updateData.City
	existingUser.Country = updateData.Country

	if _, err := u.db.UpdateUser(*existingUser); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	return c.JSON(existingUser)
}

func (u UsersHandler) deleteUser(c *fiber.Ctx) error {
	// Parse ID
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Delete user
	user, err := u.db.DeleteUser(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("User %s deleted successfully", user.Username),
	})
}

func (u UsersHandler) searchUser(c *fiber.Ctx) error {
	username := c.Params("username")
	user, err := u.db.SearchUserByUsername(username)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(user)
}

func (u UsersHandler) createUser(c *fiber.Ctx) error {
	// Parse request body
	var newUser models.User
	if err := c.BodyParser(&newUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Optional: basic validation
	if newUser.Username == "" || newUser.Email == "" || newUser.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username, email, and password are required",
		})
	}

	// Set creation timestamp
	newUser.CreatedAt = time.Now()

	// Save to DB
	createdUser, err := u.db.CreateUser(newUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(createdUser)
}
