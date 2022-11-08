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

type userRepositorySuite struct {
	suite.Suite
	mock       sqlmock.Sqlmock
	repository UserRepository
}

func (s *userRepositorySuite) SetupSuite() {
	var conn *sql.DB
	conn, s.mock = testhelper.Mock()
	gorm := testhelper.Init(conn)
	s.repository = NewUserRepository(gorm)
}

func TestUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(userRepositorySuite))
}

func (s *userRepositorySuite) TestFind() {
	query := regexp.QuoteMeta("SELECT * FROM `users`")
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	s.mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)
	s.repository.Find(&model.User{}, 1)
}

func (s *userRepositorySuite) TestFindByUsername() {
	query := regexp.QuoteMeta("SELECT * FROM `users`")
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	s.mock.ExpectQuery(query).WithArgs("user").WillReturnRows(rows)
	s.repository.FindByUsername(&model.User{}, "user")
}

func (s *userRepositorySuite) TestCreate() {
	query := regexp.QuoteMeta("INSERT INTO `users`")
	s.mock.ExpectBegin()
	s.mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.repository.Create(&model.User{})
}

func (s *userRepositorySuite) TestUpdate() {
	query := regexp.QuoteMeta("INSERT INTO `users`")
	s.mock.ExpectBegin()
	s.mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.repository.Update(&model.User{}, &request.UpdateUserRequest{})
}

func (s *userRepositorySuite) TestResetPassword() {
	query := regexp.QuoteMeta("INSERT INTO `users`")
	s.mock.ExpectBegin()
	s.mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.repository.ResetPassword(&model.User{}, "")
}
