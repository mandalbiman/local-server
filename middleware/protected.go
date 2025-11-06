package middleware

import (
	"bats.com/local-server/api/svc"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTProtectedMiddleware(c *fiber.Ctx) error {
	//c.Locals("userId", "test")
	//c.Locals("roles", []string{"admin"})
	//return c.Next()
	// Get token from the header
	authHeader := c.Get("X-F-Token")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid token"})
	}

	tokenString := authHeader
	token, err := svc.ValidateJWT(tokenString)
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	// Extract user claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
	}

	// Extract user_id
	userID, _ := claims["username"].(string)

	// Extract roles (handling array of strings)
	var roles []string
	if rolesInterface, exists := claims["roles"]; exists {
		if rolesArray, ok := rolesInterface.([]interface{}); ok {
			for _, role := range rolesArray {
				if roleStr, ok := role.(string); ok {
					roles = append(roles, roleStr)
				}
			}
		}
	}

	// Store user credentials in context
	c.Locals("userId", userID)
	c.Locals("roles", roles)

	return c.Next()
}
