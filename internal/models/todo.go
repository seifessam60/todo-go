package models

import (
	"time"

	"gorm.io/gorm"
)

type Todo struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	UserID      uint   `json:"user_id" gorm:"not null;index"`
	Title       string `json:"title" gorm:"not null"`
	Description string `json:"description"`
	Priority    string `json:"priority" gorm:"default:'medium'"`
	Category    string `json:"category"`
	DueDate     *time.Time `json:"due_date"`
	Completed 	bool   `json:"completed" gorm:"default:false"`
	CreatedAt 	time.Time `json:"created_at"`
	UpdatedAt 	time.Time `json:"updated_at"`
	DeletedAt	gorm.DeletedAt `json:"-" gorm:"index"`

	User User `json:"-" gorm:"foreignKey:UserID"`
}
