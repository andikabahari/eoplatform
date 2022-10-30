package response

import "github.com/andikabahari/eoplatform/model"

type ServiceResponse struct {
	ID          uint          `json:"id"`
	Name        string        `json:"name"`
	Cost        float64       `json:"cost"`
	Phone       string        `json:"phone"`
	Email       string        `json:"email"`
	IsPublished bool          `json:"is_published"`
	Description string        `json:"description"`
	User        *UserResponse `json:"user,omitempty"`
}

func NewServiceResponse(service model.Service) *ServiceResponse {
	res := ServiceResponse{}
	res.ID = service.ID
	res.Name = service.Name
	res.Cost = service.Cost
	res.Phone = service.Phone
	res.Email = service.Email
	res.IsPublished = service.IsPublished
	res.Description = service.Description
	res.User = NewUserResponse(service.User)

	return &res
}

func NewServicesResponse(services []model.Service) *[]ServiceResponse {
	res := make([]ServiceResponse, 0)

	for _, service := range services {
		tmp := ServiceResponse{}
		tmp.ID = service.ID
		tmp.Name = service.Name
		tmp.Cost = service.Cost
		tmp.Phone = service.Phone
		tmp.Email = service.Email
		tmp.IsPublished = service.IsPublished
		tmp.Description = service.Description
		tmp.User = NewUserResponse(service.User)
		res = append(res, tmp)
	}

	return &res
}
