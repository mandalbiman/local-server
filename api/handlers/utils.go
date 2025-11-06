package handlers

import (
	"encoding/json"
	"fmt"

	"bats.com/local-server/api/models"
	"bats.com/local-server/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func validationErr(err error, c *fiber.Ctx) fiber.Map {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return fiber.Map{
			"error": "Invalid validation error format",
		}
	}

	// Format validation error messages as structured objects
	var errorDetails []fiber.Map
	for _, e := range validationErrors {
		errorDetails = append(errorDetails, fiber.Map{
			"field":   e.Field(),
			"tag":     e.Tag(),
			"message": fmt.Sprintf("The '%s' field must satisfy the '%s' condition", e.Field(), e.Tag()),
		})
	}

	// Return JSON response with structured error details
	return fiber.Map{
		"status":  "error",
		"message": "Validation failed",
		"errors":  errorDetails,
	}
}

func isValidJSON(s string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(s), &js) == nil
}

func parseAndValidateRequest(c *fiber.Ctx, req interface{}, validator *validator.Validate) interface{} {
	fmt.Println(string(c.Body()))
	// Parse the request body
	if err := c.BodyParser(req); err != nil {
		fmt.Println(err)
		errorResponse := models.N400BadResponse{
			Error:   utils.StringPtr("Invalid JSON Structure"),
			Message: utils.StringPtr("Ensure valid JSON format"),
		}
		return errorResponse
	}
	// Validate the request struct
	if err := validator.Struct(req); err != nil {
		return validationErr(err, c)
	}

	return nil
}
