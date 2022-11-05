package test

import (
	"database/sql/driver"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

type query struct {
	Raw  string
	Args []driver.Value
	Rows *sqlmock.Rows
}

type pathParam struct {
	Names  []string
	Values []string
}

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type AnyString struct{}

func (a AnyString) Match(v driver.Value) bool {
	_, ok := v.(string)
	return ok
}
