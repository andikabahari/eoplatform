package handler

import (
	"net/http"

	"github.com/andikabahari/eoplatform/helper"
	"github.com/andikabahari/eoplatform/model"
	"github.com/andikabahari/eoplatform/request"
	"github.com/andikabahari/eoplatform/response"
	u "github.com/andikabahari/eoplatform/usecase"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type AccountHandler struct {
	usecase u.AccountUsecase
}

func NewAccountHandler(usecase u.AccountUsecase) *AccountHandler {
	return &AccountHandler{usecase}
}

func (h *AccountHandler) GetAccount(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*helper.JWTCustomClaims)

	user := model.User{}

	if apiError := h.usecase.GetAccount(claims, &user); apiError != nil {
		code, message := apiError.APIError()
		return c.JSON(code, echo.Map{
			"message": "fetch account failure",
			"error":   message,
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "fetch account successful",
		"data":    response.NewUserResponse(user),
	})
}

func (h *AccountHandler) UpdateAccount(c echo.Context) error {
	req := request.UpdateUserRequest{}

	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "validation error",
			"error":   err,
		})
	}

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*helper.JWTCustomClaims)

	user := model.User{}

	if apiError := h.usecase.UpdateAccount(claims, &user, &req); apiError != nil {
		code, message := apiError.APIError()
		return c.JSON(code, echo.Map{
			"message": "update account failure",
			"error":   message,
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "update account successful",
		"data":    response.NewUserResponse(user),
	})
}

func (h *AccountHandler) ResetPassword(c echo.Context) error {
	req := request.UpdateUserPasswordRequest{}

	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "validation error",
			"error":   err,
		})
	}

	if req.Password != req.ConfirmPassword {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "validation error",
			"error": echo.Map{
				"confirm_password": "must match \"password\"",
			},
		})
	}

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*helper.JWTCustomClaims)

	user := model.User{}

	if apiError := h.usecase.ResetPassword(claims, &user, &req); apiError != nil {
		code, message := apiError.APIError()
		return c.JSON(code, echo.Map{
			"message": "reset password failure",
			"error":   message,
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "reset password successful",
		"data": echo.Map{
			"kind":    "user",
			"id":      user.ID,
			"updated": true,
		},
	})
}
