package request

import validation "github.com/go-ozzo/ozzo-validation"

type CreateBankAccountRequest struct {
	Bank     string `json:"bank"`
	VANumber string `json:"va_number"`
}

func (r CreateBankAccountRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Bank, validation.Required),
		validation.Field(&r.VANumber, validation.Required),
	)
}
