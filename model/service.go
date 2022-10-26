package model

import "gorm.io/gorm"

type Service struct {
	gorm.Model
	UserID      uint
	Name        string
	Description string
	Cost        float64
}
