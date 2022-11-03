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

type orderSuite struct {
	suite.Suite
	mock    sqlmock.Sqlmock
	server  *server.Server
	handler *handler.OrderHandler
}

func (s *orderSuite) SetupSuite() {
	os.Setenv("APP_ENV", "production")

	var conn *sql.DB
	conn, s.mock = testhelper.Mock()
	s.server = testhelper.NewServer(conn)
	s.handler = handler.NewOrderHandler(s.server)
}

func TestOrderSuite(t *testing.T) {
	suite.Run(t, new(orderSuite))
}

func (s *orderSuite) TestGetOrders() {
	testCases := []struct {
		Name         string
		Endpoint     string
		PathParam    *pathParam
		Method       string
		Body         any
		ExpectedCode int
		Token        *jwt.Token
		Queries      []query
	}{
		{
			"ok",
			"/v1/orders",
			nil,
			http.MethodGet,
			nil,
			http.StatusOK,
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{
				ID:   1,
				Role: "customer",
			}),
			nil,
		},
		{
			"ok",
			"/v1/orders",
			nil,
			http.MethodGet,
			nil,
			http.StatusOK,
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{
				ID:   1,
				Role: "organizer",
			}),
			nil,
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
			if testCase.PathParam != nil {
				ctx.SetParamNames(testCase.PathParam.Names...)
				ctx.SetParamValues(testCase.PathParam.Values...)
			}

			s.NoError(s.handler.GetOrders(ctx))
			s.Equal(testCase.ExpectedCode, rec.Code)
		})
	}
}

func (s *orderSuite) TestCreateOrder() {
	testCases := []struct {
		Name         string
		Endpoint     string
		PathParam    *pathParam
		Method       string
		Body         *request.CreateOrderRequest
		ExpectedCode int
		Token        *jwt.Token
		Queries      []query
	}{
		{
			"unauthorized",
			"/v1/orders",
			nil,
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
			"/v1/orders",
			nil,
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
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{
				ID:   1,
				Role: "customer",
			}),
			nil,
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
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{
				ID:   1,
				Role: "customer",
			}),
			[]query{
				{
					Raw:  regexp.QuoteMeta("SELECT * FROM `services` WHERE id = ? AND `services`.`deleted_at` IS NULL"),
					Rows: sqlmock.NewRows([]string{"id", "user_id"}).AddRow(1, 2),
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
			if testCase.PathParam != nil {
				ctx.SetParamNames(testCase.PathParam.Names...)
				ctx.SetParamValues(testCase.PathParam.Values...)
			}

			s.NoError(s.handler.CreateOrder(ctx))
			s.Equal(testCase.ExpectedCode, rec.Code)
		})
	}
}

func (s *orderSuite) TestAcceptOrCompleteOrder() {
	testCases := []struct {
		Name         string
		Endpoint     string
		PathParam    *pathParam
		Method       string
		Body         any
		ExpectedCode int
		Token        *jwt.Token
		Queries      []query
	}{
		{
			"not found",
			"/v1/orders/:id/accept",
			nil,
			http.MethodPost,
			nil,
			http.StatusNotFound,
			nil,
			nil,
		},
		{
			"unauthorized",
			"/v1/orders/:id/accept",
			nil,
			http.MethodPost,
			nil,
			http.StatusUnauthorized,
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{
				ID:   2,
				Role: "organizer",
			}),
			[]query{
				{
					Raw:  regexp.QuoteMeta("SELECT * FROM `orders` WHERE id = ? AND `orders`.`deleted_at` IS NULL"),
					Rows: sqlmock.NewRows([]string{"id"}).AddRow(1),
				},
				{
					Raw:  regexp.QuoteMeta("SELECT * FROM `order_services`"),
					Rows: sqlmock.NewRows([]string{"order_id", "service_id"}).AddRow(1, 1),
				},
				{
					Raw:  regexp.QuoteMeta("SELECT * FROM `services`"),
					Rows: sqlmock.NewRows([]string{"id", "user_id"}).AddRow(1, 1),
				},
			},
		},
		{
			"ok",
			"/v1/orders/:id/accept",
			nil,
			http.MethodPost,
			nil,
			http.StatusOK,
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{
				ID:   1,
				Role: "organizer",
			}),
			[]query{
				{
					Raw:  regexp.QuoteMeta("SELECT * FROM `orders` WHERE id = ? AND `orders`.`deleted_at` IS NULL"),
					Rows: sqlmock.NewRows([]string{"id"}).AddRow(1),
				},
				{
					Raw:  regexp.QuoteMeta("SELECT * FROM `order_services`"),
					Rows: sqlmock.NewRows([]string{"order_id", "service_id"}).AddRow(1, 1),
				},
				{
					Raw:  regexp.QuoteMeta("SELECT * FROM `services`"),
					Rows: sqlmock.NewRows([]string{"id", "user_id"}).AddRow(1, 1),
				},
			},
		},
		{
			"ok",
			"/v1/orders/:id/complete",
			nil,
			http.MethodPost,
			nil,
			http.StatusOK,
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{
				ID:   1,
				Role: "organizer",
			}),
			[]query{
				{
					Raw:  regexp.QuoteMeta("SELECT * FROM `orders` WHERE id = ? AND `orders`.`deleted_at` IS NULL"),
					Rows: sqlmock.NewRows([]string{"id", "is_accepted"}).AddRow(1, 1),
				},
				{
					Raw:  regexp.QuoteMeta("SELECT * FROM `order_services`"),
					Rows: sqlmock.NewRows([]string{"order_id", "service_id"}).AddRow(1, 1),
				},
				{
					Raw:  regexp.QuoteMeta("SELECT * FROM `services`"),
					Rows: sqlmock.NewRows([]string{"id", "user_id"}).AddRow(1, 1),
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

func (s *orderSuite) TestCancelOrder() {
	testCases := []struct {
		Name         string
		Endpoint     string
		PathParam    *pathParam
		Method       string
		Body         any
		ExpectedCode int
		Token        *jwt.Token
		Queries      []query
	}{
		{
			"not found",
			"/v1/orders/:id/cancel",
			nil,
			http.MethodPost,
			nil,
			http.StatusNotFound,
			nil,
			nil,
		},
		{
			"unauthorized",
			"/v1/orders/:id/cancel",
			nil,
			http.MethodPost,
			nil,
			http.StatusUnauthorized,
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{ID: 1}),
			[]query{
				{
					Raw:  regexp.QuoteMeta("SELECT * FROM `orders` WHERE id = ? AND `orders`.`deleted_at` IS NULL"),
					Rows: sqlmock.NewRows([]string{"id", "user_id"}).AddRow(1, 2),
				},
			},
		},
		{
			"ok",
			"/v1/orders/:id/cancel",
			nil,
			http.MethodPost,
			nil,
			http.StatusOK,
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{ID: 1}),
			[]query{
				{
					Raw:  regexp.QuoteMeta("SELECT * FROM `orders` WHERE id = ? AND `orders`.`deleted_at` IS NULL"),
					Rows: sqlmock.NewRows([]string{"id", "user_id"}).AddRow(1, 1),
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

func (s *orderSuite) TestPaymentStatus() {
	testCases := []struct {
		Name         string
		Endpoint     string
		PathParam    *pathParam
		Method       string
		Body         *request.MidtransTransactionNotificationRequest
		ExpectedCode int
		Token        *jwt.Token
		Queries      []query
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
			nil,
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
			nil,
			[]query{
				{
					Raw:  regexp.QuoteMeta("SELECT * FROM `payments`"),
					Rows: sqlmock.NewRows([]string{"id", "order_id"}).AddRow(1, 1),
				},
			},
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
			nil,
			[]query{
				{
					Raw:  regexp.QuoteMeta("SELECT * FROM `payments`"),
					Rows: sqlmock.NewRows([]string{"id", "order_id"}).AddRow(1, 1),
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
