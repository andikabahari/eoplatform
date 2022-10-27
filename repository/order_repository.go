package repository

import (
	"github.com/andikabahari/eoplatform/model"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *model.Order)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *orderRepository {
	return &orderRepository{db}
}

func (r *orderRepository) Create(order *model.Order) {
	r.db.Debug().Omit("Services.*").Save(order)
}
