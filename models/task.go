package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description"`
	Status      string         `json:"status" gorm:"default:pending"`  // pending, in_progress, completed
	Priority    string         `json:"priority" gorm:"default:medium"` // low, medium, high
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"` // Soft delete
	User        User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// In-memory storage for backward compatibility (will be removed after DB migration)
var PublicTasks = []Task{
	{ID: 1, UserID: 0, Title: "Örnek Görev 1", Description: "Bu public bir görevdir."},
	{ID: 2, UserID: 0, Title: "Örnek Görev 2", Description: "Herkes görebilir."},
}

var Tasks = []Task{}
