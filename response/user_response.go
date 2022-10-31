package response

import "github.com/andikabahari/eoplatform/model"

type UserResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

func NewUserResponse(user model.User) *UserResponse {
	res := UserResponse{}
	res.ID = user.ID
	res.Name = user.Name
	res.Username = user.Username

	return &res
}
