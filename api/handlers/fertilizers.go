package handlers

import (
	"fmt"
	"strconv"
	"time"

	"bats.com/local-server/api/models"
	"bats.com/local-server/io/db"
	"bats.com/local-server/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type FertilizerHandler struct {
	db *db.FertilizerDB
}

func SetUpfertilizerRoutes(v1 fiber.Router) {
	fertilizerDB := db.GetFertilizerDB()
	uh := FertilizerHandler{db: fertilizerDB}
	fertilizerGrp := v1.Group("/fertilizers", middleware.JWTProtectedMiddleware)
	fertilizerGrp.Get("/", uh.getFertilizer)
	fertilizerGrp.Get("/:id", uh.getFertilizerById)
	fertilizerGrp.Put("/:id", uh.updateFertilizer)
	fertilizerGrp.Post("/", uh.createFertilizer)
	fertilizerGrp.Delete("/:id", uh.deleteFertilizer)
	fertilizerGrp.Get("/search/:fertilizername", uh.searchFertilizer)
}

func (u FertilizerHandler) getFertilizer(c *fiber.Ctx) error {
	fertilizers, err := u.db.GetFertilizers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch fertilizers",
		})
	}
	return c.JSON(fertilizers)
}

func (u FertilizerHandler) getFertilizerById(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid fertilizer ID",
		})
	}
	fertilizer, err := u.db.GetFertilizerByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Fertilizer not found",
		})
	}
	return c.JSON(fertilizer)
}

func (u FertilizerHandler) updateFertilizer(c *fiber.Ctx) error {
	// Parse ID param
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid fertilizer ID",
		})
	}

	// Check if fertilizer exists
	existingFertilizer, err := u.db.GetFertilizerByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Fertilizer not found",
		})
	}

	// Parse request body
	var updateData models.Fertilizer
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update fields
	existingFertilizer.ProductName = updateData.ProductName
	existingFertilizer.Category = updateData.Category
	existingFertilizer.Dosage = updateData.Dosage
	existingFertilizer.UniqueID = updateData.UniqueID
	existingFertilizer.BatchNumber = updateData.BatchNumber
	existingFertilizer.ManufactureDate = updateData.ManufactureDate
	existingFertilizer.ExpiryDate = updateData.ExpiryDate
	existingFertilizer.CautionaryLogo = updateData.CautionaryLogo
	existingFertilizer.CautionaryText = updateData.CautionaryText
	existingFertilizer.AntidoteStatement = updateData.AntidoteStatement

	if _, err := u.db.UpdateFertilizer(*existingFertilizer); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update fertilizer",
		})
	}

	return c.JSON(existingFertilizer)
}

func (u FertilizerHandler) deleteFertilizer(c *fiber.Ctx) error {
	// Parse ID
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid fertilizer ID",
		})
	}

	// Delete fertilizer
	fertilizer, err := u.db.DeleteFertilizer(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Fertilizer not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Fertilizer %s deleted successfully", fertilizer.ProductName),
	})
}

func (u FertilizerHandler) searchFertilizer(c *fiber.Ctx) error {
	fertilizername := c.Params("fertilizername")
	fertilizer, err := u.db.SearchByProductName(fertilizername)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Fertilizer not found",
		})
	}

	return c.JSON(fertilizer)
}

func (u FertilizerHandler) createFertilizer(c *fiber.Ctx) error {
	// Parse request body
	var newFertilizer models.Fertilizer
	if err := c.BodyParser(&newFertilizer); err != nil {
		log.Error(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Set creation timestamp
	newFertilizer.CreatedAt = time.Now()

	// Save to DB
	createdFertilizer, err := u.db.CreateFertilizer(newFertilizer)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create fertilizer",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(createdFertilizer)
}
