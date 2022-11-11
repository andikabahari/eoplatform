package handler

import (
	"net/http"

	"github.com/andikabahari/eoplatform/request"
	u "github.com/andikabahari/eoplatform/usecase"
	"github.com/labstack/echo/v4"
)

type LoginHandler struct {
	usecase u.LoginUsecase
}

func NewLoginHandler(usecase u.LoginUsecase) *LoginHandler {
	return &LoginHandler{usecase}
}

func (h *LoginHandler) Login(c echo.Context) error {
	req := request.LoginRequest{}

	if err := c.Bind(&req); err != nil {
		return err
	}

	token := ""

	if apiError := h.usecase.Login(&token, &req); apiError != nil {
		code, message := apiError.APIError()
		return c.JSON(code, echo.Map{
			"message": "login failure",
			"error":   message,
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "login successful",
		"data":    token,
	})
}
