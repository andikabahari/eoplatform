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

	accountHandler := handler.NewAccountHandler(server)
	restricted.GET("/v1/account", accountHandler.GetAccount)
	restricted.PUT("/v1/account", accountHandler.UpdateAccount)
	restricted.PUT("/v1/account/password", accountHandler.ResetPassword)
	restricted.GET("/v1/account/orders", accountHandler.GetOrders)

	serviceHandler := handler.NewServiceHandler(server)
	server.Echo.GET("/v1/services", serviceHandler.GetServices)
	server.Echo.GET("/v1/services/:id", serviceHandler.FindService)
	restricted.POST("/v1/services", serviceHandler.CreateService)
	restricted.PUT("/v1/services/:id", serviceHandler.UpdateService)
	restricted.DELETE("/v1/services/:id", serviceHandler.DeleteService)

	orderHandler := handler.NewOrderHandler(server)
	restricted.POST("/v1/orders", orderHandler.CreateOrder)
	restricted.GET("/v1/orders/:id/accept", orderHandler.AcceptOrCompleteOrder)
	restricted.GET("/v1/orders/:id/complete", orderHandler.AcceptOrCompleteOrder)
}
