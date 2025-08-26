package handlers

import (
	"fmt"
	"os"
	"time"

	"go_taskmanagement/database"
	"go_taskmanagement/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// RegisterRequest kullanıcı kayıt isteği modeli
type RegisterRequest struct {
	Username string `json:"username" example:"hakan"`
	Email    string `json:"email" example:"hakan@example.com"`
	Password string `json:"password" example:"1234"`
}

// LoginRequest kullanıcı giriş isteği modeli
type LoginRequest struct {
	Email    string `json:"email" example:"hakan@example.com"`
	Password string `json:"password" example:"1234"`
}

// Get JWT secret from environment variable
func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "gizliAnahtar" // Fallback for development
	}
	return []byte(secret)
}

// RegisterHandler kullanıcı kaydı oluşturur
// @Summary Kullanıcı kaydı
// @Description Yeni kullanıcı oluşturur
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "Kullanıcı kayıt bilgileri"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Router /register [post]
// @ID RegisterHandler
func RegisterHandler(c *fiber.Ctx) error {
	var input RegisterRequest

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz veri"})
	}

	if input.Username == "" || input.Email == "" || input.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Kullanıcı adı, email ve şifre zorunlu"})
	}

	// Database mode
	if database.IsConnected {
		// Check if user already exists
		var existingUser models.User
		if err := database.DB.Where("username = ? OR email = ?", input.Username, input.Email).First(&existingUser).Error; err == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Kullanıcı adı veya email zaten mevcut"})
		}

		// Hash password
		hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Şifre hashlenemedi"})
		}

		// Create user
		user := models.User{
			Username: input.Username,
			Email:    input.Email,
			Password: string(hash),
		}

		if err := database.DB.Create(&user).Error; err != nil {
			// Log the actual error for debugging
			fmt.Printf("Database error creating user: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Kullanıcı oluşturulamadı"})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Kullanıcı başarıyla oluşturuldu",
			"user": fiber.Map{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
			},
		})
	}

	// In-memory mode (fallback)
	// Check if user already exists
	for _, u := range models.Users {
		if u.Username == input.Username || u.Email == input.Email {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Kullanıcı adı veya email zaten mevcut"})
		}
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Şifre hashlenemedi"})
	}

	// Create user
	user := models.User{
		ID:       uint(len(models.Users) + 1),
		Username: input.Username,
		Email:    input.Email,
		Password: string(hash),
	}

	models.Users = append(models.Users, user)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Kullanıcı başarıyla oluşturuldu",
		"user": fiber.Map{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// LoginHandler kullanıcı girişi yapar ve JWT token döner
// @Summary Kullanıcı girişi
// @Description Email ve şifre ile giriş yapar
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Email ve şifre"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /login [post]
// @ID LoginHandler
func LoginHandler(c *fiber.Ctx) error {
	var input LoginRequest

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz veri"})
	}

	if input.Email == "" || input.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email ve şifre zorunlu"})
	}

	var user models.User
	var found bool

	// Database mode
	if database.IsConnected {
		if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Email veya şifre yanlış"})
		}
		found = true
	} else {
		// In-memory mode (fallback)
		for _, u := range models.Users {
			if u.Email == input.Email {
				user = u
				found = true
				break
			}
		}
	}

	if !found {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Email veya şifre yanlış"})
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Email veya şifre yanlış"})
	}

	// Create JWT token
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(getJWTSecret())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Token oluşturulamadı"})
	}

	return c.JSON(fiber.Map{
		"token": tokenString,
		"user": fiber.Map{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// LogoutHandler kullanıcıyı çıkış yaptırır
// @Summary Çıkış
// @Description Kullanıcıyı çıkış yaptırır
// @Tags Auth
// @Produce json
// @Success 200 {object} map[string]string
// @Security BearerAuth
// @Router /logout [post]
// @ID LogoutHandler
func LogoutHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Çıkış başarılı. Token client tarafından silinmeli."})
}
