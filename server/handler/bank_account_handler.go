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
	server *s.Server
}

func NewBankAccountHandler(server *s.Server) *BankAccountHandler {
	return &BankAccountHandler{server}
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
