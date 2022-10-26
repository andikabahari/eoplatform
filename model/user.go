package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string
	Username string `gorm:"index:,unique"`
	Password string
	Email    string
	Address  string
}
