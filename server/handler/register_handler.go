package handler

import (
	"net/http"

	"github.com/andikabahari/eoplatform/model"
	"github.com/andikabahari/eoplatform/request"
	"github.com/andikabahari/eoplatform/response"
	u "github.com/andikabahari/eoplatform/usecase"
	"github.com/labstack/echo/v4"
)

type RegisterHandler struct {
	usecase u.RegisterUsecase
}

func NewRegisterHandler(usecase u.RegisterUsecase) *RegisterHandler {
	return &RegisterHandler{usecase}
}

func (h *RegisterHandler) Register(c echo.Context) error {
	req := request.CreateUserRequest{}

	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "validation error",
			"error":   err,
		})
	}

	user := model.User{}

	if apiError := h.usecase.Register(&user, &req); apiError != nil {
		code, message := apiError.APIError()
		return c.JSON(code, echo.Map{
			"message": "registration failure",
			"error":   message,
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "registration successful",
		"data":    response.NewUserResponse(user),
	})
}
