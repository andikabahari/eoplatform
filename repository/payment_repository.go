package repository

import (
	"github.com/andikabahari/eoplatform/model"
	"github.com/andikabahari/eoplatform/request"
	"gorm.io/gorm"
)

type IPaymentRepository interface {
	Create(payment *model.Payment)
	Update(payment *model.Payment, req *request.MidtransTransactionNotificationRequest)
	GetOnlyByOrderID(payments *model.Payment, orderID any)
	FindOnlyByOrderID(payment *model.Payment, orderID any)
}

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db}
}

func (r *PaymentRepository) Create(payment *model.Payment) {
	r.db.Debug().Omit("Order").Save(payment)
}

func (r *PaymentRepository) Update(payment *model.Payment, req *request.MidtransTransactionNotificationRequest) {
	payment.Status = req.Status

	r.db.Debug().Omit("Order").Save(payment)
}

func (r *PaymentRepository) GetOnlyByOrderID(payments *[]model.Payment, orderID any) {
	r.db.Debug().Where("order_id = ?", orderID).Find(payments)
}

func (r *PaymentRepository) FindOnlyByOrderID(payment *model.Payment, orderID any) {
	r.db.Debug().Where("order_id = ?", orderID).Find(payment)
}
