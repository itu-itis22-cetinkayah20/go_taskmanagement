package handlers

import (
	"go_taskmanagement/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// PublicTasksHandler herkese açık görevleri listeler
// @ID PublicTasksHandler
// @Summary Public görevleri listele
// @Description Herkesin görebileceği görevleri döner
// @Tags Tasks
// @Produce json
// @Success 200 {array} models.Task
// @Router /tasks/public [get]
func PublicTasksHandler(c *fiber.Ctx) error {
	return c.JSON(models.PublicTasks)
}

// TasksListHandler kullanıcının kendi görevlerini listeler
// @ID TasksListHandler
// @Summary Kullanıcı görevlerini listele
// @Description Sadece giriş yapan kullanıcının görevlerini döner
// @Tags Tasks
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Task
// @Router /tasks [get]
func TasksListHandler(c *fiber.Ctx) error {
	uid := c.Locals("user_id")
	userID, ok := uid.(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Kullanıcı bilgisi alınamadı"})
	}
	var userTasks []models.Task
	for _, t := range models.Tasks {
		if t.UserID == userID {
			userTasks = append(userTasks, t)
		}
	}
	return c.JSON(userTasks)
}

// TaskCreateHandler yeni görev ekler
// @ID TaskCreateHandler
// @Summary Görev ekle
// @Description Yeni görev oluşturur
// @Tags Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param task body models.Task true "Görev"
// @Success 201 {object} models.Task
// @Failure 400 {object} map[string]string
// @Router /tasks [post]
func TaskCreateHandler(c *fiber.Ctx) error {
	var task models.Task
	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz veri"})
	}
	if task.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Başlık zorunlu"})
	}
	uid := c.Locals("user_id")
	userID, ok := uid.(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Kullanıcı bilgisi alınamadı"})
	}
	task.ID = len(models.Tasks) + 1
	task.UserID = userID
	models.Tasks = append(models.Tasks, task)
	return c.Status(fiber.StatusCreated).JSON(task)
}

// TaskDetailHandler görev detayını döner
// @ID TaskDetailHandler
// @Summary Görev detayını görüntüle
// @Description Belirli bir görevin detayını döner
// @Tags Tasks
// @Produce json
// @Security BearerAuth
// @Param id path int true "Görev ID"
// @Success 200 {object} models.Task
// @Failure 404 {object} map[string]string
// @Router /tasks/{id} [get]
func TaskDetailHandler(c *fiber.Ctx) error {
	uid := c.Locals("user_id")
	userID, ok := uid.(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Kullanıcı bilgisi alınamadı"})
	}
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz görev ID"})
	}
	for _, t := range models.Tasks {
		if t.ID == id && t.UserID == userID {
			return c.JSON(t)
		}
	}
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Görev bulunamadı veya yetkiniz yok"})
}

// TaskUpdateHandler görevi günceller
// @ID TaskUpdateHandler
// @Summary Görev güncelle
// @Description Belirli bir görevi günceller
// @Tags Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Görev ID"
// @Param task body models.Task true "Görev"
// @Success 200 {object} models.Task
// @Failure 404 {object} map[string]string
// @Router /tasks/{id} [put]
func TaskUpdateHandler(c *fiber.Ctx) error {
	uid := c.Locals("user_id")
	userID, ok := uid.(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Kullanıcı bilgisi alınamadı"})
	}
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz görev ID"})
	}
	var updated models.Task
	if err := c.BodyParser(&updated); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz veri"})
	}
	for i, t := range models.Tasks {
		if t.ID == id && t.UserID == userID {
			models.Tasks[i].Title = updated.Title
			models.Tasks[i].Details = updated.Details
			return c.JSON(models.Tasks[i])
		}
	}
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Görev bulunamadı veya yetkiniz yok"})
}

// TaskDeleteHandler görevi siler
// @ID TaskDeleteHandler
// @Summary Görev sil
// @Description Belirli bir görevi siler
// @Tags Tasks
// @Security BearerAuth
// @Param id path int true "Görev ID"
// @Success 204 {string} string ""
// @Failure 404 {object} map[string]string
// @Router /tasks/{id} [delete]
func TaskDeleteHandler(c *fiber.Ctx) error {
	uid := c.Locals("user_id")
	userID, ok := uid.(int)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Kullanıcı bilgisi alınamadı"})
	}
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz görev ID"})
	}
	for i, t := range models.Tasks {
		if t.ID == id && t.UserID == userID {
			models.Tasks = append(models.Tasks[:i], models.Tasks[i+1:]...)
			return c.SendStatus(fiber.StatusNoContent)
		}
	}
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Görev bulunamadı veya yetkiniz yok"})
}
