package model

import "gorm.io/gorm"

type Service struct {
	gorm.Model
	UserID      uint
	User        User
	Name        string
	Cost        float64
	Phone       string
	Email       string
	IsPublished bool
	Description string
}
