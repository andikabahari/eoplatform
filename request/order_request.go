package request

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type CreateOrderRequest struct {
	Phone      string `json:"phone"`
	Email      string `json:"email"`
	Address    string `json:"address"`
	ServiceIDs []uint `json:"service_ids"`
}

func (r CreateOrderRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Phone, validation.Required),
		validation.Field(&r.Email, validation.Required),
		validation.Field(&r.Address, validation.Required),
		validation.Field(&r.ServiceIDs, validation.Required),
	)
}
