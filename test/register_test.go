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
	"github.com/andikabahari/eoplatform/request"
	"github.com/andikabahari/eoplatform/server"
	"github.com/andikabahari/eoplatform/server/handler"
	"github.com/andikabahari/eoplatform/test/testhelper"
	"github.com/stretchr/testify/suite"
)

type registerSuite struct {
	suite.Suite
	mock    sqlmock.Sqlmock
	server  *server.Server
	handler *handler.RegisterHandler
}

func (s *registerSuite) SetupSuite() {
	os.Setenv("APP_ENV", "production")

	var conn *sql.DB
	conn, s.mock = testhelper.Mock()
	s.server = testhelper.NewServer(conn)
	s.handler = handler.NewRegisterHandler(s.server)
}

func TestRegisterSuite(t *testing.T) {
	suite.Run(t, new(registerSuite))
}

func (s *registerSuite) TestRegister() {
	testCases := []struct {
		Name         string
		Endpoint     string
		Method       string
		Body         *request.CreateUserRequest
		ExpectedCode int
	}{
		{
			"bad request",
			"/v1/register",
			http.MethodPost,
			nil,
			http.StatusBadRequest,
		},
		{
			"ok",
			"/v1/register",
			http.MethodPost,
			&request.CreateUserRequest{
				Name:     "Organizer",
				Username: "organizer",
				Password: "password",
				Role:     "organizer",
			},
			http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		s.T().Run(testCase.Name, func(t *testing.T) {
			query := regexp.QuoteMeta("SELECT * FROM `users` WHERE username = ? AND `users`.`deleted_at` IS NULL ORDER BY `users`.`id` LIMIT 1")
			s.mock.ExpectQuery(query).WillReturnRows(s.mock.NewRows(nil))

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

			s.NoError(s.handler.Register(ctx))
			s.Equal(testCase.ExpectedCode, rec.Code)
		})
	}
}
