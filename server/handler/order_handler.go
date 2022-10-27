package handler

import (
	"fmt"
	"net/http"

	"github.com/andikabahari/eoplatform/helper"
	"github.com/andikabahari/eoplatform/model"
	"github.com/andikabahari/eoplatform/repository"
	"github.com/andikabahari/eoplatform/request"
	"github.com/andikabahari/eoplatform/response"
	s "github.com/andikabahari/eoplatform/server"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type OrderHandler struct {
	server *s.Server
}

func NewOrderHandler(server *s.Server) *OrderHandler {
	return &OrderHandler{server}
}

func (h *OrderHandler) CreateOrder(c echo.Context) error {
	req := request.CreateOrderRequest{}

	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err,
		})
	}

	serviceRepository := repository.NewServiceRepository(h.server.DB)
	services := make([]model.Service, 0)

	first := model.Service{}
	for i, id := range req.ServiceIDs {
		service := model.Service{}
		serviceRepository.Find(&service, fmt.Sprintf("%d", id))

		if i == 0 {
			first = service
		}

		if service.ID == 0 || service.UserID != first.UserID {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error": "cannot proceed your order",
			})
		}

		services = append(services, service)
	}

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*helper.JWTCustomClaims)

	user := model.User{}
	userRepository := repository.NewUserRepository(h.server.DB)
	userRepository.Find(&user, fmt.Sprintf("%d", claims.ID))

	order := model.Order{}
	order.IsAccepted = false
	order.IsCompleted = false
	order.UserID = claims.ID
	order.Services = services

	orderRepository := repository.NewOrderRepository(h.server.DB)
	orderRepository.Create(&order)

	order.User = user

	return c.JSON(http.StatusOK, echo.Map{
		"data": response.NewOrderResponse(order),
	})
}
