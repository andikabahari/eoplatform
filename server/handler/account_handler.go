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
	"golang.org/x/crypto/bcrypt"
)

type AccountHandler struct {
	server *s.Server
}

func NewAccountHandler(server *s.Server) *AccountHandler {
	return &AccountHandler{server}
}

func (h *AccountHandler) GetAccount(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*helper.JWTCustomClaims)

	user := model.User{}
	userRepository := repository.NewUserRepository(h.server.DB)
	userRepository.Find(&user, fmt.Sprintf("%d", claims.ID))

	if user.ID == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "user not found",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": response.NewUserResponse(user),
	})
}

func (h *AccountHandler) GetMyOrders(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*helper.JWTCustomClaims)

	orders := make([]model.Order, 0)
	orderRepository := repository.NewOrderRepository(h.server.DB)
	orderRepository.GetMyOrders(&orders, claims.ID)

	return c.JSON(http.StatusOK, echo.Map{
		"data": response.NewMyOrdersResponse(orders),
	})
}

func (h *AccountHandler) GetCustomerOrders(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*helper.JWTCustomClaims)

	fmt.Println("hey", claims.ID)
	orders := make([]model.Order, 0)
	orderRepository := repository.NewOrderRepository(h.server.DB)
	orderRepository.GetCustomerOrders(&orders, claims.ID)

	return c.JSON(http.StatusOK, echo.Map{
		"data": response.NewCustomerOrdersResponse(orders),
	})
}

func (h *AccountHandler) UpdateAccount(c echo.Context) error {
	req := request.UpdateUserRequest{}

	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err,
		})
	}

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*helper.JWTCustomClaims)

	user := model.User{}
	userRepository := repository.NewUserRepository(h.server.DB)
	userRepository.Find(&user, fmt.Sprintf("%d", claims.ID))

	if user.ID == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "user not found",
		})
	}

	userRepository.Update(&user, &req)

	return c.JSON(http.StatusOK, echo.Map{
		"data": response.NewUserResponse(user),
	})
}

func (h *AccountHandler) ResetPassword(c echo.Context) error {
	req := request.UpdateUserPasswordRequest{}

	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err,
		})
	}

	if req.Password != req.ConfirmPassword {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": echo.Map{
				"confirm_password": "must match \"password\"",
			},
		})
	}

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*helper.JWTCustomClaims)

	user := model.User{}
	userRepository := repository.NewUserRepository(h.server.DB)
	userRepository.Find(&user, fmt.Sprintf("%d", claims.ID))

	if user.ID == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "user not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "unauthorized",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), h.server.Config.Auth.Cost)
	if err != nil {
		return err
	}

	userRepository.ResetPassword(&user, string(hashedPassword))

	return c.JSON(http.StatusOK, echo.Map{
		"data": echo.Map{
			"kind":    "user",
			"id":      user.ID,
			"updated": true,
		},
	})
}
