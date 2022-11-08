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

type serviceRepositorySuite struct {
	suite.Suite
	mock       sqlmock.Sqlmock
	repository ServiceRepository
}

func (s *serviceRepositorySuite) SetupSuite() {
	var conn *sql.DB
	conn, s.mock = testhelper.Mock()
	gorm := testhelper.Init(conn)
	s.repository = NewServiceRepository(gorm)
}

func TestServiceRepositorySuite(t *testing.T) {
	suite.Run(t, new(serviceRepositorySuite))
}

func (s *serviceRepositorySuite) TestGet() {
	var query string
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	query = regexp.QuoteMeta("SELECT * FROM `services`")
	s.mock.ExpectQuery(query).WillReturnRows(rows)
	query = regexp.QuoteMeta("SELECT * FROM `services`")
	s.mock.ExpectQuery(query).WillReturnRows(rows)
	s.repository.Get(&[]model.Service{}, "")
	s.repository.Get(&[]model.Service{}, "any")
}

func (s *serviceRepositorySuite) TestFind() {
	query := regexp.QuoteMeta("SELECT * FROM `services`")
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	s.mock.ExpectQuery(query).WithArgs("1").WillReturnRows(rows)
	s.repository.Find(&model.Service{}, "1")
}

func (s *serviceRepositorySuite) TestCreate() {
	query := regexp.QuoteMeta("INSERT INTO `services`")
	s.mock.ExpectBegin()
	s.mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.repository.Create(&model.Service{})
}

func (s *serviceRepositorySuite) TestUpdate() {
	query := regexp.QuoteMeta("INSERT INTO `services`")
	s.mock.ExpectBegin()
	s.mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.repository.Update(&model.Service{}, &request.UpdateServiceRequest{})
}

func (s *serviceRepositorySuite) TestDelete() {
	query := regexp.QuoteMeta("UPDATE `services`")
	s.mock.ExpectBegin()
	s.mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(0, 1))
	s.mock.ExpectCommit()
	s.repository.Delete(&model.Service{Model: gorm.Model{ID: 1}})
}
