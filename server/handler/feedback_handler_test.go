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
)

type feedbackHandlerSuite struct {
	suite.Suite

	ctrl               *gomock.Controller
	feedbackRepository *mock_repository.MockFeedbackRepository
	userRepository     *mock_repository.MockUserRepository

	server  *server.Server
	handler *FeedbackHandler
}

func (s *feedbackHandlerSuite) SetupSuite() {
	os.Setenv("APP_ENV", "production")

	s.ctrl = gomock.NewController(s.T())
	s.feedbackRepository = mock_repository.NewMockFeedbackRepository(s.ctrl)
	s.userRepository = mock_repository.NewMockUserRepository(s.ctrl)

	conn, _ := testhelper.Mock()
	s.server = testhelper.NewServer(conn)
	s.handler = NewFeedbackHandler(s.server, s.feedbackRepository, s.userRepository)
}

func TestFeedbackHandlerSuite(t *testing.T) {
	suite.Run(t, new(feedbackHandlerSuite))
}

func (s *feedbackHandlerSuite) TestGetFeedbacks() {
	testCases := []struct {
		Name         string
		Endpoint     string
		Method       string
		Body         any
		ExpectedCode int
		ExpectedFunc func()
		Token        *jwt.Token
	}{
		{
			"ok",
			"/v1/feedbacks",
			http.MethodGet,
			nil,
			http.StatusOK,
			func() {
				s.feedbackRepository.EXPECT().Get(
					gomock.Eq(&[]model.Feedback{}),
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

			s.NoError(s.handler.GetFeedbacks(ctx))
			s.Equal(testCase.ExpectedCode, rec.Code)
		})
	}
}

func (s *feedbackHandlerSuite) TestCreateFeedback() {
	testCases := []struct {
		Name         string
		Endpoint     string
		Method       string
		Body         *request.CreateFeedbackRequest
		ExpectedCode int
		ExpectedFunc func()
		Token        *jwt.Token
	}{
		{
			"unauthorized",
			"/v1/feedbacks",
			http.MethodPost,
			nil,
			http.StatusUnauthorized,
			func() {},
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{ID: 1, Role: "organizer"}),
		},
		{
			"bad request",
			"/v1/feedbacks",
			http.MethodPost,
			nil,
			http.StatusBadRequest,
			func() {},
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{ID: 1, Role: "customer"}),
		},
		{
			"forbidden",
			"/v1/feedbacks",
			http.MethodPost,
			&request.CreateFeedbackRequest{
				Description: "Good job!",
				Rating:      5,
				ToUserID:    1,
			},
			http.StatusForbidden,
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
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{ID: 1, Role: "customer"}),
		},
		// {
		// 	"ok",
		// 	"/v1/feedbacks",
		// 	http.MethodPost,
		// 	&request.CreateFeedbackRequest{
		// 		Description: "Good job!",
		// 		Rating:      5,
		// 		ToUserID:    1,
		// 	},
		// 	http.StatusOK,
		// 	func() {},
		// 	jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{ID: 1, Role: "customer"}),
		// },
		// {
		// 	"ok",
		// 	"/v1/feedbacks",
		// 	http.MethodPost,
		// 	&request.CreateFeedbackRequest{
		// 		Description: "This is bad!",
		// 		Rating:      1,
		// 		ToUserID:    1,
		// 	},
		// 	http.StatusOK,
		// 	func() {},
		// 	jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{ID: 1, Role: "customer"}),
		// },
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

			s.NoError(s.handler.CreateFeedback(ctx))
			s.Equal(testCase.ExpectedCode, rec.Code)
		})
	}
}
