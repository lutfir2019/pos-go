package model

import "gorm.io/gorm"

// User struct
type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null;size:255" json:"username"`
	Email    string `gorm:"uniqueIndex;not null;size:255" json:"email"`
	Password string `gorm:"not null" json:"-"`
	Name     string `gorm:"size:255" json:"name"`
}
