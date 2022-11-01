package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/andikabahari/eoplatform/helper"
	"github.com/andikabahari/eoplatform/model"
	"github.com/andikabahari/eoplatform/repository"
	"github.com/andikabahari/eoplatform/request"
	"github.com/andikabahari/eoplatform/response"
	s "github.com/andikabahari/eoplatform/server"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm/clause"
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
	userRepository.Find(&user, claims.ID)

	order := model.Order{}
	order.IsAccepted = false
	order.IsCompleted = false
	order.Phone = req.Phone
	order.Email = req.Email
	order.Address = req.Address
	order.Note = req.Note
	order.UserID = claims.ID
	order.Services = services

	orderRepository := repository.NewOrderRepository(h.server.DB)
	orderRepository.Create(&order)

	order.User = user

	return c.JSON(http.StatusOK, echo.Map{
		"data": response.NewOrderResponse(order),
	})
}

func (h *OrderHandler) AcceptOrCompleteOrder(c echo.Context) error {
	order := model.Order{}
	orderRepository := repository.NewOrderRepository(h.server.DB)
	orderRepository.Find(&order, c.Param("id"))

	if order.ID == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "order not found",
		})
	}

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*helper.JWTCustomClaims)

	if order.UserID != claims.ID {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "unauthorized",
		})
	}

	message := ""
	segment := strings.Split(c.Path(), "/")[4]
	if segment == "accept" {
		order.IsAccepted = true
		message = "Your request has been accepted, please contact us as soon as possible."
	}
	if segment == "complete" && order.IsAccepted {
		order.IsCompleted = true
		message = "Your order has been completed."
	}

	if message != "" {
		if err := helper.SendEmail([]string{order.Email}, message); err != nil {
			return err
		}
	}

	h.server.DB.Debug().Omit(clause.Associations).Save(order)

	return c.JSON(http.StatusOK, echo.Map{
		"data": response.NewOrderResponse(order),
	})
}
