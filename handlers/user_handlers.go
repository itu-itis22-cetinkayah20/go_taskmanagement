package handlers

import (
	"go_taskmanagement/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("gizliAnahtar")

// RegisterHandler kullanıcı kaydı oluşturur
// @Summary Kullanıcı kaydı
// @Description Yeni kullanıcı oluşturur
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body models.User true "Kullanıcı"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Router /register [post]
// @ID RegisterHandler
func RegisterHandler(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz veri"})
	}
	if user.Username == "" || user.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Kullanıcı adı ve şifre zorunlu"})
	}
	// Kullanıcı adı benzersiz mi kontrolü
	for _, u := range models.Users {
		if u.Username == user.Username {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Kullanıcı adı zaten mevcut"})
		}
	}
	// Şifreyi hashle
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Şifre hashlenemedi"})
	}
	user.Password = string(hash)
	user.ID = len(models.Users) + 1
	models.Users = append(models.Users, user)
	return c.Status(fiber.StatusCreated).JSON(user)
}

// LoginHandler kullanıcı girişi yapar ve JWT token döner
// @Summary Kullanıcı girişi
// @Description Kullanıcı adı ve şifre ile giriş yapar
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body models.User true "Kullanıcı"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /login [post]
// @ID LoginHandler
func LoginHandler(c *fiber.Ctx) error {
	var creds models.User
	if err := c.BodyParser(&creds); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz veri"})
	}
	var user models.User
	found := false
	for _, u := range models.Users {
		if u.Username == creds.Username {
			user = u
			found = true
			break
		}
	}
	if !found {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Kullanıcı adı veya şifre yanlış"})
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Kullanıcı adı veya şifre yanlış"})
	}
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Token oluşturulamadı"})
	}
	return c.JSON(fiber.Map{"token": tokenString})
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
