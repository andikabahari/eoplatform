package response

import (
	"time"

	"github.com/andikabahari/eoplatform/model"
)

type OrderResponse struct {
	ID          uint               `json:"id"`
	CreatedAt   time.Time          `json:"created_at"`
	TotalCost   float64            `json:"total_cost"`
	IsAccepted  bool               `json:"is_accepted"`
	IsCompleted bool               `json:"is_completed"`
	Phone       string             `json:"phone"`
	Email       string             `json:"email"`
	Address     string             `json:"address"`
	User        *UserResponse      `json:"user,omitempty"`
	Services    *[]ServiceResponse `json:"services,omitempty"`
}

func NewOrderResponse(order model.Order) *OrderResponse {
	res := OrderResponse{}
	res.ID = order.ID
	res.CreatedAt = order.CreatedAt
	res.IsAccepted = order.IsAccepted
	res.IsCompleted = order.IsCompleted
	res.Phone = order.Phone
	res.Email = order.Email
	res.Address = order.Address
	res.User = NewUserResponse(order.User)

	services := make([]ServiceResponse, 0)
	for _, service := range order.Services {
		res.TotalCost += service.Cost

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

func NewMyOrdersResponse(orders []model.Order) *[]OrderResponse {
	res := make([]OrderResponse, 0)
	for i, order := range orders {
		tmp := OrderResponse{}
		tmp.ID = order.ID
		tmp.CreatedAt = order.CreatedAt
		tmp.IsAccepted = order.IsAccepted
		tmp.IsCompleted = order.IsCompleted
		tmp.Phone = order.Phone
		tmp.Email = order.Email
		tmp.Address = order.Address
		tmp.User = nil
		res = append(res, tmp)

		var totalCost float64

		services := make([]ServiceResponse, 0)
		for _, service := range order.Services {
			totalCost += service.Cost

			tmp := ServiceResponse{}
			tmp.ID = service.ID
			tmp.Name = service.Name
			tmp.Description = service.Description
			tmp.Cost = service.Cost
			tmp.User = nil
			services = append(services, tmp)
		}

		res[i].TotalCost = totalCost
		res[i].Services = &services
	}

	return &res
}

func NewCustomerOrdersResponse(orders []model.Order) *[]OrderResponse {
	res := make([]OrderResponse, 0)
	for i, order := range orders {
		tmp := OrderResponse{}
		tmp.ID = order.ID
		tmp.CreatedAt = order.CreatedAt
		tmp.IsAccepted = order.IsAccepted
		tmp.IsCompleted = order.IsCompleted
		tmp.Phone = order.Phone
		tmp.Email = order.Email
		tmp.Address = order.Address
		tmp.User = NewUserResponse(order.User)
		res = append(res, tmp)

		var totalCost float64

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

		res[i].TotalCost = totalCost
		res[i].Services = &services
	}

	return &res
}
