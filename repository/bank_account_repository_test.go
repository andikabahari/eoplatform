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
	"gorm.io/gorm"
)

type bankAccountRepositorySuite struct {
	suite.Suite
	mock       sqlmock.Sqlmock
	repository *BankAccountRepository
}

func (s *bankAccountRepositorySuite) SetupSuite() {
	var conn *sql.DB
	conn, s.mock = testhelper.Mock()
	gorm := testhelper.Init(conn)
	s.repository = NewBankAccountRepository(gorm)
}

func TestBankAccountRepositorySuite(t *testing.T) {
	suite.Run(t, new(bankAccountRepositorySuite))
}

func (s *bankAccountRepositorySuite) TestGet() {
	query := regexp.QuoteMeta("SELECT * FROM `bank_accounts`")
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	s.mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)
	s.repository.Get(&[]model.BankAccount{}, 1)
}

func (s *bankAccountRepositorySuite) TestFind() {
	query := regexp.QuoteMeta("SELECT * FROM `bank_accounts`")
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	s.mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)
	s.repository.Find(&model.BankAccount{}, 1)
}

func (s *bankAccountRepositorySuite) TestFindByUserID() {
	query := regexp.QuoteMeta("SELECT * FROM `bank_accounts`")
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	s.mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)
	s.repository.FindByUserID(&model.BankAccount{}, 1)
}

func (s *bankAccountRepositorySuite) TestCreate() {
	query := regexp.QuoteMeta("INSERT INTO `bank_accounts`")
	s.mock.ExpectBegin()
	s.mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.repository.Create(&model.BankAccount{})
}

func (s *bankAccountRepositorySuite) TestUpdate() {
	query := regexp.QuoteMeta("INSERT INTO `bank_accounts`")
	s.mock.ExpectBegin()
	s.mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.repository.Update(&model.BankAccount{}, &request.UpdateBankAccountRequest{})
}

func (s *bankAccountRepositorySuite) TestDelete() {
	query := regexp.QuoteMeta("UPDATE `bank_accounts`")
	s.mock.ExpectBegin()
	s.mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()
	s.repository.Delete(&model.BankAccount{Model: gorm.Model{ID: 1}})
}
