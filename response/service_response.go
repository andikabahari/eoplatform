package response

import "github.com/andikabahari/eoplatform/model"

type ServiceResponse struct {
	ID          uint          `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Cost        float64       `json:"cost"`
	User        *UserResponse `json:"user,omitempty"`
}

func NewServiceResponse(service model.Service) *ServiceResponse {
	res := ServiceResponse{}
	res.ID = service.ID
	res.Name = service.Name
	res.Description = service.Description
	res.Cost = service.Cost
	res.User = NewUserResponse(service.User)

	return &res
}

func NewServicesResponse(services []model.Service) *[]ServiceResponse {
	res := make([]ServiceResponse, 0)

	for _, service := range services {
		tmp := ServiceResponse{}
		tmp.ID = service.ID
		tmp.Name = service.Name
		tmp.Description = service.Description
		tmp.Cost = service.Cost
		tmp.User = NewUserResponse(service.User)
		res = append(res, tmp)
	}

	return &res
}
