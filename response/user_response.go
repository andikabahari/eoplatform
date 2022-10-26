package response

import "github.com/andikabahari/eoplatform/model"

type UserResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Address  string `json:"address"`
}

func NewUserResponse(user model.User) *UserResponse {
	res := UserResponse{}
	res.ID = user.ID
	res.Name = user.Name
	res.Username = user.Username
	res.Email = user.Email
	res.Address = user.Address

	return &res
}
