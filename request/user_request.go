package request

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type BasicUser struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (bu BasicUser) Validate() error {
	return validation.ValidateStruct(&bu,
		validation.Field(&bu.Name, validation.Required),
		validation.Field(&bu.Username, validation.Required),
		validation.Field(&bu.Email, validation.Required, is.Email),
		validation.Field(&bu.Password, validation.Required, validation.Length(8, 0)),
		validation.Field(&bu.Role, validation.Required, validation.Match(regexp.MustCompile("^(organizer|customer)$"))),
	)
}

type CreateUserRequest struct {
	BasicUser
}
