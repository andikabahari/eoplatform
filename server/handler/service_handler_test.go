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

type serviceHandlerSuite struct {
	suite.Suite

	ctrl              *gomock.Controller
	serviceRepository *mock_repository.MockServiceRepository

	server  *server.Server
	handler *ServiceHandler
}

func (s *serviceHandlerSuite) SetupSuite() {
	os.Setenv("APP_ENV", "production")

	s.ctrl = gomock.NewController(s.T())
	s.serviceRepository = mock_repository.NewMockServiceRepository(s.ctrl)

	conn, _ := testhelper.Mock()
	s.server = testhelper.NewServer(conn)
	s.handler = NewServiceHandler(s.server, s.serviceRepository)
}

func TestServiceHandlerSuite(t *testing.T) {
	suite.Run(t, new(serviceHandlerSuite))
}

func (s *serviceHandlerSuite) TestGetServices() {
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
			"/v1/services",
			nil,
			http.MethodGet,
			nil,
			http.StatusOK,
			func() {
				s.serviceRepository.EXPECT().Get(
					gomock.Eq(&[]model.Service{}),
					gomock.Eq(""),
				)
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
			if testCase.PathParam != nil {
				ctx.SetParamNames(testCase.PathParam.Names...)
				ctx.SetParamValues(testCase.PathParam.Values...)
			}

			s.NoError(s.handler.GetServices(ctx))
			s.Equal(testCase.ExpectedCode, rec.Code)
		})
	}
}

func (s *serviceHandlerSuite) TestFindService() {
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
			"/v1/services",
			&testhelper.PathParam{
				Names:  []string{"id"},
				Values: []string{"1"},
			},
			http.MethodGet,
			nil,
			http.StatusNotFound,
			func() {
				s.serviceRepository.EXPECT().Find(
					gomock.Eq(&model.Service{}),
					gomock.Eq("1"),
				)
			},
			nil,
		},
		{
			"ok",
			"/v1/services",
			&testhelper.PathParam{
				Names:  []string{"id"},
				Values: []string{"1"},
			},
			http.MethodGet,
			nil,
			http.StatusOK,
			func() {
				s.serviceRepository.EXPECT().Find(
					gomock.Eq(&model.Service{}),
					gomock.Eq("1"),
				).SetArg(0, model.Service{Model: gorm.Model{ID: 1}})
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
			if testCase.PathParam != nil {
				ctx.SetParamNames(testCase.PathParam.Names...)
				ctx.SetParamValues(testCase.PathParam.Values...)
			}

			s.NoError(s.handler.FindService(ctx))
			s.Equal(testCase.ExpectedCode, rec.Code)
		})
	}
}

func (s *serviceHandlerSuite) TestCreateService() {
	testCases := []struct {
		Name         string
		Endpoint     string
		PathParam    *testhelper.PathParam
		Method       string
		Body         *request.CreateServiceRequest
		ExpectedCode int
		ExpectedFunc func()
		Token        *jwt.Token
	}{
		{
			"unauthorized",
			"/v1/services",
			nil,
			http.MethodPost,
			nil,
			http.StatusUnauthorized,
			func() {},
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{}),
		},
		{
			"bad request",
			"/v1/services",
			nil,
			http.MethodPost,
			&request.CreateServiceRequest{},
			http.StatusBadRequest,
			func() {},
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{
				ID:   1,
				Role: "organizer",
			}),
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
			func() {
				s.serviceRepository.EXPECT().Create(gomock.Any())
			},
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{
				ID:   1,
				Role: "organizer",
			}),
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

			s.NoError(s.handler.CreateService(ctx))
			s.Equal(testCase.ExpectedCode, rec.Code)
		})
	}
}

func (s *serviceHandlerSuite) TestUpdateService() {
	testCases := []struct {
		Name         string
		Endpoint     string
		PathParam    *testhelper.PathParam
		Method       string
		Body         *request.UpdateServiceRequest
		ExpectedCode int
		ExpectedFunc func()
		Token        *jwt.Token
	}{
		{
			"bad request",
			"/v1/services/:id",
			&testhelper.PathParam{
				Names:  []string{"id"},
				Values: []string{"1"},
			},
			http.MethodPost,
			nil,
			http.StatusBadRequest,
			func() {},
			nil,
		},
		{
			"not found",
			"/v1/services/:id",
			&testhelper.PathParam{
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
			func() {
				s.serviceRepository.EXPECT().Find(
					gomock.Eq(&model.Service{}),
					gomock.Eq("1"),
				)
			},
			nil,
		},
		{
			"unauthorized",
			"/v1/services/:id",
			&testhelper.PathParam{
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
			func() {
				s.serviceRepository.EXPECT().Find(
					gomock.Eq(&model.Service{}),
					gomock.Eq("1"),
				).SetArg(0, model.Service{Model: gorm.Model{ID: 1}})
			},
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{ID: 1}),
		},
		{
			"ok",
			"/v1/services/:id",
			&testhelper.PathParam{
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
			func() {
				s.serviceRepository.EXPECT().Find(
					gomock.Eq(&model.Service{}),
					gomock.Eq("1"),
				).SetArg(0, model.Service{Model: gorm.Model{ID: 1}, UserID: 2})

				s.serviceRepository.EXPECT().Update(gomock.Any(), gomock.Any())
			},
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{ID: 2}),
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

			s.NoError(s.handler.UpdateService(ctx))
			s.Equal(testCase.ExpectedCode, rec.Code)
		})
	}
}

func (s *serviceHandlerSuite) TestDeleteService() {
	testCases := []struct {
		Name         string
		Endpoint     string
		PathParam    *testhelper.PathParam
		Method       string
		Body         *request.UpdateServiceRequest
		ExpectedCode int
		ExpectedFunc func()
		Token        *jwt.Token
	}{
		{
			"not found",
			"/v1/services/:id",
			&testhelper.PathParam{
				Names:  []string{"id"},
				Values: []string{"1"},
			},
			http.MethodDelete,
			nil,
			http.StatusNotFound,
			func() {
				s.serviceRepository.EXPECT().Find(
					gomock.Eq(&model.Service{}),
					gomock.Eq("1"),
				)
			},
			nil,
		},
		{
			"unauthorized",
			"/v1/services/:id",
			&testhelper.PathParam{
				Names:  []string{"id"},
				Values: []string{"1"},
			},
			http.MethodDelete,
			nil,
			http.StatusUnauthorized,
			func() {
				s.serviceRepository.EXPECT().Find(
					gomock.Eq(&model.Service{}),
					gomock.Eq("1"),
				).SetArg(0, model.Service{Model: gorm.Model{ID: 1}, UserID: 2})
			},
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{ID: 1}),
		},
		{
			"ok",
			"/v1/services/:id",
			&testhelper.PathParam{
				Names:  []string{"id"},
				Values: []string{"1"},
			},
			http.MethodDelete,
			nil,
			http.StatusOK,
			func() {
				s.serviceRepository.EXPECT().Find(
					gomock.Eq(&model.Service{}),
					gomock.Eq("1"),
				).SetArg(0, model.Service{Model: gorm.Model{ID: 1}, UserID: 1})

				s.serviceRepository.EXPECT().Delete(gomock.Any())
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
			if testCase.PathParam != nil {
				ctx.SetParamNames(testCase.PathParam.Names...)
				ctx.SetParamValues(testCase.PathParam.Values...)
			}

			s.NoError(s.handler.DeleteService(ctx))
			s.Equal(testCase.ExpectedCode, rec.Code)
		})
	}
}
