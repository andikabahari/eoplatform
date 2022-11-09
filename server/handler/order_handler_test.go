package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/andikabahari/eoplatform/helper"
	"github.com/andikabahari/eoplatform/model"
	"github.com/andikabahari/eoplatform/repository/mock_repository"
	"github.com/andikabahari/eoplatform/request"
	"github.com/andikabahari/eoplatform/server"
	"github.com/andikabahari/eoplatform/testhelper"
	"github.com/golang-jwt/jwt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type orderHandlerSuite struct {
	suite.Suite

	ctrl                  *gomock.Controller
	orderRepository       *mock_repository.MockOrderRepository
	paymentRepository     *mock_repository.MockPaymentRepository
	userRepository        *mock_repository.MockUserRepository
	serviceRepository     *mock_repository.MockServiceRepository
	bankAccountRepository *mock_repository.MockBankAccountRepository

	server  *server.Server
	handler *OrderHandler
}

func (s *orderHandlerSuite) SetupSuite() {
	os.Setenv("APP_ENV", "production")

	s.ctrl = gomock.NewController(s.T())
	s.orderRepository = mock_repository.NewMockOrderRepository(s.ctrl)
	s.paymentRepository = mock_repository.NewMockPaymentRepository(s.ctrl)
	s.userRepository = mock_repository.NewMockUserRepository(s.ctrl)
	s.serviceRepository = mock_repository.NewMockServiceRepository(s.ctrl)
	s.bankAccountRepository = mock_repository.NewMockBankAccountRepository(s.ctrl)

	conn, _ := testhelper.Mock()
	s.server = testhelper.NewServer(conn)
	s.handler = NewOrderHandler(
		s.server,
		s.orderRepository,
		s.paymentRepository,
		s.userRepository,
		s.serviceRepository,
		s.bankAccountRepository,
	)
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
				s.orderRepository.EXPECT().GetOrdersForCustomer(
					gomock.Eq(&[]model.Order{}),
					gomock.Eq(uint(1)),
				).SetArg(0, []model.Order{{Model: gorm.Model{ID: 1}}})

				s.paymentRepository.EXPECT().FindOnlyByOrderID(
					gomock.Eq(&model.Payment{}),
					gomock.Eq(uint(1)),
				)
			},
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{ID: 1, Role: "customer"}),
		},
		{
			"ok",
			"/v1/orders",
			nil,
			http.MethodGet,
			nil,
			http.StatusOK,
			func() {
				s.orderRepository.EXPECT().GetOrdersForOrganizer(
					gomock.Eq(&[]model.Order{}),
					gomock.Eq(uint(1)),
				)
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
				s.serviceRepository.EXPECT().Find(
					gomock.Eq(&model.Service{}),
					gomock.Eq("1"),
				)
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
				s.serviceRepository.EXPECT().Find(
					gomock.Eq(&model.Service{}),
					gomock.Eq("1"),
				).SetArg(0, model.Service{Model: gorm.Model{ID: 1}})

				s.userRepository.EXPECT().Find(
					gomock.Eq(&model.User{}),
					gomock.Eq(uint(1)),
				)

				s.orderRepository.EXPECT().Create(gomock.Any())
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
				s.orderRepository.EXPECT().Find(
					gomock.Eq(&model.Order{}),
					gomock.Eq(""),
				)
			},
			nil,
		},
		{
			"unauthorized",
			"/v1/orders/:id/accept",
			nil,
			http.MethodPost,
			nil,
			http.StatusUnauthorized,
			func() {
				s.orderRepository.EXPECT().Find(
					gomock.Eq(&model.Order{}),
					gomock.Eq(""),
				).SetArg(0, model.Order{
					Model: gorm.Model{ID: 1},
					Services: []model.Service{
						{
							Model:  gorm.Model{ID: 1},
							UserID: 1,
						},
					},
				})
			},
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{ID: 2, Role: "organizer"}),
		},
		// {
		// 	"ok",
		// 	"/v1/orders/:id/accept",
		// 	nil,
		// 	http.MethodPost,
		// 	nil,
		// 	http.StatusOK,
		// 	func() {},
		// 	jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{ID: 1, Role: "organizer"}),
		// },
		{
			"ok",
			"/v1/orders/:id/complete",
			nil,
			http.MethodPost,
			nil,
			http.StatusOK,
			func() {
				s.orderRepository.EXPECT().Find(
					gomock.Eq(&model.Order{}),
					gomock.Eq(""),
				).SetArg(0, model.Order{
					Model: gorm.Model{ID: 1},
					Services: []model.Service{
						{
							Model:  gorm.Model{ID: 1},
							UserID: 1,
						},
					},
				})
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
				s.orderRepository.EXPECT().Find(
					gomock.Eq(&model.Order{}),
					gomock.Eq(""),
				)
			},
			nil,
		},
		{
			"unauthorized",
			"/v1/orders/:id/cancel",
			nil,
			http.MethodPost,
			nil,
			http.StatusUnauthorized,
			func() {
				s.orderRepository.EXPECT().Find(
					gomock.Eq(&model.Order{}),
					gomock.Eq(""),
				).SetArg(0, model.Order{Model: gorm.Model{ID: 1}, UserID: 2})
			},
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{ID: 1}),
		},
		{
			"ok",
			"/v1/orders/:id/cancel",
			nil,
			http.MethodPost,
			nil,
			http.StatusOK,
			func() {
				s.orderRepository.EXPECT().Find(
					gomock.Eq(&model.Order{}),
					gomock.Eq(""),
				).SetArg(0, model.Order{Model: gorm.Model{ID: 1}, UserID: 1})

				s.orderRepository.EXPECT().Delete(gomock.Any())
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
				s.paymentRepository.EXPECT().FindOnlyByOrderID(
					gomock.Eq(&model.Payment{}),
					gomock.Eq("1"),
				)
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
				s.paymentRepository.EXPECT().FindOnlyByOrderID(
					gomock.Eq(&model.Payment{}),
					gomock.Eq("1"),
				).SetArg(0, model.Payment{Model: gorm.Model{ID: 1}, OrderID: 1})

				s.paymentRepository.EXPECT().Update(gomock.Any(), gomock.Any())
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
				Status:  "deny",
			},
			http.StatusOK,
			func() {
				s.paymentRepository.EXPECT().FindOnlyByOrderID(
					gomock.Eq(&model.Payment{}),
					gomock.Eq("1"),
				).SetArg(0, model.Payment{Model: gorm.Model{ID: 1}, OrderID: 1})

				s.orderRepository.EXPECT().FindOnly(
					gomock.Eq(&model.Order{}),
					gomock.Eq("1"),
				).SetArg(0, model.Order{Model: gorm.Model{ID: 1}})

				s.orderRepository.EXPECT().Delete(gomock.Any())

				s.paymentRepository.EXPECT().Update(gomock.Any(), gomock.Any())
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
