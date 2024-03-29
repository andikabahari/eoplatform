package repository

import (
	"github.com/andikabahari/eoplatform/model"
	"github.com/andikabahari/eoplatform/request"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	Create(payment *model.Payment)
	Update(payment *model.Payment, req *request.MidtransTransactionNotificationRequest)
	GetOnlyByOrderID(payments *[]model.Payment, orderID any)
	FindOnlyByOrderID(payment *model.Payment, orderID any)
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db}
}

func (r *paymentRepository) Create(payment *model.Payment) {
	r.db.Debug().Omit("Order").Save(payment)
}

func (r *paymentRepository) Update(payment *model.Payment, req *request.MidtransTransactionNotificationRequest) {
	payment.Status = req.Status

	r.db.Debug().Omit("Order").Save(payment)
}

func (r *paymentRepository) GetOnlyByOrderID(payments *[]model.Payment, orderID any) {
	r.db.Debug().Where("order_id = ?", orderID).Find(payments)
}

func (r *paymentRepository) FindOnlyByOrderID(payment *model.Payment, orderID any) {
	r.db.Debug().Where("order_id = ?", orderID).Find(payment)
}
