package model

import "gorm.io/gorm"

// Product struct
type Product struct {
	gorm.Model
	Code     string  `gorm:"uniqueIndex;not null;size:255" json:"code"`
	Name     string  `gorm:"not null" json:"name"`
	Quantity uint    `gorm:"not null" json:"quantity"`
	Price    float64 `gorm:"not null" json:"price"`
}
