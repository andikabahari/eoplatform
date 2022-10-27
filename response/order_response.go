package response

import (
	"time"

	"github.com/andikabahari/eoplatform/model"
)

type OrderResponse struct {
	ID          uint               `json:"id"`
	CreatedAt   time.Time          `json:"created_at"`
	IsAccepted  bool               `json:"is_accepted"`
	IsCompleted bool               `json:"is_completed"`
	User        *UserResponse      `json:"user,omitempty"`
	Services    *[]ServiceResponse `json:"services,omitempty"`
}

func NewOrderResponse(order model.Order) *OrderResponse {
	res := OrderResponse{}
	res.ID = order.ID
	res.CreatedAt = order.CreatedAt
	res.IsAccepted = order.IsAccepted
	res.IsCompleted = order.IsCompleted
	res.User = NewUserResponse(order.User)

	services := make([]ServiceResponse, 0)
	for _, service := range order.Services {
		tmp := ServiceResponse{}
		tmp.ID = service.ID
		tmp.Name = service.Name
		tmp.Description = service.Description
		tmp.Cost = service.Cost
		tmp.User = nil
		services = append(services, tmp)
	}

	res.Services = &services

	return &res
}
