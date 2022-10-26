package response

import "github.com/andikabahari/eoplatform/model"

type ServiceResponse struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Cost        float64 `json:"cost"`
}

func NewServiceResponse(service model.Service) ServiceResponse {
	res := ServiceResponse{}
	res.ID = service.ID
	res.Name = service.Name
	res.Description = service.Description
	res.Cost = service.Cost

	return res
}
