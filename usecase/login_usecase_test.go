package usecase

import (
	"net/http"
	"os"
	"testing"

	"github.com/andikabahari/eoplatform/config"
	"github.com/andikabahari/eoplatform/model"
	mr "github.com/andikabahari/eoplatform/repository/mock_repository"
	"github.com/andikabahari/eoplatform/request"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type loginUsecaseSuite struct {
	suite.Suite

	ctrl           *gomock.Controller
	userRepository *mr.MockUserRepository

	usecase LoginUsecase
}

func (s *loginUsecaseSuite) SetupSuite() {
	os.Setenv("APP_ENV", "production")

	s.ctrl = gomock.NewController(s.T())
	s.userRepository = mr.NewMockUserRepository(s.ctrl)

	s.usecase = NewLoginUsecase(s.userRepository)
}

func (s *loginUsecaseSuite) TearDownSuite() {
	s.ctrl.Finish()
}

func TestLoginUsecaseSuite(t *testing.T) {
	suite.Run(t, new(loginUsecaseSuite))
}

func (s *loginUsecaseSuite) TestLogin() {
	testCases := []struct {
		Name         string
		Body         *request.LoginRequest
		ExpectedFunc func()
		ExpectedCode int
	}{
		{
			"not found",
			&request.LoginRequest{},
			func() {
				s.userRepository.EXPECT().FindByUsername(
					gomock.Eq(&model.User{}),
					gomock.Eq(""),
				)
			},
			http.StatusNotFound,
		},
		{
			"bad request",
			&request.LoginRequest{
				Username: "organizer",
				Password: "password",
			},
			func() {
				s.userRepository.EXPECT().FindByUsername(
					gomock.Eq(&model.User{}),
					gomock.Eq("organizer"),
				).SetArg(0, model.User{Model: gorm.Model{ID: 1}})
			},
			http.StatusBadRequest,
		},
		{
			"ok",
			&request.LoginRequest{
				Username: "organizer",
				Password: "password",
			},
			func() {
				hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), config.LoadAuthConfig().Cost)
				s.userRepository.EXPECT().FindByUsername(
					gomock.Eq(&model.User{}),
					gomock.Eq("organizer"),
				).SetArg(0, model.User{
					Model:    gorm.Model{ID: 1},
					Password: string(hashedPassword),
				})
			},
			0,
		},
	}

	for _, testCase := range testCases {
		s.T().Run(testCase.Name, func(t *testing.T) {
			testCase.ExpectedFunc()
			if apiError := s.usecase.Login(new(string), testCase.Body); apiError != nil {
				code, _ := apiError.APIError()
				s.Equal(testCase.ExpectedCode, code)
			}
		})
	}
}
