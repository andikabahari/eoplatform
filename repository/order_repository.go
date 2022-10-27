package repository

import (
	"github.com/andikabahari/eoplatform/model"
	"gorm.io/gorm"
)

type OrderRepository interface {
	GetMyOrders(orders *[]model.Order, userID uint)
	Find(order *model.Order, id string)
	Create(order *model.Order)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *orderRepository {
	return &orderRepository{db}
}

func (r *orderRepository) GetMyOrders(orders *[]model.Order, userID uint) {
	r.db.Debug().Preload("User").Preload("Services").Where("user_id = ?", userID).Find(orders)
}

func (r *orderRepository) Create(order *model.Order) {
	r.db.Debug().Omit("Services.*").Save(order)
}

func (r *orderRepository) Find(order *model.Order, id string) {
	r.db.Debug().Preload("User").Preload("Services").Where("id = ?", id).Find(order)
}
