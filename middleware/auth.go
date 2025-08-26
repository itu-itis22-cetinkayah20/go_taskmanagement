package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Get JWT secret from environment variable
func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "gizliAnahtar" // Fallback for development
	}
	return []byte(secret)
}

// AuthMiddleware JWT doğrulaması yapar
func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token gerekli"})
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return getJWTSecret(), nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Geçersiz token"})
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// Store both user_id and other claims for handlers
		if uidFloat, ok := claims["user_id"].(float64); ok {
			c.Locals("user_id", uint(uidFloat))
		}
		if username, ok := claims["username"].(string); ok {
			c.Locals("username", username)
		}
		if email, ok := claims["email"].(string); ok {
			c.Locals("email", email)
		}
	}
	return c.Next()
}
