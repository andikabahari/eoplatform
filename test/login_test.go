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
	"golang.org/x/crypto/bcrypt"
)

type loginSuite struct {
	suite.Suite
	mock    sqlmock.Sqlmock
	server  *server.Server
	handler *handler.LoginHandler
}

func (s *loginSuite) SetupSuite() {
	os.Setenv("APP_ENV", "production")

	var conn *sql.DB
	conn, s.mock = testhelper.Mock()
	s.server = testhelper.NewServer(conn)
	s.handler = handler.NewLoginHandler(s.server)
}

func TestLoginSuite(t *testing.T) {
	suite.Run(t, new(loginSuite))
}

func (s *loginSuite) TestLogin() {
	testCases := []struct {
		Name         string
		Endpoint     string
		Method       string
		Body         *request.LoginRequest
		ExpectedCode int
	}{
		{
			"not found",
			"/v1/login",
			http.MethodPost,
			nil,
			http.StatusNotFound,
		},
		{
			"bad request",
			"/v1/login",
			http.MethodPost,
			&request.LoginRequest{
				Username: "organizer",
				Password: "password",
			},
			http.StatusBadRequest,
		},
		{
			"ok",
			"/v1/login",
			http.MethodPost,
			&request.LoginRequest{
				Username: "organizer",
				Password: "password",
			},
			http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		s.T().Run(testCase.Name, func(t *testing.T) {
			if testCase.ExpectedCode == http.StatusNotFound {
				query := regexp.QuoteMeta("SELECT * FROM `users`")
				s.mock.ExpectQuery(query).WillReturnRows(s.mock.NewRows(nil))
			}

			if testCase.ExpectedCode == http.StatusBadRequest {
				query := regexp.QuoteMeta("SELECT * FROM `users`")
				s.mock.ExpectQuery(query).WillReturnRows(s.mock.NewRows([]string{"id", "role", "password"}).AddRow(1, "organizer", "ok"))
			}

			if testCase.ExpectedCode == http.StatusOK {
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(testCase.Body.Password), s.server.Config.Auth.Cost)
				s.NoError(err)

				query := regexp.QuoteMeta("SELECT * FROM `users`")
				s.mock.ExpectQuery(query).WillReturnRows(s.mock.NewRows([]string{"id", "role", "password"}).AddRow(1, "organizer", hashedPassword))
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

			s.NoError(s.handler.Login(ctx))
			s.Equal(testCase.ExpectedCode, rec.Code)
		})
	}
}
