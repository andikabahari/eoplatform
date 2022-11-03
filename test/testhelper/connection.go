package testhelper

import (
	"database/sql"
	"database/sql/driver"
	"log"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

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

func Init(conn *sql.DB) *gorm.DB {
	db, err := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      conn,
	}))
	if err != nil {
		log.Fatal("Can't connect to DB!")
	}

	return db
}

func Mock() (*sql.DB, sqlmock.Sqlmock) {
	conn, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal("Can't mock DB!")
	}

	return conn, mock
}
