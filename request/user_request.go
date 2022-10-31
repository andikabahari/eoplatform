package request

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r CreateUserRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Username, validation.Required),
		validation.Field(&r.Password, validation.Required, validation.Length(8, 0)),
	)
}

type UpdateUserRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
}

func (r UpdateUserRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Username, validation.Required),
	)
}

type UpdateUserPasswordRequest struct {
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	OldPassword     string `json:"old_password"`
}

func (r UpdateUserPasswordRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Password, validation.Required, validation.Length(8, 0)),
		validation.Field(&r.OldPassword, validation.Required),
	)
}
