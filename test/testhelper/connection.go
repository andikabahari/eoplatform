package testhelper

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Init(conn *sql.DB) *gorm.DB {
	db, err := gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      conn,
	}))
	if err != nil {
		panic(err)
	}

	return db
}

func Mock() (*sql.DB, sqlmock.Sqlmock) {
	conn, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	return conn, mock
}
