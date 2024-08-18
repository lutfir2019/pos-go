package model

import "gorm.io/gorm"

// Product struct
type Product struct {
	gorm.Model
	Code          string `gorm:"uniqueIndex;not null;size:255" json:"product_code"`
	Name          string `gorm:"not null" json:"name"`
	Quantity      uint   `gorm:"not null" json:"quantity"`
	PricePurchase uint   `gorm:"not null" json:"price_purchase"`
	PriceSelling  uint   `gorm:"not null" json:"price_selling"`
	Image         string `gorm:"not null" json:"file"`
	ReferUser     uint   `json:"-"`
}

type ProductCounter struct {
	Date    string `gorm:"primaryKey"`
	Counter int
}
