package usecase

import (
	"net/http"

	"github.com/andikabahari/eoplatform/helper"
	"github.com/andikabahari/eoplatform/model"
	r "github.com/andikabahari/eoplatform/repository"
	"github.com/andikabahari/eoplatform/request"
)

type BankAccountUsecase interface {
	GetBankAccount(claims *helper.JWTCustomClaims, bankAccount *model.BankAccount)
	CreateBankAccount(claims *helper.JWTCustomClaims, bankAccount *model.BankAccount, req *request.CreateBankAccountRequest)
	UpdateBankAccount(claims *helper.JWTCustomClaims, bankAccount *model.BankAccount, req *request.UpdateBankAccountRequest) helper.APIError
}

type bankAccountUsecase struct {
	bankAccountRepository r.BankAccountRepository
}

func NewBankAccountUsecase(bankAccountRepository r.BankAccountRepository) BankAccountUsecase {
	return &bankAccountUsecase{bankAccountRepository}
}

func (u *bankAccountUsecase) GetBankAccount(claims *helper.JWTCustomClaims, bankAccount *model.BankAccount) {
	u.bankAccountRepository.FindByUserID(bankAccount, claims.ID)
}

func (u *bankAccountUsecase) CreateBankAccount(claims *helper.JWTCustomClaims, bankAccount *model.BankAccount, req *request.CreateBankAccountRequest) {
	bankAccount.Bank = req.Bank
	bankAccount.VANumber = req.VANumber
	bankAccount.UserID = claims.ID

	u.bankAccountRepository.Create(bankAccount)
}

func (u *bankAccountUsecase) UpdateBankAccount(claims *helper.JWTCustomClaims, bankAccount *model.BankAccount, req *request.UpdateBankAccountRequest) helper.APIError {
	u.bankAccountRepository.FindByUserID(bankAccount, claims.ID)

	if bankAccount.ID == 0 {
		return helper.NewAPIError(http.StatusNotFound, "bank account not found")
	}

	u.bankAccountRepository.Update(bankAccount, req)

	return nil
}
