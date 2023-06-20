package models

import (
	"gorm.io/gorm"
)

type Shopping struct {
	gorm.Model
	UserID   uint   `gorm:"not null"`
	ItemName string `gorm:"not null"`
	Quantity string `gorm:"not null"`
	Note     string
	IsBought bool `gorm:"default:false"`
}
