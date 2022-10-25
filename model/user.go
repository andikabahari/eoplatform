package model

import (
	"database/sql/driver"

	"gorm.io/gorm"
)

type role string

const (
	ORGANIZER role = "organizer"
	CUSTOMER  role = "customer"
)

func (r *role) Scan(value any) error {
	*r = role(value.([]byte))
	return nil
}

func (r role) Value() (driver.Value, error) {
	return string(r), nil
}

type User struct {
	gorm.Model
	Name     string
	Email    string
	Password string
	Role     role `sql:"type:ENUM('organizer', 'customer')"`
}
