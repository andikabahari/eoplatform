package repository

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/andikabahari/eoplatform/model"
	"github.com/andikabahari/eoplatform/testhelper"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type orderRepositorySuite struct {
	suite.Suite
	mock       sqlmock.Sqlmock
	repository OrderRepository
}

func (s *orderRepositorySuite) SetupSuite() {
	var conn *sql.DB
	conn, s.mock = testhelper.Mock()
	gorm := testhelper.Init(conn)
	s.repository = NewOrderRepository(gorm)
}

func TestOrderRepositorySuite(t *testing.T) {
	suite.Run(t, new(orderRepositorySuite))
}

func (s *orderRepositorySuite) TestGetOrdersForCustomer() {
	var query string
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	query = regexp.QuoteMeta("SELECT * FROM `orders`")
	s.mock.ExpectQuery(query).WillReturnRows(rows)
	query = regexp.QuoteMeta("SELECT * FROM `order_services`")
	s.mock.ExpectQuery(query).WillReturnRows(rows)
	s.repository.GetOrdersForCustomer(&[]model.Order{}, 1)
}

func (s *orderRepositorySuite) TestGetOrdersForOrganizer() {
	var query string
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	query = regexp.QuoteMeta("SELECT * FROM `orders` WHERE id IN (SELECT DISTINCT o.id FROM orders o JOIN users u ON u.id=o.user_id JOIN order_services os ON os.order_id=o.id JOIN services s ON s.id=os.service_id WHERE u.id!=? AND s.user_id=?) AND `orders`.`deleted_at` IS NULL")
	s.mock.ExpectQuery(query).WillReturnRows(rows)
	query = regexp.QuoteMeta("SELECT * FROM `order_services`")
	s.mock.ExpectQuery(query).WillReturnRows(rows)
	s.repository.GetOrdersForOrganizer(&[]model.Order{}, 1)
}

func (s *orderRepositorySuite) TestFind() {
	var query string
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	query = regexp.QuoteMeta("SELECT * FROM `orders`")
	s.mock.ExpectQuery(query).WithArgs("1").WillReturnRows(rows)
	query = regexp.QuoteMeta("SELECT * FROM `order_services`")
	s.mock.ExpectQuery(query).WillReturnRows(rows)
	s.repository.Find(&model.Order{}, "1")
}

func (s *orderRepositorySuite) TestFindOnly() {
	var query string
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	query = regexp.QuoteMeta("SELECT * FROM `orders`")
	s.mock.ExpectQuery(query).WithArgs("1").WillReturnRows(rows)
	s.repository.FindOnly(&model.Order{}, "1")
}

func (s *orderRepositorySuite) TestCreate() {
	query := regexp.QuoteMeta("INSERT INTO `orders`")
	s.mock.ExpectBegin()
	s.mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.repository.Create(&model.Order{})
}

func (s *orderRepositorySuite) TestDelete() {
	query := regexp.QuoteMeta("UPDATE `orders`")
	s.mock.ExpectBegin()
	s.mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()
	s.repository.Delete(&model.Order{Model: gorm.Model{ID: 1}})
}

func (s *orderRepositorySuite) TestSave() {
	query := regexp.QuoteMeta("INSERT INTO `orders`")
	s.mock.ExpectBegin()
	s.mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.repository.Save(&model.Order{})
}
