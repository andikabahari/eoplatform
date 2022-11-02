package repository

import (
	"github.com/andikabahari/eoplatform/model"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	Create(payment *model.Payment)
	FindOnlyByOrderID(payment *model.Payment, orderID any)
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *paymentRepository {
	return &paymentRepository{db}
}

func (r *paymentRepository) Create(payment *model.Payment) {
	r.db.Debug().Omit("Order").Save(&payment)
}

func (r *paymentRepository) FindOnlyByOrderID(payment *model.Payment, orderID any) {
	r.db.Debug().Where("order_id = ?", orderID).Find(payment)
}
