package models

import (
	"time"

	"gorm.io/gorm"
)

type Ingredient struct {
	gorm.Model
	RefrigeratorID uint      `gorm:"not null"`
	UserID         uint      `gorm:"not null"`
	IngredientName string    `gorm:"not null"`
	Quantity       string    `gorm:"not null"`
	ValidUntil     time.Time `gorm:"not null"`
	CategoryID     string    `gorm:"not null"`
}
