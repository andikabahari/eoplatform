package usecase

import (
	"net/http"
	"os"
	"testing"

	"github.com/andikabahari/eoplatform/helper"
	"github.com/andikabahari/eoplatform/model"
	mr "github.com/andikabahari/eoplatform/repository/mock_repository"
	"github.com/andikabahari/eoplatform/request"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type feedbackUsecaseSuite struct {
	suite.Suite

	ctrl               *gomock.Controller
	feedbackRepository *mr.MockFeedbackRepository
	userRepository     *mr.MockUserRepository

	usecase FeedbackUsecase
}

func (s *feedbackUsecaseSuite) SetupSuite() {
	os.Setenv("APP_ENV", "production")

	s.ctrl = gomock.NewController(s.T())
	s.feedbackRepository = mr.NewMockFeedbackRepository(s.ctrl)
	s.userRepository = mr.NewMockUserRepository(s.ctrl)

	s.usecase = NewFeedbackUsecase(s.feedbackRepository, s.userRepository)
}

func (s *feedbackUsecaseSuite) TearDownSuite() {
	s.ctrl.Finish()
}

func TestFeedbackUsecaseSuite(t *testing.T) {
	suite.Run(t, new(feedbackUsecaseSuite))
}

func (s *feedbackUsecaseSuite) TestGetFeedbacks() {
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
				s.feedbackRepository.EXPECT().Get(
					gomock.Eq(&[]model.Feedback{}),
					gomock.Eq(""),
				)
			},
			http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		s.T().Run(testCase.Name, func(t *testing.T) {
			testCase.ExpectedFunc()
			s.usecase.GetFeedbacks(&[]model.Feedback{}, "")
		})
	}
}

func (s *feedbackUsecaseSuite) TestCreateFeedback() {
	testCases := []struct {
		Name         string
		Body         *request.CreateFeedbackRequest
		Claims       *helper.JWTCustomClaims
		ExpectedFunc func()
		ExpectedCode int
	}{
		{
			"forbidden",
			&request.CreateFeedbackRequest{
				Description: "Good job!",
				Rating:      5,
				ToUserID:    1,
			},
			&helper.JWTCustomClaims{ID: 1, Role: "customer"},
			func() {
				s.feedbackRepository.EXPECT().GetFeedbacksCount(
					gomock.Eq(uint(1)),
					gomock.Eq(uint(1)),
				)

				s.feedbackRepository.EXPECT().GetOrdersCount(
					gomock.Eq(uint(1)),
					gomock.Eq(uint(1)),
				)
			},
			http.StatusForbidden,
		},
	}

	for _, testCase := range testCases {
		s.T().Run(testCase.Name, func(t *testing.T) {
			testCase.ExpectedFunc()
			if apiError := s.usecase.CreateFeedback(testCase.Claims, &model.Feedback{}, testCase.Body); apiError != nil {
				code, _ := apiError.APIError()
				s.Equal(testCase.ExpectedCode, code)
			}
		})
	}
}
