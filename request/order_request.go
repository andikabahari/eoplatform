package request

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type BasicOrder struct {
	ServiceIDs []uint `json:"service_ids"`
}

func (b BasicOrder) Validate() error {
	return validation.ValidateStruct(&b,
		validation.Field(&b.ServiceIDs, validation.Required),
	)
}

type CreateOrderRequest struct {
	BasicOrder
}
