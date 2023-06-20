package models

import (
	"time"

	"gorm.io/gorm"
)

type Refrigerator struct {
	gorm.Model
	UserID           uint   `gorm:"not null"`
	RefrigeratorName string `gorm:"not null"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        gorm.DeletedAt `gorm:"index"`
}
