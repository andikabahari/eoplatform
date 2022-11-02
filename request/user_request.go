package request

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (r CreateUserRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, validation.Length(1, 100)),
		validation.Field(&r.Username, validation.Required, validation.Length(4, 30), is.Alphanumeric),
		validation.Field(&r.Password, validation.Required, validation.Length(8, 100)),
		validation.Field(&r.Role, validation.Match(regexp.MustCompile("^(organizer|customer)$"))),
	)
}

type UpdateUserRequest struct {
	Name string `json:"name"`
}

func (r UpdateUserRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, validation.Length(1, 100)),
	)
}

type UpdateUserPasswordRequest struct {
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	OldPassword     string `json:"old_password"`
}

func (r UpdateUserPasswordRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Password, validation.Required, validation.Length(8, 100)),
		validation.Field(&r.OldPassword, validation.Required),
	)
}
