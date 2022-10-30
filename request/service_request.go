package request

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type BasicService struct {
	Name        string  `json:"name"`
	Cost        float64 `json:"cost"`
	Phone       string  `json:"phone"`
	Email       string  `json:"email"`
	IsPublished bool    `json:"is_published"`
	Description string  `json:"description"`
}

func (b BasicService) Validate() error {
	return validation.ValidateStruct(&b,
		validation.Field(&b.Name, validation.Required),
		validation.Field(&b.Cost, validation.Required),
		validation.Field(&b.Phone, validation.Required),
		validation.Field(&b.Email, validation.Required, is.Email),
		validation.Field(&b.IsPublished, validation.Required),
		validation.Field(&b.Description, validation.Required),
	)
}

type CreateServiceRequest struct {
	BasicService
}

type UpdateServiceRequest struct {
	BasicService
}
