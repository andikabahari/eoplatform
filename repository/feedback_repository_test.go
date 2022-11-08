package repository

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/andikabahari/eoplatform/model"
	"github.com/andikabahari/eoplatform/test/testhelper"
	"github.com/stretchr/testify/suite"
)

type feedbackRepositorySuite struct {
	suite.Suite
	mock       sqlmock.Sqlmock
	repository *FeedbackRepository
}

func (s *feedbackRepositorySuite) SetupSuite() {
	var conn *sql.DB
	conn, s.mock = testhelper.Mock()
	gorm := testhelper.Init(conn)
	s.repository = NewFeedbackRepository(gorm)
}

func TestFeedbackRepositorySuite(t *testing.T) {
	suite.Run(t, new(feedbackRepositorySuite))
}

func (s *feedbackRepositorySuite) TestGet() {
	var query string
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	query = regexp.QuoteMeta("SELECT * FROM `feedbacks`")
	s.mock.ExpectQuery(query).WillReturnRows(rows)
	query = regexp.QuoteMeta("SELECT * FROM `feedbacks`")
	s.mock.ExpectQuery(query).WillReturnRows(rows)
	s.repository.Get(&[]model.Feedback{}, "")
	s.repository.Get(&[]model.Feedback{}, "1")
}

func (s *feedbackRepositorySuite) TestCreate() {
	query := regexp.QuoteMeta("INSERT INTO `feedbacks`")
	s.mock.ExpectBegin()
	s.mock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	s.repository.Create(&model.Feedback{})
}

func (s *feedbackRepositorySuite) TestGetFeedbacksCount() {
	rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	query := regexp.QuoteMeta("SELECT COUNT(1) FROM feedbacks WHERE from_user_id=? AND to_user_id=?")
	s.mock.ExpectQuery(query).WillReturnRows(rows)
	s.repository.GetFeedbacksCount(1, 2)
}

func (s *feedbackRepositorySuite) TestGetOrdersCount() {
	rows := sqlmock.NewRows([]string{"count"}).AddRow(1)
	query := regexp.QuoteMeta("SELECT COUNT(1) FROM(SELECT DISTINCT o.id FROM orders o JOIN order_services os ON os.order_id=o.id JOIN services s ON s.id=os.service_id WHERE o.user_id=? AND s.user_id=? AND is_completed>0) AS t")
	s.mock.ExpectQuery(query).WillReturnRows(rows)
	s.repository.GetOrdersCount(1, 2)
}
