package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Username string `gorm:"not null" json:"username"`
	Email    string `gorm:"not null" json:"email"`
	Tasks    []Task `json:"tasks"`
}
