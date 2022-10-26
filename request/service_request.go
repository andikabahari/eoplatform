package request

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type BasicService struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Cost        float64 `json:"cost"`
}

func (b BasicService) Validate() error {
	return validation.ValidateStruct(&b,
		validation.Field(&b.Name, validation.Required),
		validation.Field(&b.Description, validation.Required),
		validation.Field(&b.Cost, validation.Required),
	)
}

type CreateServiceRequest struct {
	BasicService
}
