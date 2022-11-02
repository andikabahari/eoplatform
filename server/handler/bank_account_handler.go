package handler

import (
	"net/http"
	"strconv"

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
	server *s.Server
}

func NewBankAccountHandler(server *s.Server) *BankAccountHandler {
	return &BankAccountHandler{server}
}

func (h *BankAccountHandler) GetBankAccounts(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*helper.JWTCustomClaims)

	if claims.Role != "organizer" {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "unauthorized",
		})
	}

	bankAccounts := make([]model.BankAccount, 0)
	bankAccountRepository := repository.NewBankAccountRepository(h.server.DB)
	bankAccountRepository.Get(&bankAccounts, claims.ID)

	return c.JSON(http.StatusOK, echo.Map{
		"data": response.NewBankAccountsResponse(bankAccounts),
	})
}

func (h *BankAccountHandler) CreateBankAccount(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*helper.JWTCustomClaims)

	if claims.Role != "organizer" {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "unauthorized",
		})
	}

	req := request.CreateBankAccountRequest{}

	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err,
		})
	}

	bankAccount := model.BankAccount{}
	bankAccount.Bank = req.Bank
	bankAccount.VANumber = req.VANumber
	bankAccount.UserID = claims.ID

	bankAccountRepository := repository.NewBankAccountRepository(h.server.DB)
	bankAccountRepository.Create(&bankAccount)

	return c.JSON(http.StatusOK, echo.Map{
		"data": response.NewBankAccountResponse(bankAccount),
	})
}

func (h *BankAccountHandler) UpdateBankAccount(c echo.Context) error {
	req := request.UpdateBankAccountRequest{}

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

	bankAccount := model.BankAccount{}
	bankAccountRepository := repository.NewBankAccountRepository(h.server.DB)
	bankAccountRepository.FindByUserID(&bankAccount, claims.ID)

	if bankAccount.ID == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "bank account not found",
		})
	}

	bankAccountRepository.Update(&bankAccount, &req)

	return c.JSON(http.StatusOK, echo.Map{
		"data": response.NewBankAccountResponse(bankAccount),
	})
}

func (h *BankAccountHandler) DeleteBankAccount(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}

	bankAccount := model.BankAccount{}
	bankAccountRepository := repository.NewBankAccountRepository(h.server.DB)
	bankAccountRepository.Find(&bankAccount, uint(id))

	if bankAccount.ID == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "bank account not found",
		})
	}

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*helper.JWTCustomClaims)

	if bankAccount.UserID != claims.ID {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "unauthorized",
		})
	}

	bankAccountRepository.Delete(&bankAccount)

	return c.JSON(http.StatusOK, echo.Map{
		"data": echo.Map{
			"kind":    "bank account",
			"id":      id,
			"deleted": true,
		},
	})
}
