package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

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
	server                *s.Server
	orderRepository       *repository.OrderRepository
	paymentRepository     *repository.PaymentRepository
	userRepository        *repository.UserRepository
	serviceRepository     *repository.ServiceRepository
	bankAccountRepository *repository.BankAccountRepository
}

func NewOrderHandler(server *s.Server) *OrderHandler {
	return &OrderHandler{
		server,
		repository.NewOrderRepository(server.DB),
		repository.NewPaymentRepository(server.DB),
		repository.NewUserRepository(server.DB),
		repository.NewServiceRepository(server.DB),
		repository.NewBankAccountRepository(server.DB),
	}
}

func (h *OrderHandler) GetOrders(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*helper.JWTCustomClaims)

	orders := make([]model.Order, 0)
	if claims.Role == "customer" {
		h.orderRepository.GetOrdersForCustomer(&orders, claims.ID)
	}
	if claims.Role == "organizer" {
		h.orderRepository.GetOrdersForOrganizer(&orders, claims.ID)
	}

	payments := make([]model.Payment, len(orders))
	for i, order := range orders {
		payment := model.Payment{}
		h.paymentRepository.FindOnlyByOrderID(&payment, order.ID)
		payments[i] = payment
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "fetch orders successful",
		"data":    response.NewOrdersWithPaymentStatusResponse(orders, payments),
	})
}

func (h *OrderHandler) CreateOrder(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*helper.JWTCustomClaims)

	if claims.Role != "customer" {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "create order failure",
			"error":   "unauthorized",
		})
	}

	req := request.CreateOrderRequest{}

	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "validation error",
			"error":   err,
		})
	}

	dateOfEvent, err := time.Parse("2006-01-02", req.DateOfEvent)
	if err != nil {
		return err
	}

	services := make([]model.Service, 0)

	first := model.Service{}
	for i, id := range req.ServiceIDs {
		service := model.Service{}
		h.serviceRepository.Find(&service, fmt.Sprintf("%d", id))

		if i == 0 {
			first = service
		}

		if service.ID == 0 || service.UserID != first.UserID {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "cannot proceed your order",
			})
		}

		services = append(services, service)
	}

	user := model.User{}
	h.userRepository.Find(&user, claims.ID)

	order := model.Order{}
	order.IsAccepted = false
	order.IsCompleted = false
	order.DateOfEvent = dateOfEvent
	order.FirstName = req.FirstName
	order.LastName = req.LastName
	order.Phone = req.Phone
	order.Email = req.Email
	order.Address = req.Address
	order.Note = req.Note
	order.UserID = claims.ID
	order.Services = services

	h.orderRepository.Create(&order)

	order.User = user

	return c.JSON(http.StatusOK, echo.Map{
		"message": "create order successful",
		"data":    response.NewOrderResponse(order),
	})
}

func (h *OrderHandler) AcceptOrCompleteOrder(c echo.Context) error {
	order := model.Order{}
	h.orderRepository.Find(&order, c.Param("id"))

	if order.ID == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "accept or complete order failure",
			"error":   "order not found",
		})
	}

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*helper.JWTCustomClaims)

	if order.Services[0].UserID != claims.ID {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "accept or complete order failure",
			"error":   "unauthorized",
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
		h.paymentRepository.Create(&payment)

		bankAccount := model.BankAccount{}
		h.bankAccountRepository.FindByUserID(&bankAccount, claims.ID)

		transaction := map[string]any{
			"payment_type": "bank_transfer",
			"transaction_details": map[string]any{
				"order_id":     fmt.Sprintf("EOP-%d", order.ID),
				"gross_amount": totalCost,
			},
			"bank_transfer": map[string]any{
				"bank":      bankAccount.Bank,
				"va_number": bankAccount.VANumber,
			},
			"customer_details": map[string]any{
				"first_name": order.FirstName,
				"last_name":  order.LastName,
				"phone":      order.Phone,
				"email":      order.Email,
				"address":    order.Address,
			},
		}
		helper.ChargeOrder(transaction)
	}
	if segment == "complete" && order.IsAccepted {
		order.IsCompleted = true
	}

	h.server.DB.Debug().Omit(clause.Associations).Save(order)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "accept or complete order successful",
		"data":    response.NewOrderResponse(order),
	})
}

func (h *OrderHandler) CancelOrder(c echo.Context) error {
	order := model.Order{}
	h.orderRepository.Find(&order, c.Param("id"))

	if order.ID == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "cancel order failure",
			"error":   "order not found",
		})
	}

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*helper.JWTCustomClaims)

	if order.UserID != claims.ID {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "cancel service failure",
			"error":   "unauthorized",
		})
	}

	h.orderRepository.Delete(&order)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "cancel order successful",
		"data": echo.Map{
			"kind":    "order",
			"id":      c.Param("id"),
			"deleted": true,
		},
	})
}

func (h *OrderHandler) PaymentStatus(c echo.Context) error {
	req := request.MidtransTransactionNotificationRequest{}

	if err := c.Bind(&req); err != nil {
		return err
	}
	log.Println("Midtrans request:", req)

	orderID := strings.Split(req.OrderID, "-")[1]

	payment := model.Payment{}
	h.paymentRepository.FindOnlyByOrderID(&payment, orderID)

	if payment.OrderID == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "payment failure",
			"error":   "order not found",
		})
	}

	if req.Status == "settlement" || req.Status == "capture" {
		req.Status = "success"
	}
	if req.Status == "deny" || req.Status == "cancel" || req.Status == "expire" {
		req.Status = "fail"

		order := model.Order{}
		h.orderRepository.FindOnly(&order, orderID)
		h.orderRepository.Delete(&order)
	}
	h.paymentRepository.Update(&payment, &req)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "payment successful",
	})
}
