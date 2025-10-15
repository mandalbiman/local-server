package handlers

import (
	"strconv"

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
