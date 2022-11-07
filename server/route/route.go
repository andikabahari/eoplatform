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

	accountV1 := v1.Group("/account")
	accountHandler := handler.NewAccountHandler(server)
	accountV1.GET("", accountHandler.GetAccount, auth)
	accountV1.PUT("", accountHandler.UpdateAccount, auth)
	accountV1.PUT("/password", accountHandler.ResetPassword, auth)

	serviceV1 := v1.Group("/services")
	serviceHandler := handler.NewServiceHandler(server)
	serviceV1.GET("", serviceHandler.GetServices)
	serviceV1.GET("/:id", serviceHandler.FindService)
	serviceV1.POST("", serviceHandler.CreateService, auth)
	serviceV1.PUT("/:id", serviceHandler.UpdateService, auth)
	serviceV1.DELETE("/:id", serviceHandler.DeleteService, auth)

	orderV1 := v1.Group("/orders")
	orderHandler := handler.NewOrderHandler(server)
	orderV1.GET("", orderHandler.GetOrders, auth)
	orderV1.POST("", orderHandler.CreateOrder, auth)
	orderV1.POST("/:id/accept", orderHandler.AcceptOrCompleteOrder, auth)
	orderV1.POST("/:id/complete", orderHandler.AcceptOrCompleteOrder, auth)
	orderV1.POST("/:id/cancel", orderHandler.CancelOrder, auth)
	v1.POST("/MDDRlkYVFm9QOLK08MDp", orderHandler.PaymentStatus)

	bankAccountV1 := v1.Group("/bank-accounts")
	bankAccountHandler := handler.NewBankAccountHandler(server)
	bankAccountV1.GET("", bankAccountHandler.GetBankAccounts, auth)
	bankAccountV1.POST("", bankAccountHandler.CreateBankAccount, auth)
	bankAccountV1.PUT("", bankAccountHandler.UpdateBankAccount, auth)

	feedbackV1 := v1.Group("/feedbacks")
	feedbackHandler := handler.NewFeedbackHandler(server)
	feedbackV1.GET("", feedbackHandler.GetFeedbacks)
	feedbackV1.POST("", feedbackHandler.CreateFeedback, auth)
}
