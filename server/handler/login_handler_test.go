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
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type loginHandlerSuite struct {
	suite.Suite

	ctrl    *gomock.Controller
	usecase *mu.MockLoginUsecase

	server  *server.Server
	handler *LoginHandler
}

func (s *loginHandlerSuite) SetupSuite() {
	os.Setenv("APP_ENV", "production")

	s.ctrl = gomock.NewController(s.T())
	s.usecase = mu.NewMockLoginUsecase(s.ctrl)

	conn, _ := testhelper.Mock()
	s.server = testhelper.NewServer(conn)
	s.handler = NewLoginHandler(s.usecase)
}

func (s *loginHandlerSuite) TearDownSuite() {
	s.ctrl.Finish()
}

func TestLoginHandlerSuite(t *testing.T) {
	suite.Run(t, new(loginHandlerSuite))
}

func (s *loginHandlerSuite) TestLogin() {
	testCases := []struct {
		Name         string
		Endpoint     string
		Method       string
		Body         *request.LoginRequest
		ExpectedFunc func()
		ExpectedCode int
	}{
		{
			"not found",
			"/v1/login",
			http.MethodPost,
			nil,
			func() {
				apiError := helper.NewAPIError(http.StatusNotFound, "")
				s.usecase.EXPECT().Login(gomock.Any(), gomock.Any()).Return(apiError)
			},
			http.StatusNotFound,
		},
		{
			"ok",
			"/v1/login",
			http.MethodPost,
			&request.LoginRequest{
				Username: "organizer",
				Password: "password",
			},
			func() {
				s.usecase.EXPECT().Login(gomock.Any(), gomock.Any()).Return(nil)
			},
			http.StatusOK,
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

			s.NoError(s.handler.Login(ctx))
			s.Equal(testCase.ExpectedCode, rec.Code)
		})
	}
}
