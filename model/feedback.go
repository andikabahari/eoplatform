package model

import "gorm.io/gorm"

type Feedback struct {
	gorm.Model
	Description string
	Rating      uint
	Positive    float64
	Negative    float64
	FromUserID  uint
	FromUser    User
	ToUserID    uint
	ToUser      User
}
