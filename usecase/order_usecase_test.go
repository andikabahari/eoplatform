package usecase

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/andikabahari/eoplatform/helper"
	"github.com/andikabahari/eoplatform/model"
	mr "github.com/andikabahari/eoplatform/repository/mock_repository"
	"github.com/andikabahari/eoplatform/request"
	"github.com/golang-jwt/jwt"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type orderUsecaseSuite struct {
	suite.Suite

	ctrl                  *gomock.Controller
	orderRepository       *mr.MockOrderRepository
	paymentRepository     *mr.MockPaymentRepository
	userRepository        *mr.MockUserRepository
	serviceRepository     *mr.MockServiceRepository
	bankAccountRepository *mr.MockBankAccountRepository

	usecase OrderUsecase
}

func (s *orderUsecaseSuite) SetupSuite() {
	os.Setenv("APP_ENV", "production")

	s.ctrl = gomock.NewController(s.T())
	s.orderRepository = mr.NewMockOrderRepository(s.ctrl)
	s.paymentRepository = mr.NewMockPaymentRepository(s.ctrl)
	s.userRepository = mr.NewMockUserRepository(s.ctrl)
	s.serviceRepository = mr.NewMockServiceRepository(s.ctrl)
	s.bankAccountRepository = mr.NewMockBankAccountRepository(s.ctrl)

	s.usecase = NewOrderUsecase(
		s.orderRepository,
		s.paymentRepository,
		s.userRepository,
		s.serviceRepository,
		s.bankAccountRepository,
	)
}

func (s *orderUsecaseSuite) TearDownSuite() {
	s.ctrl.Finish()
}

func TestOrderUsecaseSuite(t *testing.T) {
	suite.Run(t, new(orderUsecaseSuite))
}

func (s *orderUsecaseSuite) TestGetOrders() {
	testCases := []struct {
		Name         string
		Body         any
		Claims       *helper.JWTCustomClaims
		ExpectedFunc func()
		ExpectedCode int
	}{
		{
			"ok",
			nil,
			&helper.JWTCustomClaims{ID: 1, Role: "customer"},
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
			http.StatusOK,
		},
		{
			"ok",
			nil,
			&helper.JWTCustomClaims{ID: 1, Role: "organizer"},
			func() {
				s.orderRepository.EXPECT().GetOrdersForOrganizer(
					gomock.Eq(&[]model.Order{}),
					gomock.Eq(uint(1)),
				)
			},
			http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		s.T().Run(testCase.Name, func(t *testing.T) {
			testCase.ExpectedFunc()
			s.usecase.GetOrders(testCase.Claims, &[]model.Order{}, &[]model.Payment{})
		})
	}
}

func (s *orderUsecaseSuite) TestCreateOrder() {
	testCases := []struct {
		Name         string
		Body         *request.CreateOrderRequest
		Claims       *helper.JWTCustomClaims
		ExpectedFunc func()
		ExpectedCode int
	}{
		{
			"bad request",
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
			&helper.JWTCustomClaims{ID: 1, Role: "customer"},
			func() {
				s.serviceRepository.EXPECT().Find(
					gomock.Eq(&model.Service{}),
					gomock.Eq("1"),
				)
			},
			http.StatusBadRequest,
		},
		{
			"ok",
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
			&helper.JWTCustomClaims{ID: 1, Role: "customer"},
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
			http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		s.T().Run(testCase.Name, func(t *testing.T) {
			testCase.ExpectedFunc()
			if apiError := s.usecase.CreateOrder(testCase.Claims, &model.Order{}, testCase.Body); apiError != nil {
				code, _ := apiError.APIError()
				s.Equal(testCase.ExpectedCode, code)
			}
		})
	}
}

func (s *orderUsecaseSuite) TestAcceptOrCompleteOrder() {
	createContext := func(token *jwt.Token, endpoint string) echo.Context {
		req := httptest.NewRequest("", "/", nil)
		rec := httptest.NewRecorder()
		ctx := echo.New().NewContext(req, rec)
		ctx.SetPath(endpoint)
		ctx.Set("user", token)
		ctx.SetParamNames("id")
		ctx.SetParamValues("1")
		return ctx
	}

	testCases := []struct {
		Name         string
		Body         any
		Context      echo.Context
		ExpectedFunc func()
		ExpectedCode int
	}{
		{
			"not found",
			nil,
			createContext(nil, "/v1/orders/:id/accept"),
			func() {
				s.orderRepository.EXPECT().Find(
					gomock.Eq(&model.Order{}),
					gomock.Eq("1"),
				)
			},
			http.StatusNotFound,
		},
		{
			"unauthorized",
			nil,
			createContext(jwt.NewWithClaims(
				jwt.SigningMethodHS256,
				&helper.JWTCustomClaims{ID: 2, Role: "organizer"},
			), "/v1/orders/:id/accept"),
			func() {
				s.orderRepository.EXPECT().Find(
					gomock.Eq(&model.Order{}),
					gomock.Eq("1"),
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
			http.StatusUnauthorized,
		},
		{
			"ok",
			nil,
			createContext(jwt.NewWithClaims(
				jwt.SigningMethodHS256,
				&helper.JWTCustomClaims{ID: 1, Role: "organizer"},
			), "/v1/orders/:id/complete"),
			func() {
				s.orderRepository.EXPECT().Find(
					gomock.Eq(&model.Order{}),
					gomock.Eq("1"),
				).SetArg(0, model.Order{
					Model: gorm.Model{ID: 1},
					Services: []model.Service{
						{
							Model:  gorm.Model{ID: 1},
							UserID: 1,
						},
					},
				})

				s.orderRepository.EXPECT().Save(gomock.Any())
			},
			http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		s.T().Run(testCase.Name, func(t *testing.T) {
			testCase.ExpectedFunc()
			if apiError := s.usecase.AcceptOrCompleteOrder(testCase.Context, &model.Order{}); apiError != nil {
				code, _ := apiError.APIError()
				s.Equal(testCase.ExpectedCode, code)
			}
		})
	}
}

func (s *orderUsecaseSuite) TestCancelOrder() {
	createContext := func(token *jwt.Token, id string) echo.Context {
		req := httptest.NewRequest("", "/", nil)
		rec := httptest.NewRecorder()
		ctx := echo.New().NewContext(req, rec)
		ctx.Set("user", token)
		ctx.SetParamNames("id")
		ctx.SetParamValues(id)
		return ctx
	}

	testCases := []struct {
		Name         string
		Body         any
		Context      echo.Context
		ExpectedFunc func()
		ExpectedCode int
	}{
		{
			"not found",
			nil,
			createContext(nil, "1"),
			func() {
				s.orderRepository.EXPECT().Find(
					gomock.Eq(&model.Order{}),
					gomock.Eq("1"),
				)
			},
			http.StatusNotFound,
		},
		{
			"unauthorized",
			nil,
			createContext(jwt.NewWithClaims(
				jwt.SigningMethodHS256,
				&helper.JWTCustomClaims{ID: 1},
			), "1"),
			func() {
				s.orderRepository.EXPECT().Find(
					gomock.Eq(&model.Order{}),
					gomock.Eq("1"),
				).SetArg(0, model.Order{Model: gorm.Model{ID: 1}, UserID: 2})
			},
			http.StatusUnauthorized,
		},
		{
			"ok",
			nil,
			createContext(jwt.NewWithClaims(
				jwt.SigningMethodHS256,
				&helper.JWTCustomClaims{ID: 1},
			), "1"),
			func() {
				s.orderRepository.EXPECT().Find(
					gomock.Eq(&model.Order{}),
					gomock.Eq("1"),
				).SetArg(0, model.Order{Model: gorm.Model{ID: 1}, UserID: 1})

				s.orderRepository.EXPECT().Delete(gomock.Any())
			},
			http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		s.T().Run(testCase.Name, func(t *testing.T) {
			testCase.ExpectedFunc()
			if apiError := s.usecase.CancelOrder(testCase.Context, &model.Order{}); apiError != nil {
				code, _ := apiError.APIError()
				s.Equal(testCase.ExpectedCode, code)
			}
		})
	}
}

func (s *orderUsecaseSuite) TestPaymentStatus() {
	testCases := []struct {
		Name         string
		Body         *request.MidtransTransactionNotificationRequest
		ExpectedFunc func()
		ExpectedCode int
	}{
		{
			"not found",
			&request.MidtransTransactionNotificationRequest{
				OrderID: "EOP-1",
				Status:  "",
			},
			func() {
				s.paymentRepository.EXPECT().FindOnlyByOrderID(
					gomock.Eq(&model.Payment{}),
					gomock.Eq("1"),
				)
			},
			http.StatusNotFound,
		},
		{
			"ok",
			&request.MidtransTransactionNotificationRequest{
				OrderID: "EOP-1",
				Status:  "settlement",
			},
			func() {
				s.paymentRepository.EXPECT().FindOnlyByOrderID(
					gomock.Eq(&model.Payment{}),
					gomock.Eq("1"),
				).SetArg(0, model.Payment{Model: gorm.Model{ID: 1}, OrderID: 1})

				s.paymentRepository.EXPECT().Update(gomock.Any(), gomock.Any())
			},
			http.StatusOK,
		},
		{
			"ok",
			&request.MidtransTransactionNotificationRequest{
				OrderID: "EOP-1",
				Status:  "deny",
			},
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
			http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		s.T().Run(testCase.Name, func(t *testing.T) {
			testCase.ExpectedFunc()
			if apiError := s.usecase.PaymentStatus(testCase.Body); apiError != nil {
				code, _ := apiError.APIError()
				s.Equal(testCase.ExpectedCode, code)
			}
		})
	}
}
