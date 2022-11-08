package repository

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/andikabahari/eoplatform/model"
	"github.com/andikabahari/eoplatform/request"
	"github.com/andikabahari/eoplatform/test/testhelper"
	"github.com/stretchr/testify/suite"
)

type paymentRepositorySuite struct {
	suite.Suite
	mock       sqlmock.Sqlmock
	repository *PaymentRepository
}

func (s *paymentRepositorySuite) SetupSuite() {
	var conn *sql.DB
	conn, s.mock = testhelper.Mock()
	gorm := testhelper.Init(conn)
	s.repository = NewPaymentRepository(gorm)
}

func TestPaymentRepositorySuite(t *testing.T) {
	suite.Run(t, new(paymentRepositorySuite))
}

func (s *paymentRepositorySuite) TestCreate() {
	query := regexp.QuoteMeta("INSERT INTO `payments`")
	s.mock.ExpectBegin()
	s.mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.repository.Create(&model.Payment{})
}

func (s *paymentRepositorySuite) TestUpdate() {
	query := regexp.QuoteMeta("INSERT INTO `payments`")
	s.mock.ExpectBegin()
	s.mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.repository.Update(&model.Payment{}, &request.MidtransTransactionNotificationRequest{})
}

func (s *paymentRepositorySuite) TestGetOnlyByOrderID() {
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	query := regexp.QuoteMeta("SELECT * FROM `payments`")
	s.mock.ExpectQuery(query).WillReturnRows(rows)
	s.repository.GetOnlyByOrderID(&[]model.Payment{}, 1)
}

func (s *paymentRepositorySuite) TestFindOnlyByOrderID() {
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	query := regexp.QuoteMeta("SELECT * FROM `payments`")
	s.mock.ExpectQuery(query).WillReturnRows(rows)
	s.repository.FindOnlyByOrderID(&model.Payment{}, 1)
}
