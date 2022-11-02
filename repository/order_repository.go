package repository

import (
	"database/sql"

	"github.com/andikabahari/eoplatform/model"
	"gorm.io/gorm"
)

type OrderRepository interface {
	GetMyOrders(orders *[]model.Order, userID uint)
	GetCustomerOrders(orders *[]model.Order, userID uint)
	GetOrdersForCustomer(orders *[]model.Order, userID uint)
	GetOrdersForOrganizer(orders *[]model.Order, userID uint)
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

func (r *orderRepository) GetCustomerOrders(orders *[]model.Order, userID uint) {
	query := "SELECT DISTINCT o.id FROM orders o " +
		"JOIN users u ON u.id=o.user_id " +
		"JOIN order_services os ON os.order_id=o.id " +
		"JOIN services s ON s.id=os.service_id " +
		"WHERE u.id!=@UserID AND s.user_id=@ServiceUserID"

	r.db.Debug().Preload("User").Preload("Services").Where("id IN (?)", r.db.Raw(query,
		sql.Named("UserID", userID),
		sql.Named("ServiceUserID", userID),
	)).Find(orders)
}

func (r *orderRepository) GetOrdersForCustomer(orders *[]model.Order, userID uint) {
	r.db.Debug().Preload("User").Preload("Services").Where("user_id = ?", userID).Find(orders)
}

func (r *orderRepository) GetOrdersForOrganizer(orders *[]model.Order, userID uint) {
	query := "SELECT DISTINCT o.id FROM orders o " +
		"JOIN users u ON u.id=o.user_id " +
		"JOIN order_services os ON os.order_id=o.id " +
		"JOIN services s ON s.id=os.service_id " +
		"WHERE u.id!=@UserID AND s.user_id=@ServiceUserID"

	r.db.Debug().Preload("User").Preload("Services").Where("id IN (?)", r.db.Raw(query,
		sql.Named("UserID", userID),
		sql.Named("ServiceUserID", userID),
	)).Find(orders)
}

func (r *orderRepository) Create(order *model.Order) {
	r.db.Debug().Omit("Services.*").Save(order)
}

func (r *orderRepository) Find(order *model.Order, id string) {
	r.db.Debug().Preload("User").Preload("Services").Where("id = ?", id).Find(order)
}

func (r *orderRepository) FindOnly(order *model.Order, id any) {
	r.db.Debug().Where("id = ?", id).Find(order)
}

func (r *orderRepository) Delete(order *model.Order) {
	r.db.Debug().Delete(order)
}
