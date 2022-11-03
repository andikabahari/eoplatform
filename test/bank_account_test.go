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

type bankAccountSuite struct {
	suite.Suite
	mock    sqlmock.Sqlmock
	server  *server.Server
	handler *handler.BankAccountHandler
}

func (s *bankAccountSuite) SetupSuite() {
	os.Setenv("APP_ENV", "production")

	var conn *sql.DB
	conn, s.mock = testhelper.Mock()
	s.server = testhelper.NewServer(conn)
	s.handler = handler.NewBankAccountHandler(s.server)
}

func TestBankAccountSuite(t *testing.T) {
	suite.Run(t, new(bankAccountSuite))
}

func (s *bankAccountSuite) TestGetBankAccounts() {
	testCases := []struct {
		Name         string
		Endpoint     string
		Method       string
		Body         any
		ExpectedCode int
	}{
		{
			"unauthorized",
			"/v1/bank-accounts",
			http.MethodGet,
			nil,
			http.StatusUnauthorized,
		},
		{
			"ok",
			"/v1/bank-accounts",
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

			if testCase.ExpectedCode == http.StatusUnauthorized {
				claims = &helper.JWTCustomClaims{}
			}

			if testCase.ExpectedCode == http.StatusOK {
				claims = &helper.JWTCustomClaims{
					Role: "organizer",
				}
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

			req := httptest.NewRequest(testCase.Method, testCase.Endpoint, bodyReader)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			ctx := s.server.Echo.NewContext(req, rec)
			ctx.Set("user", token)

			s.NoError(s.handler.GetBankAccounts(ctx))
			s.Equal(testCase.ExpectedCode, rec.Code)
		})
	}
}

func (s *bankAccountSuite) TestCreateBankAccount() {
	testCases := []struct {
		Name         string
		Endpoint     string
		Method       string
		Body         *request.CreateBankAccountRequest
		ExpectedCode int
	}{
		{
			"unauthorized",
			"/v1/bank-accounts",
			http.MethodPost,
			nil,
			http.StatusUnauthorized,
		},
		{
			"bad request",
			"/v1/bank-accounts",
			http.MethodPost,
			nil,
			http.StatusBadRequest,
		},
		{
			"ok",
			"/v1/bank-accounts",
			http.MethodPost,
			&request.CreateBankAccountRequest{
				BasicBankAccount: request.BasicBankAccount{
					Bank:     "bni",
					VANumber: "12345",
				},
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
				claims = &helper.JWTCustomClaims{}
			}

			if testCase.ExpectedCode == http.StatusBadRequest {
				claims = &helper.JWTCustomClaims{
					Role: "organizer",
				}
			}

			if testCase.ExpectedCode == http.StatusOK {
				claims = &helper.JWTCustomClaims{
					ID:   1,
					Role: "organizer",
				}
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

			req := httptest.NewRequest(testCase.Method, testCase.Endpoint, bodyReader)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			ctx := s.server.Echo.NewContext(req, rec)
			ctx.Set("user", token)

			s.NoError(s.handler.CreateBankAccount(ctx))
			s.Equal(testCase.ExpectedCode, rec.Code)
		})
	}
}

func (s *bankAccountSuite) TestUpdateBankAccount() {
	testCases := []struct {
		Name         string
		Endpoint     string
		Method       string
		Body         *request.UpdateBankAccountRequest
		ExpectedCode int
	}{
		{
			"bad request",
			"/v1/bank-accounts",
			http.MethodPut,
			nil,
			http.StatusBadRequest,
		},
		{
			"not found",
			"/v1/bank-accounts",
			http.MethodPut,
			&request.UpdateBankAccountRequest{
				BasicBankAccount: request.BasicBankAccount{
					Bank:     "bni",
					VANumber: "12345",
				},
			},
			http.StatusNotFound,
		},
		{
			"ok",
			"/v1/bank-accounts",
			http.MethodPut,
			&request.UpdateBankAccountRequest{
				BasicBankAccount: request.BasicBankAccount{
					Bank:     "bni",
					VANumber: "12345",
				},
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
				query := regexp.QuoteMeta("SELECT * FROM `bank_accounts`")
				s.mock.ExpectQuery(query).WillReturnRows(s.mock.NewRows([]string{"id"}).AddRow(1))
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

			req := httptest.NewRequest(testCase.Method, testCase.Endpoint, bodyReader)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			ctx := s.server.Echo.NewContext(req, rec)
			ctx.Set("user", token)

			s.NoError(s.handler.UpdateBankAccount(ctx))
			s.Equal(testCase.ExpectedCode, rec.Code)
		})
	}
}
