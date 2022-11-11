package usecase

import (
	"net/http"
	"os"
	"testing"

	"github.com/andikabahari/eoplatform/config"
	"github.com/andikabahari/eoplatform/helper"
	"github.com/andikabahari/eoplatform/model"
	mr "github.com/andikabahari/eoplatform/repository/mock_repository"
	"github.com/andikabahari/eoplatform/request"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type accountUsecaseSuite struct {
	suite.Suite

	ctrl           *gomock.Controller
	userRepository *mr.MockUserRepository

	usecase AccountUsecase
}

func (s *accountUsecaseSuite) SetupSuite() {
	os.Setenv("APP_ENV", "production")

	s.ctrl = gomock.NewController(s.T())
	s.userRepository = mr.NewMockUserRepository(s.ctrl)

	s.usecase = NewAccountUsecase(s.userRepository)
}

func (s *accountUsecaseSuite) TearDownSuite() {
	s.ctrl.Finish()
}

func TestAccountUsecaseSuite(t *testing.T) {
	suite.Run(t, new(accountUsecaseSuite))
}

func (s *accountUsecaseSuite) TestGetAccount() {
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
			&helper.JWTCustomClaims{},
			func() {
				s.userRepository.EXPECT().Find(
					gomock.Eq(&model.User{}),
					gomock.Eq(uint(0)),
				)
			},
			http.StatusNotFound,
		},
		{
			"ok",
			nil,
			&helper.JWTCustomClaims{ID: 1},
			func() {
				s.userRepository.EXPECT().Find(
					gomock.Eq(&model.User{}),
					gomock.Eq(uint(1)),
				).SetArg(0, model.User{Model: gorm.Model{ID: 1}})
			},
			http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		s.T().Run(testCase.Name, func(t *testing.T) {
			testCase.ExpectedFunc()
			if apiError := s.usecase.GetAccount(testCase.Claims, &model.User{}); apiError != nil {
				code, _ := apiError.APIError()
				s.Equal(testCase.ExpectedCode, code)
			}
		})
	}
}

func (s *accountUsecaseSuite) TestUpdateAccount() {
	testCases := []struct {
		Name         string
		Body         *request.UpdateUserRequest
		Claims       *helper.JWTCustomClaims
		ExpectedFunc func()
		ExpectedCode int
	}{
		{
			"not found",
			&request.UpdateUserRequest{
				Name: "User",
			},
			&helper.JWTCustomClaims{},
			func() {
				s.userRepository.EXPECT().Find(
					gomock.Eq(&model.User{}),
					gomock.Eq(uint(0)),
				)
			},
			http.StatusNotFound,
		},
		{
			"ok",
			&request.UpdateUserRequest{
				Name: "User",
			},
			&helper.JWTCustomClaims{ID: 1},
			func() {
				s.userRepository.EXPECT().Find(
					gomock.Eq(&model.User{}),
					gomock.Eq(uint(1)),
				).SetArg(0, model.User{Model: gorm.Model{ID: 1}})

				s.userRepository.EXPECT().Update(gomock.Any(), gomock.Any())
			},
			http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		s.T().Run(testCase.Name, func(t *testing.T) {
			testCase.ExpectedFunc()
			if apiError := s.usecase.UpdateAccount(testCase.Claims, &model.User{}, testCase.Body); apiError != nil {
				code, _ := apiError.APIError()
				s.Equal(testCase.ExpectedCode, code)
			}
		})
	}
}

func (s *accountUsecaseSuite) TestResetPassword() {
	testCases := []struct {
		Name         string
		Body         *request.UpdateUserPasswordRequest
		Claims       *helper.JWTCustomClaims
		ExpectedFunc func()
		ExpectedCode int
	}{
		{
			"not found",
			&request.UpdateUserPasswordRequest{
				Password:        "password",
				ConfirmPassword: "password",
				OldPassword:     "password",
			},
			&helper.JWTCustomClaims{},
			func() {
				s.userRepository.EXPECT().Find(
					gomock.Eq(&model.User{}),
					gomock.Eq(uint(0)),
				)
			},
			http.StatusNotFound,
		},
		{
			"unauthorized",
			&request.UpdateUserPasswordRequest{
				Password:        "password",
				ConfirmPassword: "password",
				OldPassword:     "wrong",
			},
			&helper.JWTCustomClaims{ID: 1},
			func() {
				s.userRepository.EXPECT().Find(
					gomock.Eq(&model.User{}),
					gomock.Eq(uint(1)),
				).SetArg(0, model.User{Model: gorm.Model{ID: 1}, Password: "password"})
			},
			http.StatusUnauthorized,
		},
		{
			"ok",
			&request.UpdateUserPasswordRequest{
				Password:        "password",
				ConfirmPassword: "password",
				OldPassword:     "password",
			},
			&helper.JWTCustomClaims{ID: 1},
			func() {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), config.LoadAuthConfig().Cost)
				s.userRepository.EXPECT().Find(
					gomock.Eq(&model.User{}),
					gomock.Eq(uint(1)),
				).SetArg(0, model.User{Model: gorm.Model{ID: 1}, Password: string(hashedPassword)})

				s.userRepository.EXPECT().ResetPassword(gomock.Any(), gomock.Any())
			},
			http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		s.T().Run(testCase.Name, func(t *testing.T) {
			testCase.ExpectedFunc()
			if apiError := s.usecase.ResetPassword(testCase.Claims, &model.User{}, testCase.Body); apiError != nil {
				code, _ := apiError.APIError()
				s.Equal(testCase.ExpectedCode, code)
			}
		})
	}
}
