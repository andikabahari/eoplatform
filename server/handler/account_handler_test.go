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

type accountHandlerSuite struct {
	suite.Suite

	ctrl    *gomock.Controller
	usecase *mu.MockAccountUsecase

	server  *server.Server
	handler *AccountHandler
}

func (s *accountHandlerSuite) SetupSuite() {
	os.Setenv("APP_ENV", "production")

	s.ctrl = gomock.NewController(s.T())
	s.usecase = mu.NewMockAccountUsecase(s.ctrl)

	conn, _ := testhelper.Mock()
	s.server = testhelper.NewServer(conn)
	s.handler = NewAccountHandler(s.usecase)
}

func (s *accountHandlerSuite) TearDownSuite() {
	s.ctrl.Finish()
}

func TestAccountHandlerSuite(t *testing.T) {
	suite.Run(t, new(accountHandlerSuite))
}

func (s *accountHandlerSuite) TestGetAccount() {
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
			"not found",
			"/v1/account",
			http.MethodGet,
			nil,
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{}),
			func() {
				apiError := helper.NewAPIError(http.StatusNotFound, "")
				s.usecase.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(apiError)
			},
			http.StatusNotFound,
		},
		{
			"ok",
			"/v1/account",
			http.MethodGet,
			nil,
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{ID: 1}),
			func() {
				s.usecase.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Return(nil)
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

			s.NoError(s.handler.GetAccount(ctx))
			s.Equal(testCase.ExpectedCode, rec.Code)
		})
	}
}

func (s *accountHandlerSuite) TestUpdateAccount() {
	testCases := []struct {
		Name         string
		Endpoint     string
		Method       string
		Body         *request.UpdateUserRequest
		Token        *jwt.Token
		ExpectedFunc func()
		ExpectedCode int
	}{
		{
			"bad request",
			"/v1/account",
			http.MethodPost,
			nil,
			nil,
			func() {
			},
			http.StatusBadRequest,
		},
		{
			"not found",
			"/v1/account",
			http.MethodPost,
			&request.UpdateUserRequest{
				Name: "User",
			},
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{}),
			func() {
				apiError := helper.NewAPIError(http.StatusNotFound, "")
				s.usecase.EXPECT().UpdateAccount(gomock.Any(), gomock.Any(), gomock.Any()).Return(apiError)
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
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{ID: 1}),
			func() {
				s.usecase.EXPECT().UpdateAccount(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
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

			s.NoError(s.handler.UpdateAccount(ctx))
			s.Equal(testCase.ExpectedCode, rec.Code)
		})
	}
}

func (s *accountHandlerSuite) TestResetPassword() {
	testCases := []struct {
		Name         string
		Endpoint     string
		Method       string
		Body         *request.UpdateUserPasswordRequest
		Token        *jwt.Token
		ExpectedFunc func()
		ExpectedCode int
	}{
		{
			"bad request",
			"/v1/account/password",
			http.MethodPost,
			nil,
			nil,
			func() {},
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
			nil,
			func() {},
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
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{}),
			func() {
				apiError := helper.NewAPIError(http.StatusNotFound, "")
				s.usecase.EXPECT().ResetPassword(gomock.Any(), gomock.Any(), gomock.Any()).Return(apiError)
			},
			http.StatusNotFound,
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
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{ID: 1}),
			func() {
				s.usecase.EXPECT().ResetPassword(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
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

			s.NoError(s.handler.ResetPassword(ctx))
			s.Equal(testCase.ExpectedCode, rec.Code)
		})
	}
}
