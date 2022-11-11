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

type serviceHandlerSuite struct {
	suite.Suite

	ctrl    *gomock.Controller
	usecase *mu.MockServiceUsecase

	server  *server.Server
	handler *ServiceHandler
}

func (s *serviceHandlerSuite) SetupSuite() {
	os.Setenv("APP_ENV", "production")

	s.ctrl = gomock.NewController(s.T())
	s.usecase = mu.NewMockServiceUsecase(s.ctrl)

	conn, _ := testhelper.Mock()
	s.server = testhelper.NewServer(conn)
	s.handler = NewServiceHandler(s.usecase)
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
				s.usecase.EXPECT().GetServices(gomock.Any(), gomock.Any())
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
				apiError := helper.NewAPIError(http.StatusNotFound, "")
				s.usecase.EXPECT().FindService(gomock.Any(), gomock.Any()).Return(apiError)
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
				s.usecase.EXPECT().FindService(gomock.Any(), gomock.Any()).Return(nil)
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
				s.usecase.EXPECT().CreateService(gomock.Any(), gomock.Any(), gomock.Any())
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
				apiError := helper.NewAPIError(http.StatusNotFound, "")
				s.usecase.EXPECT().UpdateService(gomock.Any(), gomock.Any(), gomock.Any()).Return(apiError)
			},
			nil,
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
				s.usecase.EXPECT().UpdateService(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
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
				apiError := helper.NewAPIError(http.StatusNotFound, "")
				s.usecase.EXPECT().DeleteService(gomock.Any(), gomock.Any()).Return(apiError)
			},
			nil,
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
				s.usecase.EXPECT().DeleteService(gomock.Any(), gomock.Any()).Return(nil)
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
