package test

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
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

type serviceSuite struct {
	suite.Suite
	mock    sqlmock.Sqlmock
	server  *server.Server
	handler *handler.ServiceHandler
}

func (s *serviceSuite) SetupSuite() {
	os.Setenv("APP_ENV", "production")

	var conn *sql.DB
	conn, s.mock = testhelper.Mock()
	s.server = testhelper.NewServer(conn)
	s.handler = handler.NewServiceHandler(s.server)
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(serviceSuite))
}

func (s *serviceSuite) TestGetServices() {
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
			"/v1/services",
			nil,
			http.MethodGet,
			nil,
			http.StatusOK,
			nil,
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

			s.NoError(s.handler.GetServices(ctx))
			s.Equal(testCase.ExpectedCode, rec.Code)
		})
	}
}

func (s *serviceSuite) TestFindService() {
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
			"/v1/services",
			&pathParam{
				Names:  []string{"id"},
				Values: []string{"1"},
			},
			http.MethodGet,
			nil,
			http.StatusNotFound,
			nil,
			[]query{
				{
					Raw:  regexp.QuoteMeta("SELECT * FROM `services` WHERE id = ? AND `services`.`deleted_at` IS NULL"),
					Rows: sqlmock.NewRows(nil),
				},
			},
		},
		{
			"ok",
			"/v1/services",
			&pathParam{
				Names:  []string{"id"},
				Values: []string{"1"},
			},
			http.MethodGet,
			nil,
			http.StatusOK,
			nil,
			[]query{
				{
					Raw:  regexp.QuoteMeta("SELECT * FROM `services` WHERE id = ? AND `services`.`deleted_at` IS NULL"),
					Rows: sqlmock.NewRows([]string{"id"}).AddRow(1),
				},
			},
		},
	}

	for _, testCase := range testCases {
		s.T().Run(testCase.Name, func(t *testing.T) {
			for _, query := range testCase.Queries {
				s.mock.ExpectQuery(query.Raw).WithArgs(query.Args...).WillReturnRows(query.Rows)
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

			s.NoError(s.handler.FindService(ctx))
			s.Equal(testCase.ExpectedCode, rec.Code)
		})
	}
}

func (s *serviceSuite) TestCreateService() {
	testCases := []struct {
		Name         string
		Endpoint     string
		PathParam    *pathParam
		Method       string
		Body         *request.CreateServiceRequest
		ExpectedCode int
		Token        *jwt.Token
		Queries      []query
	}{
		{
			"unauthorized",
			"/v1/services",
			nil,
			http.MethodPost,
			nil,
			http.StatusUnauthorized,
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{}),
			nil,
		},
		{
			"bad request",
			"/v1/services",
			nil,
			http.MethodPost,
			&request.CreateServiceRequest{},
			http.StatusBadRequest,
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{
				ID:   1,
				Role: "organizer",
			}),
			nil,
		},
		{
			"ok",
			"/v1/services",
			nil,
			http.MethodPost,
			&request.CreateServiceRequest{
				BasicService: request.BasicService{
					Name:        "Service",
					Cost:        1000000,
					Phone:       "08123456789",
					Email:       "user@example.com",
					Description: "Lorem ipsum",
				},
			},
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
				s.mock.ExpectQuery(query.Raw).WithArgs(query.Args...).WillReturnRows(query.Rows)
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

			s.NoError(s.handler.CreateService(ctx))
			s.Equal(testCase.ExpectedCode, rec.Code)
		})
	}
}

func (s *serviceSuite) TestUpdateService() {
	testCases := []struct {
		Name         string
		Endpoint     string
		PathParam    *pathParam
		Method       string
		Body         *request.UpdateServiceRequest
		ExpectedCode int
		Token        *jwt.Token
		Queries      []query
	}{
		{
			"bad request",
			"/v1/services/:id",
			&pathParam{
				Names:  []string{"id"},
				Values: []string{"1"},
			},
			http.MethodPost,
			nil,
			http.StatusBadRequest,
			nil,
			nil,
		},
		{
			"not found",
			"/v1/services/:id",
			&pathParam{
				Names:  []string{"id"},
				Values: []string{"1"},
			},
			http.MethodPost,
			&request.UpdateServiceRequest{
				BasicService: request.BasicService{
					Name:        "Service",
					Cost:        1000000,
					Phone:       "08123456789",
					Email:       "user@example.com",
					Description: "Lorem ipsum",
				},
			},
			http.StatusNotFound,
			nil,
			nil,
		},
		{
			"unauthorized",
			"/v1/services/:id",
			&pathParam{
				Names:  []string{"id"},
				Values: []string{"1"},
			},
			http.MethodPost,
			&request.UpdateServiceRequest{
				BasicService: request.BasicService{
					Name:        "Service",
					Cost:        1000000,
					Phone:       "08123456789",
					Email:       "user@example.com",
					Description: "Lorem ipsum",
				},
			},
			http.StatusUnauthorized,
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{ID: 1}),
			[]query{
				{
					Raw:  regexp.QuoteMeta("SELECT * FROM `services` WHERE id = ? AND `services`.`deleted_at` IS NULL"),
					Args: []driver.Value{"1"},
					Rows: sqlmock.NewRows([]string{"id", "user_id"}).AddRow(1, 2),
				},
			},
		},
		{
			"ok",
			"/v1/services/:id",
			&pathParam{
				Names:  []string{"id"},
				Values: []string{"1"},
			},
			http.MethodPost,
			&request.UpdateServiceRequest{
				BasicService: request.BasicService{
					Name:        "Service",
					Cost:        1000000,
					Phone:       "08123456789",
					Email:       "user@example.com",
					Description: "Lorem ipsum",
				},
			},
			http.StatusOK,
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{ID: 2}),
			[]query{
				{
					Raw:  regexp.QuoteMeta("SELECT * FROM `services` WHERE id = ? AND `services`.`deleted_at` IS NULL"),
					Args: []driver.Value{"1"},
					Rows: sqlmock.NewRows([]string{"id", "user_id"}).AddRow(1, 2),
				},
			},
		},
	}

	for _, testCase := range testCases {
		s.T().Run(testCase.Name, func(t *testing.T) {
			for _, query := range testCase.Queries {
				s.mock.ExpectQuery(query.Raw).WithArgs(query.Args...).WillReturnRows(query.Rows)
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

			s.NoError(s.handler.UpdateService(ctx))
			s.Equal(testCase.ExpectedCode, rec.Code)
		})
	}
}

func (s *serviceSuite) TestDeleteService() {
	testCases := []struct {
		Name         string
		Endpoint     string
		PathParam    *pathParam
		Method       string
		Body         *request.UpdateServiceRequest
		ExpectedCode int
		Token        *jwt.Token
		Queries      []query
	}{
		{
			"not found",
			"/v1/services/:id",
			&pathParam{
				Names:  []string{"id"},
				Values: []string{"1"},
			},
			http.MethodDelete,
			nil,
			http.StatusNotFound,
			nil,
			nil,
		},
		{
			"unauthorized",
			"/v1/services/:id",
			&pathParam{
				Names:  []string{"id"},
				Values: []string{"1"},
			},
			http.MethodDelete,
			nil,
			http.StatusUnauthorized,
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{ID: 1}),
			[]query{
				{
					Raw:  regexp.QuoteMeta("SELECT * FROM `services` WHERE id = ? AND `services`.`deleted_at` IS NULL"),
					Args: []driver.Value{"1"},
					Rows: sqlmock.NewRows([]string{"id", "user_id"}).AddRow(1, 2),
				},
			},
		},
		{
			"ok",
			"/v1/services/:id",
			&pathParam{
				Names:  []string{"id"},
				Values: []string{"1"},
			},
			http.MethodDelete,
			nil,
			http.StatusOK,
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{ID: 1}),
			[]query{
				{
					Raw:  regexp.QuoteMeta("SELECT * FROM `services` WHERE id = ? AND `services`.`deleted_at` IS NULL"),
					Args: []driver.Value{"1"},
					Rows: sqlmock.NewRows([]string{"id", "user_id"}).AddRow(1, 1),
				},
			},
		},
	}

	for _, testCase := range testCases {
		s.T().Run(testCase.Name, func(t *testing.T) {
			for _, query := range testCase.Queries {
				s.mock.ExpectQuery(query.Raw).WithArgs(query.Args...).WillReturnRows(query.Rows)
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

			s.NoError(s.handler.DeleteService(ctx))
			s.Equal(testCase.ExpectedCode, rec.Code)
		})
	}
}
