package usecase

import (
	"os"
	"testing"

	"github.com/andikabahari/eoplatform/model"
	mr "github.com/andikabahari/eoplatform/repository/mock_repository"
	"github.com/andikabahari/eoplatform/request"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type registerUsecaseSuite struct {
	suite.Suite

	ctrl           *gomock.Controller
	userRepository *mr.MockUserRepository

	usecase RegisterUsecase
}

func (s *registerUsecaseSuite) SetupSuite() {
	os.Setenv("APP_ENV", "production")

	s.ctrl = gomock.NewController(s.T())
	s.userRepository = mr.NewMockUserRepository(s.ctrl)

	s.usecase = NewRegisterUsecase(s.userRepository)
}

func (s *registerUsecaseSuite) TearDownSuite() {
	s.ctrl.Finish()
}

func TestRegisterUsecaseSuite(t *testing.T) {
	suite.Run(t, new(registerUsecaseSuite))
}

func (s *registerUsecaseSuite) TestRegister() {
	testCases := []struct {
		Name          string
		Body          *request.CreateUserRequest
		ExpectedFunc  func()
		ExpectedError bool
	}{
		{
			"bad request",
			&request.CreateUserRequest{
				Name:     "Organizer",
				Username: "organizer",
				Password: "password",
				Role:     "organizer",
			},
			func() {
				s.userRepository.EXPECT().FindByUsername(
					gomock.Eq(&model.User{}),
					gomock.Eq("organizer"),
				).SetArg(0, model.User{Model: gorm.Model{ID: 1}})
			},
			true,
		},
		{
			"ok",
			&request.CreateUserRequest{
				Name:     "Organizer",
				Username: "organizer",
				Password: "password",
				Role:     "organizer",
			},
			func() {
				s.userRepository.EXPECT().FindByUsername(
					gomock.Eq(&model.User{}),
					gomock.Eq("organizer"),
				)

				s.userRepository.EXPECT().Create(gomock.Any())
			},
			false,
		},
	}

	for _, testCase := range testCases {
		s.T().Run(testCase.Name, func(t *testing.T) {
			testCase.ExpectedFunc()
			apiError := s.usecase.Register(&model.User{}, testCase.Body)
			s.Equal(testCase.ExpectedError, apiError != nil)
		})
	}
}
