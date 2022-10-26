package route

import (
	"github.com/andikabahari/eoplatform/helper"
	s "github.com/andikabahari/eoplatform/server"
	"github.com/andikabahari/eoplatform/server/handler"
	"github.com/labstack/echo/v4/middleware"
)

func Setup(server *s.Server) {
	server.Echo.Use(middleware.Recover())
	server.Echo.Use(middleware.Logger())

	helloHandler := handler.NewHelloHandler(server)
	server.Echo.GET("/", helloHandler.Greeting)

	registerHandler := handler.NewRegisterHandler(server)
	server.Echo.POST("/v1/register", registerHandler.Register)

	loginHandler := handler.NewLoginHandler(server)
	server.Echo.POST("/v1/login", loginHandler.Login)

	restricted := server.Echo.Group("")
	restricted.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &helper.JWTCustomClaims{},
		SigningKey: []byte(server.Config.Auth.Secret),
	}))

	serviceHandler := handler.NewServiceHandler(server)
	restricted.POST("/v1/services", serviceHandler.CreateService)
}
