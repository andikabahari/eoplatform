package testhelper

import (
	"database/sql"

	"github.com/andikabahari/eoplatform/config"
	s "github.com/andikabahari/eoplatform/server"
	"github.com/labstack/echo/v4"
)

func NewServer(conn *sql.DB) *s.Server {
	return &s.Server{
		Echo:   echo.New(),
		DB:     Init(conn),
		Config: config.NewConfig(),
	}
}
