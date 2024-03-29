package request

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type CreateOrderRequest struct {
	DateOfEvent string `json:"date_of_event"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	Note        string `json:"note"`
	ServiceIDs  []uint `json:"service_ids"`
}

func (r CreateOrderRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.DateOfEvent, validation.Required, validation.Match(regexp.MustCompile(`^\d{1,4}-\d{1,2}-\d{1,2}$`))),
		validation.Field(&r.FirstName, validation.Required, validation.Length(1, 50)),
		validation.Field(&r.LastName, validation.Required, validation.Length(1, 50)),
		validation.Field(&r.Phone, validation.Required, validation.Length(1, 20)),
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(&r.Address, validation.Required, validation.Length(1, 300)),
		validation.Field(&r.Note, validation.Required, validation.Length(1, 300)),
		validation.Field(&r.ServiceIDs, validation.Required),
	)
}
