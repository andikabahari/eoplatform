package request

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type BasicUser struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Address  string `json:"address"`
}

func (b BasicUser) Validate() error {
	return validation.ValidateStruct(&b,
		validation.Field(&b.Name, validation.Required),
		validation.Field(&b.Username, validation.Required),
		validation.Field(&b.Password, validation.Required, validation.Length(8, 0)),
		validation.Field(&b.Email, validation.Required, is.Email),
		validation.Field(&b.Address, validation.Required),
	)
}

type CreateUserRequest struct {
	BasicUser
}
