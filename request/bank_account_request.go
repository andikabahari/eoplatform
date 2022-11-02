package request

import validation "github.com/go-ozzo/ozzo-validation"

type BasicBankAccount struct {
	Bank     string `json:"bank"`
	VANumber string `json:"va_number"`
}

func (b BasicBankAccount) Validate() error {
	return validation.ValidateStruct(&b,
		validation.Field(&b.Bank, validation.Required),
		validation.Field(&b.VANumber, validation.Required),
	)
}

type CreateBankAccountRequest struct {
	BasicBankAccount
}

type UpdateBankAccountRequest struct {
	BasicBankAccount
}
