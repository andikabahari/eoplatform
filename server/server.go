package server

import (
	"github.com/andikabahari/eoplatform/config"
	"github.com/andikabahari/eoplatform/db"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Server struct {
	Echo   *echo.Echo
	DB     *gorm.DB
	Config *config.Config
}

func NewServer(config *config.Config) *Server {
	return &Server{
		Echo:   echo.New(),
		DB:     db.Init(config),
		Config: config,
	}
}

func (server *Server) Run() {
	server.Echo.Logger.Fatal(server.Echo.Start(":" + server.Config.HTTP.Port))
}
