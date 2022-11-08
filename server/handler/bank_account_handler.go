package handler

import (
	"net/http"

	"github.com/andikabahari/eoplatform/helper"
	"github.com/andikabahari/eoplatform/model"
	"github.com/andikabahari/eoplatform/repository"
	"github.com/andikabahari/eoplatform/request"
	"github.com/andikabahari/eoplatform/response"
	s "github.com/andikabahari/eoplatform/server"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type BankAccountHandler struct {
	server                *s.Server
	bankAccountRepository repository.BankAccountRepository
}

func NewBankAccountHandler(server *s.Server) *BankAccountHandler {
	return &BankAccountHandler{
		server,
		repository.NewBankAccountRepository(server.DB),
	}
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
	h.bankAccountRepository.FindByUserID(&bankAccount, claims.ID)

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
	bankAccount.Bank = req.Bank
	bankAccount.VANumber = req.VANumber
	bankAccount.UserID = claims.ID

	h.bankAccountRepository.Create(&bankAccount)

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
	h.bankAccountRepository.FindByUserID(&bankAccount, claims.ID)

	if bankAccount.ID == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "update bank account failure",
			"error":   "bank account not found",
		})
	}

	h.bankAccountRepository.Update(&bankAccount, &req)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "update bank account successful",
		"data":    response.NewBankAccountResponse(bankAccount),
	})
}
