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

type BankAccountHandler struct {
	usecase u.BankAccountUsecase
}

func NewBankAccountHandler(usecase u.BankAccountUsecase) *BankAccountHandler {
	return &BankAccountHandler{usecase}
}

func (h *BankAccountHandler) GetBankAccounts(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*helper.JWTCustomClaims)

	if claims.Role != "organizer" {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "fetch bank account failure",
			"error":   "unauthorized",
		})
	}

	bankAccount := model.BankAccount{}
	h.usecase.GetBankAccount(claims, &bankAccount)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "fetch bank account successful",
		"data":    response.NewBankAccountResponse(bankAccount),
	})
}

func (h *BankAccountHandler) CreateBankAccount(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*helper.JWTCustomClaims)

	if claims.Role != "organizer" {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "create bank account failure",
			"error":   "unauthorized",
		})
	}

	req := request.CreateBankAccountRequest{}

	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "validation error",
			"error":   err,
		})
	}

	bankAccount := model.BankAccount{}
	h.usecase.CreateBankAccount(claims, &bankAccount, &req)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "create bank account successful",
		"data":    response.NewBankAccountResponse(bankAccount),
	})
}

func (h *BankAccountHandler) UpdateBankAccount(c echo.Context) error {
	req := request.UpdateBankAccountRequest{}

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

	bankAccount := model.BankAccount{}

	if apiError := h.usecase.UpdateBankAccount(claims, &bankAccount, &req); apiError != nil {
		code, message := apiError.APIError()
		return c.JSON(code, echo.Map{
			"message": "update bank account failure",
			"error":   message,
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "update bank account successful",
		"data":    response.NewBankAccountResponse(bankAccount),
	})
}
