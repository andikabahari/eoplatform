package response

import "github.com/andikabahari/eoplatform/model"

type BankAccountResponse struct {
	ID       uint          `json:"id"`
	Bank     string        `json:"bank"`
	VANumber string        `json:"va_number"`
	User     *UserResponse `json:"user,omitempty"`
}

func NewBankAccountResponse(bankAccount model.BankAccount) *BankAccountResponse {
	res := BankAccountResponse{}
	res.ID = bankAccount.ID
	res.Bank = bankAccount.Bank
	res.VANumber = bankAccount.VANumber
	if bankAccount.User.ID > 0 {
		res.User = NewUserResponse(bankAccount.User)
	}

	return &res
}

func NewBankAccountsResponse(bankAccounts []model.BankAccount) *[]BankAccountResponse {
	res := make([]BankAccountResponse, 0)

	for _, bankAccount := range bankAccounts {
		tmp := BankAccountResponse{}
		tmp.ID = bankAccount.ID
		tmp.Bank = bankAccount.Bank
		tmp.VANumber = bankAccount.VANumber
		if bankAccount.User.ID > 0 {
			tmp.User = NewUserResponse(bankAccount.User)
		}
		res = append(res, tmp)
	}

	return &res
}
