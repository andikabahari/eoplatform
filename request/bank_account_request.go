package request

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
)

type BasicBankAccount struct {
	Bank     string `json:"bank"`
	VANumber string `json:"va_number"`
}

func (b BasicBankAccount) Validate() error {
	return validation.ValidateStruct(&b,
		validation.Field(&b.Bank, validation.Required, validation.Match(regexp.MustCompile("^(bni|bri|bca)$"))),
		validation.Field(&b.VANumber, validation.Required, validation.Length(1, 50)),
	)
}

type CreateBankAccountRequest struct {
	BasicBankAccount
}

type UpdateBankAccountRequest struct {
	BasicBankAccount
}
