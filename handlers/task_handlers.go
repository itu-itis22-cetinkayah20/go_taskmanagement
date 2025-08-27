package handlers

import (
	"strconv"

	"go_taskmanagement/database"
	"go_taskmanagement/models"

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
	if database.IsConnected && database.DB != nil {
		var publicTasks []models.Task
		database.DB.Preload("User").Where("user_id = ?", 0).Find(&publicTasks)
		return c.JSON(publicTasks)
	}

	// In-memory mode (fallback)
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
	userID, ok := uid.(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Kullanıcı bilgisi alınamadı"})
	}

	if database.IsConnected && database.DB != nil {
		var userTasks []models.Task
		database.DB.Preload("User").Where("user_id = ?", userID).Find(&userTasks)
		return c.JSON(userTasks)
	}

	// In-memory mode (fallback)
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
	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
		Priority    string `json:"priority"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz veri"})
	}

	if input.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Başlık zorunlu"})
	}

	uid := c.Locals("user_id")
	userID, ok := uid.(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Kullanıcı bilgisi alınamadı"})
	}

	// Set defaults
	if input.Status == "" {
		input.Status = "pending"
	}
	if input.Priority == "" {
		input.Priority = "medium"
	}

	task := models.Task{
		UserID:      userID,
		Title:       input.Title,
		Description: input.Description,
		Status:      input.Status,
		Priority:    input.Priority,
	}

	if database.IsConnected && database.DB != nil {
		if err := database.DB.Create(&task).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Görev oluşturulamadı"})
		}
		// Preload user information for the created task
		database.DB.Preload("User").First(&task, task.ID)
	} else {
		// In-memory mode (fallback)
		task.ID = uint(len(models.Tasks) + 1)
		models.Tasks = append(models.Tasks, task)
	}

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
	userID, ok := uid.(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Kullanıcı bilgisi alınamadı"})
	}

	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz görev ID"})
	}

	var task models.Task

	if database.DB == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database connection not available"})
	}

	if err := database.DB.Preload("User").Where("id = ? AND user_id = ?", uint(id), userID).First(&task).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Görev bulunamadı veya yetkiniz yok"})
	}

	return c.JSON(task)
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
	userID, ok := uid.(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Kullanıcı bilgisi alınamadı"})
	}

	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz görev ID"})
	}

	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Status      string `json:"status"`
		Priority    string `json:"priority"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz veri"})
	}

	// Validate title if provided
	if input.Title == "" && input.Description == "" && input.Status == "" && input.Priority == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "En az bir alan güncellenmelidir"})
	}

	// Find the task
	var task models.Task

	if database.DB == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database connection not available"})
	}

	if err := database.DB.Where("id = ? AND user_id = ?", uint(id), userID).First(&task).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Görev bulunamadı veya yetkiniz yok"})
	}

	// Update fields
	updates := make(map[string]interface{})
	if input.Title != "" {
		updates["title"] = input.Title
	}
	if input.Description != "" {
		updates["description"] = input.Description
	}
	if input.Status != "" {
		updates["status"] = input.Status
	}
	if input.Priority != "" {
		updates["priority"] = input.Priority
	}

	if err := database.DB.Model(&task).Updates(updates).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Görev güncellenemedi"})
	}

	// Reload the task with user information
	database.DB.Preload("User").First(&task, task.ID)

	return c.JSON(task)
}

// TaskDeleteHandler görevi siler
// @ID TaskDeleteHandler
// @Summary Görev sil
// @Description Belirli bir görevi siler
// @Tags Tasks
// @Security BearerAuth
// @Param id path int true "Görev ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tasks/{id} [delete]
func TaskDeleteHandler(c *fiber.Ctx) error {
	uid := c.Locals("user_id")
	userID, ok := uid.(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Kullanıcı bilgisi alınamadı"})
	}

	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Geçersiz görev ID"})
	}

	// Check database connection
	if database.DB == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database connection not available"})
	}

	// Soft delete the task
	result := database.DB.Where("id = ? AND user_id = ?", uint(id), userID).Delete(&models.Task{})
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Görev silinemedi"})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Görev bulunamadı veya yetkiniz yok"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Task deleted successfully"})
}
