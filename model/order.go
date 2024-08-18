package model

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	UserID uint        `gorm:"not null"`
	Total  float64     `gorm:"not null" json:"total"`
	Items  []OrderItem `gorm:"foreignKey:OrderID" json:"items"`
}

type OrderItem struct {
	gorm.Model
	OrderID   uint    `gorm:"not null" json:"order_id"`
	ProductID uint    `gorm:"not null" json:"product_id"`
	Quantity  int     `gorm:"not null" json:"quantity"`
	Price     float64 `gorm:"not null" json:"price"`
	Note      string  `gorm:"size:255" json:"note"`
}
