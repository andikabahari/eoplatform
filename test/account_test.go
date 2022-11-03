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
	"golang.org/x/crypto/bcrypt"
)

type accountSuite struct {
	suite.Suite
	mock    sqlmock.Sqlmock
	server  *server.Server
	handler *handler.AccountHandler
}

func (s *accountSuite) SetupSuite() {
	os.Setenv("APP_ENV", "production")

	var conn *sql.DB
	conn, s.mock = testhelper.Mock()
	s.server = testhelper.NewServer(conn)
	s.handler = handler.NewAccountHandler(s.server)
}

func TestAccountSuite(t *testing.T) {
	suite.Run(t, new(accountSuite))
}

func (s *accountSuite) TestGetAccount() {
	testCases := []struct {
		Name         string
		Endpoint     string
		Method       string
		Body         any
		ExpectedCode int
	}{
		{
			"not found",
			"/v1/account",
			http.MethodGet,
			nil,
			http.StatusNotFound,
		},
		{
			"ok",
			"/v1/account",
			http.MethodGet,
			nil,
			http.StatusOK,
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

			claims := new(helper.JWTCustomClaims)

			if testCase.ExpectedCode == http.StatusNotFound {
				claims = &helper.JWTCustomClaims{}
			}

			if testCase.ExpectedCode == http.StatusOK {
				claims = &helper.JWTCustomClaims{ID: 1}
				query := regexp.QuoteMeta("SELECT * FROM `users`")
				s.mock.ExpectQuery(query).WillReturnRows(s.mock.NewRows([]string{"id"}).AddRow(1))

			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

			req := httptest.NewRequest(testCase.Method, testCase.Endpoint, bodyReader)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			ctx := s.server.Echo.NewContext(req, rec)
			ctx.Set("user", token)

			s.NoError(s.handler.GetAccount(ctx))
			s.Equal(testCase.ExpectedCode, rec.Code)
		})
	}
}

func (s *accountSuite) TestUpdateAccount() {
	testCases := []struct {
		Name         string
		Endpoint     string
		Method       string
		Body         *request.UpdateUserRequest
		ExpectedCode int
	}{
		{
			"bad request",
			"/v1/account",
			http.MethodPost,
			nil,
			http.StatusBadRequest,
		},
		{
			"not found",
			"/v1/account",
			http.MethodPost,
			&request.UpdateUserRequest{
				Name: "User",
			},
			http.StatusNotFound,
		},
		{
			"ok",
			"/v1/account",
			http.MethodPost,
			&request.UpdateUserRequest{
				Name: "User",
			},
			http.StatusOK,
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

			claims := new(helper.JWTCustomClaims)

			if testCase.ExpectedCode == http.StatusOK {
				claims = &helper.JWTCustomClaims{ID: 1}
				query := regexp.QuoteMeta("SELECT * FROM `users`")
				s.mock.ExpectQuery(query).WillReturnRows(s.mock.NewRows([]string{"id"}).AddRow(1))
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

			req := httptest.NewRequest(testCase.Method, testCase.Endpoint, bodyReader)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			ctx := s.server.Echo.NewContext(req, rec)
			ctx.Set("user", token)

			s.NoError(s.handler.UpdateAccount(ctx))
			s.Equal(testCase.ExpectedCode, rec.Code)
		})
	}
}

func (s *accountSuite) TestResetPassword() {
	testCases := []struct {
		Name         string
		Endpoint     string
		Method       string
		Body         *request.UpdateUserPasswordRequest
		ExpectedCode int
	}{
		{
			"bad request",
			"/v1/account/password",
			http.MethodPost,
			nil,
			http.StatusBadRequest,
		},
		{
			"bad request",
			"/v1/account/password",
			http.MethodPost,
			&request.UpdateUserPasswordRequest{
				Password:        "password",
				ConfirmPassword: "password1",
				OldPassword:     "password",
			},
			http.StatusBadRequest,
		},
		{
			"not found",
			"/v1/account/password",
			http.MethodPost,
			&request.UpdateUserPasswordRequest{
				Password:        "password",
				ConfirmPassword: "password",
				OldPassword:     "password",
			},
			http.StatusNotFound,
		},
		{
			"unauthorized",
			"/v1/account/password",
			http.MethodPost,
			&request.UpdateUserPasswordRequest{
				Password:        "password",
				ConfirmPassword: "password",
				OldPassword:     "wrong",
			},
			http.StatusUnauthorized,
		},
		{
			"ok",
			"/v1/account/password",
			http.MethodPost,
			&request.UpdateUserPasswordRequest{
				Password:        "password",
				ConfirmPassword: "password",
				OldPassword:     "password",
			},
			http.StatusOK,
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

			claims := new(helper.JWTCustomClaims)

			if testCase.ExpectedCode == http.StatusUnauthorized {
				claims = &helper.JWTCustomClaims{ID: 1}

				hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), s.server.Config.Auth.Cost)
				s.NoError(err)

				query := regexp.QuoteMeta("SELECT * FROM `users`")
				s.mock.ExpectQuery(query).WillReturnRows(s.mock.NewRows([]string{"id", "password"}).AddRow(1, hashedPassword))
			}

			if testCase.ExpectedCode == http.StatusOK {
				claims = &helper.JWTCustomClaims{ID: 1}

				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(testCase.Body.OldPassword), s.server.Config.Auth.Cost)
				s.NoError(err)

				query := regexp.QuoteMeta("SELECT * FROM `users`")
				s.mock.ExpectQuery(query).WillReturnRows(s.mock.NewRows([]string{"id", "password"}).AddRow(1, hashedPassword))
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

			req := httptest.NewRequest(testCase.Method, testCase.Endpoint, bodyReader)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			ctx := s.server.Echo.NewContext(req, rec)
			ctx.Set("user", token)

			s.NoError(s.handler.ResetPassword(ctx))
			s.Equal(testCase.ExpectedCode, rec.Code)
		})
	}
}
