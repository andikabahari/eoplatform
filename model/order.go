package model

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	IsAccepted  bool
	IsCompleted bool
	UserID      uint
	User        User
	Services    []Service `gorm:"many2many:order_services;"`
}
