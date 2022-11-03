package test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/andikabahari/eoplatform/helper"
	"github.com/andikabahari/eoplatform/request"
	"github.com/andikabahari/eoplatform/server"
	"github.com/andikabahari/eoplatform/server/handler"
	"github.com/andikabahari/eoplatform/test/testhelper"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/suite"
)

type feedbackSuite struct {
	suite.Suite
	mock    sqlmock.Sqlmock
	server  *server.Server
	handler *handler.FeedbackHandler
}

func (s *feedbackSuite) SetupSuite() {
	os.Setenv("APP_ENV", "production")

	var conn *sql.DB
	conn, s.mock = testhelper.Mock()
	s.server = testhelper.NewServer(conn)
	s.handler = handler.NewFeedbackHandler(s.server)
}

func TestFeedbackSuite(t *testing.T) {
	suite.Run(t, new(feedbackSuite))
}

func (s *feedbackSuite) TestGetFeedbacks() {
	testCases := []struct {
		Name         string
		Endpoint     string
		Method       string
		Body         any
		ExpectedCode int
		Token        *jwt.Token
	}{
		{
			"ok",
			"/v1/feedbacks",
			http.MethodGet,
			nil,
			http.StatusOK,
			nil,
		},
	}

	for _, testCase := range testCases {
		s.T().Run(testCase.Name, func(t *testing.T) {
			bodyReader := new(bytes.Reader)
			if testCase.Body != nil {
				body, err := json.Marshal(testCase.Body)
				s.NoError(err)
				bodyReader = bytes.NewReader(body)
			}

			req := httptest.NewRequest(testCase.Method, testCase.Endpoint, bodyReader)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			ctx := s.server.Echo.NewContext(req, rec)
			ctx.Set("user", testCase.Token)

			s.NoError(s.handler.GetFeedbacks(ctx))
			s.Equal(testCase.ExpectedCode, rec.Code)
		})
	}
}

func (s *feedbackSuite) TestCreateFeedback() {
	type query struct {
		Raw  string
		Rows *sqlmock.Rows
	}

	testCases := []struct {
		Name         string
		Endpoint     string
		Method       string
		Body         *request.CreateFeedbackRequest
		ExpectedCode int
		Token        *jwt.Token
		Queries      []query
	}{
		{
			"unauthorized",
			"/v1/feedbacks",
			http.MethodPost,
			nil,
			http.StatusUnauthorized,
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{
				ID:   1,
				Role: "organizer",
			}),
			nil,
		},
		{
			"bad request",
			"/v1/feedbacks",
			http.MethodPost,
			nil,
			http.StatusBadRequest,
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{
				ID:   1,
				Role: "customer",
			}),
			nil,
		},
		{
			"forbidden",
			"/v1/feedbacks",
			http.MethodPost,
			&request.CreateFeedbackRequest{
				Description: "Good job!",
				Rating:      5,
				ToUserID:    1,
			},
			http.StatusForbidden,
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{
				ID:   1,
				Role: "customer",
			}),
			nil,
		},
		{
			"ok",
			"/v1/feedbacks",
			http.MethodPost,
			&request.CreateFeedbackRequest{
				Description: "Good job!",
				Rating:      5,
				ToUserID:    1,
			},
			http.StatusOK,
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{
				ID:   1,
				Role: "customer",
			}),
			[]query{
				{
					Raw:  regexp.QuoteMeta("SELECT COUNT(1) FROM feedbacks WHERE from_user_id=? AND to_user_id=?"),
					Rows: sqlmock.NewRows([]string{"count"}).AddRow(0),
				},
				{
					Raw:  regexp.QuoteMeta("SELECT COUNT(1) FROM(SELECT DISTINCT o.id FROM orders o JOIN order_services os ON os.order_id=o.id JOIN services s ON s.id=os.service_id WHERE o.user_id=? AND s.user_id=? AND is_completed>0) AS t"),
					Rows: sqlmock.NewRows([]string{"count"}).AddRow(1),
				},
			},
		},
		{
			"ok",
			"/v1/feedbacks",
			http.MethodPost,
			&request.CreateFeedbackRequest{
				Description: "This is bad!",
				Rating:      1,
				ToUserID:    1,
			},
			http.StatusOK,
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{
				ID:   1,
				Role: "customer",
			}),
			[]query{
				{
					Raw:  regexp.QuoteMeta("SELECT COUNT(1) FROM feedbacks WHERE from_user_id=? AND to_user_id=?"),
					Rows: sqlmock.NewRows([]string{"count"}).AddRow(0),
				},
				{
					Raw:  regexp.QuoteMeta("SELECT COUNT(1) FROM(SELECT DISTINCT o.id FROM orders o JOIN order_services os ON os.order_id=o.id JOIN services s ON s.id=os.service_id WHERE o.user_id=? AND s.user_id=? AND is_completed>0) AS t"),
					Rows: sqlmock.NewRows([]string{"count"}).AddRow(1),
				},
			},
		},
	}

	for _, testCase := range testCases {
		s.T().Run(testCase.Name, func(t *testing.T) {
			for _, query := range testCase.Queries {
				s.mock.ExpectQuery(query.Raw).WillReturnRows(query.Rows)
			}

			bodyReader := new(bytes.Reader)
			if testCase.Body != nil {
				body, err := json.Marshal(testCase.Body)
				s.NoError(err)
				bodyReader = bytes.NewReader(body)
			}

			req := httptest.NewRequest(testCase.Method, testCase.Endpoint, bodyReader)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			ctx := s.server.Echo.NewContext(req, rec)
			ctx.Set("user", testCase.Token)

			s.NoError(s.handler.CreateFeedback(ctx))
			s.Equal(testCase.ExpectedCode, rec.Code)
		})
	}
}
