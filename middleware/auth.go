package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("gizliAnahtar")

// AuthMiddleware JWT doğrulaması yapar
func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token gerekli"})
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Geçersiz token"})
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if uidFloat, ok := claims["user_id"].(float64); ok {
			c.Locals("user_id", int(uidFloat))
		}
	}
	return c.Next()
}
