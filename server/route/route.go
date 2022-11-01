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

	auth := middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &helper.JWTCustomClaims{},
		SigningKey: []byte(server.Config.Auth.Secret),
	})

	helloHandler := handler.NewHelloHandler(server)
	server.Echo.GET("/", helloHandler.Greeting)

	v1 := server.Echo.Group("/v1")

	registerHandler := handler.NewRegisterHandler(server)
	v1.POST("/register", registerHandler.Register)

	loginHandler := handler.NewLoginHandler(server)
	v1.POST("/login", loginHandler.Login)

	accountHandler := handler.NewAccountHandler(server)
	v1.GET("/account", accountHandler.GetAccount, auth)
	v1.PUT("/account", accountHandler.UpdateAccount, auth)
	v1.PUT("/account/password", accountHandler.ResetPassword, auth)

	serviceHandler := handler.NewServiceHandler(server)
	v1.GET("/services", serviceHandler.GetServices)
	v1.GET("/services/:id", serviceHandler.FindService)
	v1.POST("/services", serviceHandler.CreateService, auth)
	v1.PUT("/services/:id", serviceHandler.UpdateService, auth)
	v1.DELETE("/services/:id", serviceHandler.DeleteService, auth)

	orderHandler := handler.NewOrderHandler(server)
	v1.GET("/orders", orderHandler.GetOrders, auth)
	v1.POST("/orders", orderHandler.CreateOrder, auth)
	v1.POST("/orders/:id/accept", orderHandler.AcceptOrCompleteOrder, auth)
	v1.POST("/orders/:id/complete", orderHandler.AcceptOrCompleteOrder, auth)
	v1.POST("/MDDRlkYVFm9QOLK08MDp", orderHandler.PaymentStatus)

	bankAccountHandler := handler.NewBankAccountHandler(server)
	v1.POST("/bank-accounts", bankAccountHandler.CreateBankAccount, auth)

	feedbackHandler := handler.NewFeedbackHandler(server)
	v1.POST("/feedbacks", feedbackHandler.CreateFeedback, auth)
}
