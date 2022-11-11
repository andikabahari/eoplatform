package usecase

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/andikabahari/eoplatform/helper"
	"github.com/andikabahari/eoplatform/model"
	r "github.com/andikabahari/eoplatform/repository"
	"github.com/andikabahari/eoplatform/request"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type OrderUsecase interface {
	GetOrders(claims *helper.JWTCustomClaims, orders *[]model.Order, payments *[]model.Payment)
	CreateOrder(claims *helper.JWTCustomClaims, order *model.Order, req *request.CreateOrderRequest) helper.APIError
	AcceptOrCompleteOrder(ctx echo.Context, order *model.Order) helper.APIError
	CancelOrder(ctx echo.Context, order *model.Order) helper.APIError
	PaymentStatus(req *request.MidtransTransactionNotificationRequest) helper.APIError
}

type orderUsecase struct {
	orderRepository       r.OrderRepository
	paymentRepository     r.PaymentRepository
	userRepository        r.UserRepository
	serviceRepository     r.ServiceRepository
	bankAccountRepository r.BankAccountRepository
}

func NewOrderUsecase(
	orderRepository r.OrderRepository,
	paymentRepository r.PaymentRepository,
	userRepository r.UserRepository,
	serviceRepository r.ServiceRepository,
	bankAccountRepository r.BankAccountRepository,
) OrderUsecase {
	return &orderUsecase{
		orderRepository,
		paymentRepository,
		userRepository,
		serviceRepository,
		bankAccountRepository,
	}
}

func (u *orderUsecase) GetOrders(claims *helper.JWTCustomClaims, orders *[]model.Order, payments *[]model.Payment) {
	if claims.Role == "customer" {
		u.orderRepository.GetOrdersForCustomer(orders, claims.ID)
	}
	if claims.Role == "organizer" {
		u.orderRepository.GetOrdersForOrganizer(orders, claims.ID)
	}

	tmpPayments := make([]model.Payment, len(*orders))
	for i, order := range *orders {
		payment := model.Payment{}
		u.paymentRepository.FindOnlyByOrderID(&payment, order.ID)
		tmpPayments[i] = payment
	}
	*payments = tmpPayments
}

func (u *orderUsecase) CreateOrder(claims *helper.JWTCustomClaims, order *model.Order, req *request.CreateOrderRequest) helper.APIError {
	dateOfEvent, err := time.Parse("2006-01-02", req.DateOfEvent)
	if err != nil {
		log.Printf("Error: %s", err)
		return helper.NewAPIError(http.StatusInternalServerError, "internal server error")
	}

	services := make([]model.Service, 0)

	first := model.Service{}
	for i, id := range req.ServiceIDs {
		service := model.Service{}
		u.serviceRepository.Find(&service, fmt.Sprintf("%d", id))

		if i == 0 {
			first = service
		}

		if service.ID == 0 || service.UserID != first.UserID {
			return helper.NewAPIError(http.StatusBadRequest, "cannot proceed your order")
		}

		services = append(services, service)
	}

	user := model.User{}
	u.userRepository.Find(&user, claims.ID)

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

	u.orderRepository.Create(order)

	order.User = user

	return nil
}

func (u *orderUsecase) AcceptOrCompleteOrder(ctx echo.Context, order *model.Order) helper.APIError {
	u.orderRepository.Find(order, ctx.Param("id"))

	if order.ID == 0 {
		return helper.NewAPIError(http.StatusNotFound, "order not found")
	}

	userToken := ctx.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*helper.JWTCustomClaims)

	if order.Services[0].UserID != claims.ID {
		return helper.NewAPIError(http.StatusUnauthorized, "unauthorized")
	}

	segment := strings.Split(ctx.Path(), "/")[4]
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
		u.paymentRepository.Create(&payment)

		bankAccount := model.BankAccount{}
		u.bankAccountRepository.FindByUserID(&bankAccount, claims.ID)

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

	u.orderRepository.Save(order)

	return nil
}

func (u *orderUsecase) CancelOrder(ctx echo.Context, order *model.Order) helper.APIError {
	u.orderRepository.Find(order, ctx.Param("id"))

	if order.ID == 0 {
		return helper.NewAPIError(http.StatusNotFound, "order not found")
	}

	userToken := ctx.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*helper.JWTCustomClaims)

	if order.UserID != claims.ID {
		return helper.NewAPIError(http.StatusUnauthorized, "unauthorized")
	}

	u.orderRepository.Delete(order)

	return nil
}

func (u *orderUsecase) PaymentStatus(req *request.MidtransTransactionNotificationRequest) helper.APIError {
	orderID := strings.Split(req.OrderID, "-")[1]

	payment := model.Payment{}
	u.paymentRepository.FindOnlyByOrderID(&payment, orderID)

	if payment.OrderID == 0 {
		return helper.NewAPIError(http.StatusNotFound, "order not found")
	}

	if req.Status == "settlement" || req.Status == "capture" {
		req.Status = "success"
	}
	if req.Status == "deny" || req.Status == "cancel" || req.Status == "expire" {
		req.Status = "fail"

		order := model.Order{}
		u.orderRepository.FindOnly(&order, orderID)
		u.orderRepository.Delete(&order)
	}
	u.paymentRepository.Update(&payment, req)

	return nil
}
