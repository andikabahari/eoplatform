package route

import (
	s "github.com/andikabahari/eoplatform/server"
	"github.com/andikabahari/eoplatform/server/handler"
	"github.com/labstack/echo/v4/middleware"
)

func Setup(server *s.Server) {
	server.Echo.Use(middleware.Recover())
	server.Echo.Use(middleware.Logger())

	helloHandler := handler.NewHelloHandler(server)
	server.Echo.GET("/", helloHandler.Greeting)
}
