package handler

import (
	"fmt"
	"log"
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

func (h *OrderHandler) GetOrders(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*helper.JWTCustomClaims)

	orders := make([]model.Order, 0)
	orderRepository := repository.NewOrderRepository(h.server.DB)
	if claims.Role == "customer" {
		orderRepository.GetOrdersForCustomer(&orders, claims.ID)
	}
	if claims.Role == "organizer" {
		orderRepository.GetOrdersForOrganizer(&orders, claims.ID)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": response.NewOrdersResponse(orders),
	})
}

func (h *OrderHandler) CreateOrder(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*helper.JWTCustomClaims)

	if claims.Role != "customer" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "unauthorized",
		})
	}

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

	if order.Services[0].UserID != claims.ID {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "unauthorized",
		})
	}

	segment := strings.Split(c.Path(), "/")[4]
	if segment == "accept" && !order.IsAccepted {
		order.IsAccepted = true

		var totalCost float64
		for _, service := range order.Services {
			totalCost += service.Cost
		}

		payment := model.Payment{}
		payment.OrderID = order.ID
		payment.Amount = totalCost
		payment.Status = "pending"
		h.server.DB.Debug().Omit("Order").Save(&payment)

		bankAccount := model.BankAccount{}
		h.server.DB.Debug().Where("user_id = ?", claims.ID).First(&bankAccount)

		helper.ChargeOrder(map[string]any{
			"payment_type": "bank_transfer",
			"transaction_details": map[string]any{
				"order_id":     order.ID,
				"gross_amount": totalCost,
			},
			"bank_transfer": map[string]any{
				"bank":      bankAccount.Bank,
				"va_number": bankAccount.VANumber,
			},
			"customer_details": map[string]any{
				"phone":   order.Phone,
				"email":   order.Email,
				"address": order.Address,
			},
		})
	}
	if segment == "complete" && order.IsAccepted {
		order.IsCompleted = true
	}

	h.server.DB.Debug().Omit(clause.Associations).Save(order)

	return c.JSON(http.StatusOK, echo.Map{
		"data": response.NewOrderResponse(order),
	})
}

func (h *OrderHandler) PaymentStatus(c echo.Context) error {
	req := request.MidtransTransactionNotificationRequest{}
	log.Println("Midtrans request:", req)

	if err := c.Bind(&req); err != nil {
		return err
	}

	payment := model.Payment{}
	h.server.DB.Debug().Where("order_id = ?", req.OrderID).Find(&payment)

	if payment.OrderID == 0 {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "order not found",
		})
	}

	switch req.Status {
	case "settlement":
	case "capture":
		payment.Status = "success"
	case "deny":
	case "cancel":
	case "expire":
		payment.Status = "fail"
	}

	h.server.DB.Debug().Omit("Order").Save(&payment)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "order status updated",
	})
}
