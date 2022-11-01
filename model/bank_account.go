package model

import "gorm.io/gorm"

type BankAccount struct {
	gorm.Model
	Bank     string
	VANumber string
	UserID   uint
	User     User
}
