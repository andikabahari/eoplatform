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

type orderHandlerSuite struct {
	suite.Suite

	ctrl    *gomock.Controller
	usecase *mu.MockOrderUsecase

	server  *server.Server
	handler *OrderHandler
}

func (s *orderHandlerSuite) SetupSuite() {
	os.Setenv("APP_ENV", "production")

	s.ctrl = gomock.NewController(s.T())
	s.usecase = mu.NewMockOrderUsecase(s.ctrl)

	conn, _ := testhelper.Mock()
	s.server = testhelper.NewServer(conn)
	s.handler = NewOrderHandler(s.usecase)
}

func (s *orderHandlerSuite) TearDownSuite() {
	s.ctrl.Finish()
}

func TestOrderHandlerSuite(t *testing.T) {
	suite.Run(t, new(orderHandlerSuite))
}

func (s *orderHandlerSuite) TestGetOrders() {
	testCases := []struct {
		Name         string
		Endpoint     string
		PathParam    *testhelper.PathParam
		Method       string
		Body         any
		ExpectedCode int
		ExpectedFunc func()
		Token        *jwt.Token
	}{
		{
			"ok",
			"/v1/orders",
			nil,
			http.MethodGet,
			nil,
			http.StatusOK,
			func() {
				s.usecase.EXPECT().GetOrders(gomock.Any(), gomock.Any(), gomock.Any())
			},
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{ID: 1, Role: "customer"}),
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
			if testCase.PathParam != nil {
				ctx.SetParamNames(testCase.PathParam.Names...)
				ctx.SetParamValues(testCase.PathParam.Values...)
			}

			s.NoError(s.handler.GetOrders(ctx))
			s.Equal(testCase.ExpectedCode, rec.Code)
		})
	}
}

func (s *orderHandlerSuite) TestCreateOrder() {
	testCases := []struct {
		Name         string
		Endpoint     string
		PathParam    *testhelper.PathParam
		Method       string
		Body         *request.CreateOrderRequest
		ExpectedCode int
		ExpectedFunc func()
		Token        *jwt.Token
	}{
		{
			"unauthorized",
			"/v1/orders",
			nil,
			http.MethodPost,
			nil,
			http.StatusUnauthorized,
			func() {},
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{ID: 1, Role: "organizer"}),
		},
		{
			"bad request",
			"/v1/orders",
			nil,
			http.MethodPost,
			nil,
			http.StatusBadRequest,
			func() {},
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{ID: 1, Role: "customer"}),
		},
		{
			"bad request",
			"/v1/orders",
			nil,
			http.MethodPost,
			&request.CreateOrderRequest{
				DateOfEvent: "2022-12-12",
				FirstName:   "Example",
				LastName:    "User",
				Phone:       "08123456789",
				Email:       "user@example.com",
				Address:     "Mars",
				Note:        "Ok.",
				ServiceIDs:  []uint{1},
			},
			http.StatusBadRequest,
			func() {
				apiError := helper.NewAPIError(http.StatusBadRequest, "")
				s.usecase.EXPECT().CreateOrder(gomock.Any(), gomock.Any(), gomock.Any()).Return(apiError)
			},
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{ID: 1, Role: "customer"}),
		},
		{
			"ok",
			"/v1/orders",
			nil,
			http.MethodPost,
			&request.CreateOrderRequest{
				DateOfEvent: "2022-12-12",
				FirstName:   "Example",
				LastName:    "User",
				Phone:       "08123456789",
				Email:       "user@example.com",
				Address:     "Mars",
				Note:        "Ok.",
				ServiceIDs:  []uint{1},
			},
			http.StatusOK,
			func() {
				s.usecase.EXPECT().CreateOrder(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{ID: 1, Role: "customer"}),
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
			if testCase.PathParam != nil {
				ctx.SetParamNames(testCase.PathParam.Names...)
				ctx.SetParamValues(testCase.PathParam.Values...)
			}

			s.NoError(s.handler.CreateOrder(ctx))
			s.Equal(testCase.ExpectedCode, rec.Code)
		})
	}
}

func (s *orderHandlerSuite) TestAcceptOrCompleteOrder() {
	testCases := []struct {
		Name         string
		Endpoint     string
		PathParam    *testhelper.PathParam
		Method       string
		Body         any
		ExpectedCode int
		ExpectedFunc func()
		Token        *jwt.Token
	}{
		{
			"not found",
			"/v1/orders/:id/accept",
			nil,
			http.MethodPost,
			nil,
			http.StatusNotFound,
			func() {
				apiError := helper.NewAPIError(http.StatusNotFound, "")
				s.usecase.EXPECT().AcceptOrCompleteOrder(gomock.Any(), gomock.Any()).Return(apiError)
			},
			nil,
		},
		{
			"ok",
			"/v1/orders/:id/complete",
			nil,
			http.MethodPost,
			nil,
			http.StatusOK,
			func() {
				s.usecase.EXPECT().AcceptOrCompleteOrder(gomock.Any(), gomock.Any()).Return(nil)
			},
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{ID: 1, Role: "organizer"}),
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
			ctx.SetPath(testCase.Endpoint)
			if testCase.PathParam != nil {
				ctx.SetParamNames(testCase.PathParam.Names...)
				ctx.SetParamValues(testCase.PathParam.Values...)
			}

			s.NoError(s.handler.AcceptOrCompleteOrder(ctx))
			s.Equal(testCase.ExpectedCode, rec.Code)
		})
	}
}

func (s *orderHandlerSuite) TestCancelOrder() {
	testCases := []struct {
		Name         string
		Endpoint     string
		PathParam    *testhelper.PathParam
		Method       string
		Body         any
		ExpectedCode int
		ExpectedFunc func()
		Token        *jwt.Token
	}{
		{
			"not found",
			"/v1/orders/:id/cancel",
			nil,
			http.MethodPost,
			nil,
			http.StatusNotFound,
			func() {
				apiError := helper.NewAPIError(http.StatusNotFound, "")
				s.usecase.EXPECT().CancelOrder(gomock.Any(), gomock.Any()).Return(apiError)
			},
			nil,
		},
		{
			"ok",
			"/v1/orders/:id/cancel",
			nil,
			http.MethodPost,
			nil,
			http.StatusOK,
			func() {
				s.usecase.EXPECT().CancelOrder(gomock.Any(), gomock.Any()).Return(nil)
			},
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{ID: 1}),
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
			ctx.SetPath(testCase.Endpoint)
			if testCase.PathParam != nil {
				ctx.SetParamNames(testCase.PathParam.Names...)
				ctx.SetParamValues(testCase.PathParam.Values...)
			}

			s.NoError(s.handler.CancelOrder(ctx))
			s.Equal(testCase.ExpectedCode, rec.Code)
		})
	}
}

func (s *orderHandlerSuite) TestPaymentStatus() {
	testCases := []struct {
		Name         string
		Endpoint     string
		PathParam    *testhelper.PathParam
		Method       string
		Body         *request.MidtransTransactionNotificationRequest
		ExpectedCode int
		ExpectedFunc func()
		Token        *jwt.Token
	}{
		{
			"not found",
			"/v1/webhook",
			nil,
			http.MethodPost,
			&request.MidtransTransactionNotificationRequest{
				OrderID: "EOP-1",
				Status:  "",
			},
			http.StatusNotFound,
			func() {
				apiError := helper.NewAPIError(http.StatusNotFound, "")
				s.usecase.EXPECT().PaymentStatus(gomock.Any()).Return(apiError)
			},
			nil,
		},
		{
			"ok",
			"/v1/webhook",
			nil,
			http.MethodPost,
			&request.MidtransTransactionNotificationRequest{
				OrderID: "EOP-1",
				Status:  "settlement",
			},
			http.StatusOK,
			func() {
				s.usecase.EXPECT().PaymentStatus(gomock.Any()).Return(nil)
			},
			nil,
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
			ctx.SetPath(testCase.Endpoint)
			if testCase.PathParam != nil {
				ctx.SetParamNames(testCase.PathParam.Names...)
				ctx.SetParamValues(testCase.PathParam.Values...)
			}

			s.NoError(s.handler.PaymentStatus(ctx))
			s.Equal(testCase.ExpectedCode, rec.Code)
		})
	}
}
