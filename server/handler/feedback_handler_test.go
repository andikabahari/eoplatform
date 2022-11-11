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

type feedbackHandlerSuite struct {
	suite.Suite

	ctrl    *gomock.Controller
	usecase *mu.MockFeedbackUsecase

	server  *server.Server
	handler *FeedbackHandler
}

func (s *feedbackHandlerSuite) SetupSuite() {
	os.Setenv("APP_ENV", "production")

	s.ctrl = gomock.NewController(s.T())
	s.usecase = mu.NewMockFeedbackUsecase(s.ctrl)

	conn, _ := testhelper.Mock()
	s.server = testhelper.NewServer(conn)
	s.handler = NewFeedbackHandler(s.usecase)
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
				s.usecase.EXPECT().GetFeedbacks(gomock.Any(), gomock.Any())
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
			"ok",
			"/v1/feedbacks",
			http.MethodPost,
			&request.CreateFeedbackRequest{
				Description: "Good job!",
				Rating:      5,
				ToUserID:    1,
			},
			http.StatusOK,
			func() {
				s.usecase.EXPECT().CreateFeedback(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			jwt.NewWithClaims(jwt.SigningMethodHS256, &helper.JWTCustomClaims{ID: 1, Role: "customer"}),
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

			s.NoError(s.handler.CreateFeedback(ctx))
			s.Equal(testCase.ExpectedCode, rec.Code)
		})
	}
}
