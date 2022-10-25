package response

import "github.com/andikabahari/eoplatform/model"

type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

func NewUserResponse(user model.User) UserResponse {
	res := UserResponse{}
	res.ID = user.ID
	res.Name = user.Name
	res.Email = user.Email
	res.Role = user.Role

	return res
}
