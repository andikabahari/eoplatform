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

type serviceUsecaseSuite struct {
	suite.Suite

	ctrl              *gomock.Controller
	serviceRepository *mr.MockServiceRepository

	usecase ServiceUsecase
}

func (s *serviceUsecaseSuite) SetupSuite() {
	os.Setenv("APP_ENV", "production")

	s.ctrl = gomock.NewController(s.T())
	s.serviceRepository = mr.NewMockServiceRepository(s.ctrl)

	s.usecase = NewServiceUsecase(s.serviceRepository)
}

func (s *serviceUsecaseSuite) TearDownSuite() {
	s.ctrl.Finish()
}

func TestServiceUsecaseSuite(t *testing.T) {
	suite.Run(t, new(serviceUsecaseSuite))
}

func (s *serviceUsecaseSuite) TestGetServices() {
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
			nil,
			func() {
				s.serviceRepository.EXPECT().Get(
					gomock.Eq(&[]model.Service{}),
					gomock.Eq(""),
				)
			},
			http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		s.T().Run(testCase.Name, func(t *testing.T) {
			testCase.ExpectedFunc()
			s.usecase.GetServices(&[]model.Service{}, "")
		})
	}
}

func (s *serviceUsecaseSuite) TestFindService() {
	testCases := []struct {
		Name         string
		Body         any
		Claims       *helper.JWTCustomClaims
		ExpectedFunc func()
		ExpectedCode int
	}{
		{
			"not found",
			nil,
			nil,
			func() {
				s.serviceRepository.EXPECT().Find(
					gomock.Eq(&model.Service{}),
					gomock.Eq("1"),
				)
			},
			http.StatusNotFound,
		},
		{
			"ok",
			nil,
			nil,
			func() {
				s.serviceRepository.EXPECT().Find(
					gomock.Eq(&model.Service{}),
					gomock.Eq("1"),
				).SetArg(0, model.Service{Model: gorm.Model{ID: 1}})
			},
			http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		s.T().Run(testCase.Name, func(t *testing.T) {
			testCase.ExpectedFunc()
			if apiError := s.usecase.FindService(&model.Service{}, "1"); apiError != nil {
				code, _ := apiError.APIError()
				s.Equal(testCase.ExpectedCode, code)
			}
		})
	}
}

func (s *serviceUsecaseSuite) TestCreateService() {
	testCases := []struct {
		Name         string
		Body         *request.CreateServiceRequest
		Claims       *helper.JWTCustomClaims
		ExpectedFunc func()
		ExpectedCode int
	}{
		{
			"ok",
			&request.CreateServiceRequest{
				BasicService: request.BasicService{
					Name:        "Service",
					Cost:        1000000,
					Phone:       "08123456789",
					Email:       "user@example.com",
					Description: "Lorem ipsum",
				},
			},
			&helper.JWTCustomClaims{ID: 1, Role: "organizer"},
			func() {
				s.serviceRepository.EXPECT().Create(gomock.Any())
			},
			http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		s.T().Run(testCase.Name, func(t *testing.T) {
			testCase.ExpectedFunc()
			s.usecase.CreateService(testCase.Claims, &model.Service{}, testCase.Body)
		})
	}
}

func (s *serviceUsecaseSuite) TestUpdateService() {
	createContext := func(token *jwt.Token) echo.Context {
		req := httptest.NewRequest("", "/", nil)
		rec := httptest.NewRecorder()
		ctx := echo.New().NewContext(req, rec)
		ctx.Set("user", token)
		ctx.SetParamNames("id")
		ctx.SetParamValues("1")
		return ctx
	}

	testCases := []struct {
		Name         string
		Body         *request.UpdateServiceRequest
		Context      echo.Context
		ExpectedFunc func()
		ExpectedCode int
	}{
		{
			"not found",
			&request.UpdateServiceRequest{
				BasicService: request.BasicService{
					Name:        "Service",
					Cost:        1000000,
					Phone:       "08123456789",
					Email:       "user@example.com",
					Description: "Lorem ipsum",
				},
			},
			createContext(jwt.NewWithClaims(
				jwt.SigningMethodHS256,
				&helper.JWTCustomClaims{ID: 1},
			)),
			func() {
				s.serviceRepository.EXPECT().Find(
					gomock.Eq(&model.Service{}),
					gomock.Eq("1"),
				)
			},
			http.StatusNotFound,
		},
		{
			"unauthorized",
			&request.UpdateServiceRequest{
				BasicService: request.BasicService{
					Name:        "Service",
					Cost:        1000000,
					Phone:       "08123456789",
					Email:       "user@example.com",
					Description: "Lorem ipsum",
				},
			},
			createContext(jwt.NewWithClaims(
				jwt.SigningMethodHS256,
				&helper.JWTCustomClaims{ID: 1},
			)),
			func() {
				s.serviceRepository.EXPECT().Find(
					gomock.Eq(&model.Service{}),
					gomock.Eq("1"),
				).SetArg(0, model.Service{Model: gorm.Model{ID: 1}})
			},
			http.StatusUnauthorized,
		},
		{
			"ok",
			&request.UpdateServiceRequest{
				BasicService: request.BasicService{
					Name:        "Service",
					Cost:        1000000,
					Phone:       "08123456789",
					Email:       "user@example.com",
					Description: "Lorem ipsum",
				},
			},
			createContext(jwt.NewWithClaims(
				jwt.SigningMethodHS256,
				&helper.JWTCustomClaims{ID: 2},
			)),
			func() {
				s.serviceRepository.EXPECT().Find(
					gomock.Eq(&model.Service{}),
					gomock.Eq("1"),
				).SetArg(0, model.Service{Model: gorm.Model{ID: 1}, UserID: 2})

				s.serviceRepository.EXPECT().Update(gomock.Any(), gomock.Any())
			},
			http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		s.T().Run(testCase.Name, func(t *testing.T) {
			testCase.ExpectedFunc()
			if apiError := s.usecase.UpdateService(testCase.Context, &model.Service{}, testCase.Body); apiError != nil {
				code, _ := apiError.APIError()
				s.Equal(testCase.ExpectedCode, code)
			}
		})
	}
}

func (s *serviceUsecaseSuite) TestDeleteService() {
	createContext := func(token *jwt.Token) echo.Context {
		req := httptest.NewRequest("", "/", nil)
		rec := httptest.NewRecorder()
		ctx := echo.New().NewContext(req, rec)
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
			createContext(jwt.NewWithClaims(
				jwt.SigningMethodHS256,
				&helper.JWTCustomClaims{ID: 1},
			)),
			func() {
				s.serviceRepository.EXPECT().Find(
					gomock.Eq(&model.Service{}),
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
			)),
			func() {
				s.serviceRepository.EXPECT().Find(
					gomock.Eq(&model.Service{}),
					gomock.Eq("1"),
				).SetArg(0, model.Service{Model: gorm.Model{ID: 1}, UserID: 2})
			},
			http.StatusUnauthorized,
		},
		{
			"ok",
			nil,
			createContext(jwt.NewWithClaims(
				jwt.SigningMethodHS256,
				&helper.JWTCustomClaims{ID: 1},
			)),
			func() {
				s.serviceRepository.EXPECT().Find(
					gomock.Eq(&model.Service{}),
					gomock.Eq("1"),
				).SetArg(0, model.Service{Model: gorm.Model{ID: 1}, UserID: 1})

				s.serviceRepository.EXPECT().Delete(gomock.Any())
			},
			http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		s.T().Run(testCase.Name, func(t *testing.T) {
			testCase.ExpectedFunc()
			if apiError := s.usecase.DeleteService(testCase.Context, &model.Service{}); apiError != nil {
				code, _ := apiError.APIError()
				s.Equal(testCase.ExpectedCode, code)
			}
		})
	}
}
