package model

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	IsAccepted  bool
	IsCompleted bool
	DateOfEvent time.Time
	FirstName   string
	LastName    string
	Phone       string
	Email       string
	Address     string
	Note        string
	UserID      uint
	User        User
	Services    []Service `gorm:"many2many:order_services;"`
}
