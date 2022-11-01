package model

import "gorm.io/gorm"

type BankAccounts struct {
	gorm.Model
	Bank     string
	VANumber string
	UserID   uint
	User     User
}
