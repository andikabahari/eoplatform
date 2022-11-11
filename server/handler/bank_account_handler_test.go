package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/andikabahari/eoplatform/helper"
	"github.com/andikabahari/eoplatform/request"
	"github.com/andikabahari/eoplatform/server"
	"github.com/andikabahari/eoplatform/testhelper"
	mu "github.com/andikabahari/eoplatform/usecase/mock_usecase"
	"github.com/golang-jwt/jwt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type bankAccountHandlerSuite struct {
	suite.Suite

	ctrl    *gomock.Controller
	usecase *mu.MockBankAccountUsecase

	server  *server.Server
	handler *BankAccountHandler
}

func (s *bankAccountHandlerSuite) SetupSuite() {
	os.Setenv("APP_ENV", "production")

	s.ctrl = gomock.NewController(s.T())
	s.usecase = mu.NewMockBankAccountUsecase(s.ctrl)

	conn, _ := testhelper.Mock()
	s.server = testhelper.NewServer(conn)
	s.handler = NewBankAccountHandler(s.usecase)
}

func (s *bankAccountHandlerSuite) TearDownSuite() {
	s.ctrl.Finish()
}

func TestBankAccountHandlerSuite(t *testing.T) {
	suite.Run(t, new(bankAccountHandlerSuite))
}

func (s *bankAccountHandlerSuite) TestGetBankAccounts() {
	testCases := []struct {
		Name         string
		Endpoint     string
		Method       string
		Body         any
		Token        *jwt.Token
		ExpectedFunc func()
		ExpectedCode int
	}{
		{
			"unauthorized",
			"/v1/bank-accounts",
			http.MethodGet,
			nil,
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{}),
			func() {},
			http.StatusUnauthorized,
		},
		{
			"ok",
			"/v1/bank-accounts",
			http.MethodGet,
			nil,
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{ID: 1, Role: "organizer"}),
			func() {
				s.usecase.EXPECT().GetBankAccount(gomock.Any(), gomock.Any())
			},
			http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		s.T().Run(testCase.Name, func(t *testing.T) {
			testCase.ExpectedFunc()

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

			s.NoError(s.handler.GetBankAccounts(ctx))
			s.Equal(testCase.ExpectedCode, rec.Code)
		})
	}
}

func (s *bankAccountHandlerSuite) TestCreateBankAccount() {
	testCases := []struct {
		Name         string
		Endpoint     string
		Method       string
		Body         *request.CreateBankAccountRequest
		Token        *jwt.Token
		ExpectedFunc func()
		ExpectedCode int
	}{
		{
			"unauthorized",
			"/v1/bank-accounts",
			http.MethodPost,
			nil,
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{}),
			func() {},
			http.StatusUnauthorized,
		},
		{
			"bad request",
			"/v1/bank-accounts",
			http.MethodPost,
			nil,
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{ID: 1, Role: "organizer"}),
			func() {},
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
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{ID: 1, Role: "organizer"}),
			func() {
				s.usecase.EXPECT().CreateBankAccount(gomock.Any(), gomock.Any(), gomock.Any())
			},
			http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		s.T().Run(testCase.Name, func(t *testing.T) {
			testCase.ExpectedFunc()

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

			s.NoError(s.handler.CreateBankAccount(ctx))
			s.Equal(testCase.ExpectedCode, rec.Code)
		})
	}
}

func (s *bankAccountHandlerSuite) TestUpdateBankAccount() {
	testCases := []struct {
		Name         string
		Endpoint     string
		Method       string
		Body         *request.UpdateBankAccountRequest
		Token        *jwt.Token
		ExpectedFunc func()
		ExpectedCode int
	}{
		{
			"bad request",
			"/v1/bank-accounts",
			http.MethodPut,
			nil,
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{}),
			func() {},
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
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{}),
			func() {
				apiError := helper.NewAPIError(http.StatusNotFound, "")
				s.usecase.EXPECT().UpdateBankAccount(gomock.Any(), gomock.Any(), gomock.Any()).Return(apiError)
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
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{ID: 1, Role: "organizer"}),
			func() {
				s.usecase.EXPECT().UpdateBankAccount(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		s.T().Run(testCase.Name, func(t *testing.T) {
			testCase.ExpectedFunc()

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

			s.NoError(s.handler.UpdateBankAccount(ctx))
			s.Equal(testCase.ExpectedCode, rec.Code)
		})
	}
}
