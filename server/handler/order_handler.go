package handler

import (
	"log"
	"net/http"

	"github.com/andikabahari/eoplatform/helper"
	"github.com/andikabahari/eoplatform/model"
	"github.com/andikabahari/eoplatform/request"
	"github.com/andikabahari/eoplatform/response"
	u "github.com/andikabahari/eoplatform/usecase"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type OrderHandler struct {
	usecase u.OrderUsecase
}

func NewOrderHandler(usecase u.OrderUsecase) *OrderHandler {
	return &OrderHandler{usecase}
}

func (h *OrderHandler) GetOrders(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*helper.JWTCustomClaims)

	orders := make([]model.Order, 0)
	payments := make([]model.Payment, len(orders))
	h.usecase.GetOrders(claims, &orders, &payments)

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

	order := model.Order{}

	if apiError := h.usecase.CreateOrder(claims, &order, &req); apiError != nil {
		code, message := apiError.APIError()
		return c.JSON(code, echo.Map{
			"message": "create order failure",
			"error":   message,
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "create order successful",
		"data":    response.NewOrderResponse(order),
	})
}

func (h *OrderHandler) AcceptOrCompleteOrder(c echo.Context) error {
	order := model.Order{}

	if apiError := h.usecase.AcceptOrCompleteOrder(c, &order); apiError != nil {
		code, message := apiError.APIError()
		return c.JSON(code, echo.Map{
			"message": "accept or complete order failure",
			"error":   message,
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "accept or complete order successful",
		"data":    response.NewOrderResponse(order),
	})
}

func (h *OrderHandler) CancelOrder(c echo.Context) error {
	order := model.Order{}

	if apiError := h.usecase.CancelOrder(c, &order); apiError != nil {
		code, message := apiError.APIError()
		return c.JSON(code, echo.Map{
			"message": "cancel order failure",
			"error":   message,
		})
	}

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

	if apiError := h.usecase.PaymentStatus(&req); apiError != nil {
		code, message := apiError.APIError()
		return c.JSON(code, echo.Map{
			"message": "payment failure",
			"error":   message,
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "payment successful",
	})
}
