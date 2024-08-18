package model

import "gorm.io/gorm"

// User struct
type User struct {
	gorm.Model
	Username string    `gorm:"uniqueIndex;not null;size:255" json:"username"`
	Password string    `gorm:"not null" json:"-"`
	Name     string    `gorm:"size:255" json:"name"`
	Role     string    `gorm:"size:255" json:"role"`
	Product  []Product `gorm:"foreignKey:ReferUser;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}
