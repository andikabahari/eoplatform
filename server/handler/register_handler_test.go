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

type registerHandlerSuite struct {
	suite.Suite

	ctrl    *gomock.Controller
	usecase *mu.MockRegisterUsecase

	server  *server.Server
	handler *RegisterHandler
}

func (s *registerHandlerSuite) SetupSuite() {
	os.Setenv("APP_ENV", "production")

	s.ctrl = gomock.NewController(s.T())
	s.usecase = mu.NewMockRegisterUsecase(s.ctrl)

	conn, _ := testhelper.Mock()
	s.server = testhelper.NewServer(conn)
	s.handler = NewRegisterHandler(s.usecase)
}

func (s *registerHandlerSuite) TearDownSuite() {
	s.ctrl.Finish()
}

func TestRegisterHandlerSuite(t *testing.T) {
	suite.Run(t, new(registerHandlerSuite))
}

func (s *registerHandlerSuite) TestRegister() {
	testCases := []struct {
		Name         string
		Endpoint     string
		Method       string
		Body         *request.CreateUserRequest
		ExpectedFunc func()
		ExpectedCode int
	}{
		{
			"bad request",
			"/v1/register",
			http.MethodPost,
			nil,
			func() {},
			http.StatusBadRequest,
		},
		{
			"bad request",
			"/v1/register",
			http.MethodPost,
			&request.CreateUserRequest{
				Name:     "Organizer",
				Username: "organizer",
				Password: "password",
				Role:     "organizer",
			},
			func() {
				apiError := helper.NewAPIError(http.StatusBadRequest, "")
				s.usecase.EXPECT().Register(gomock.Any(), gomock.Any()).Return(apiError)
			},
			http.StatusBadRequest,
		},
		{
			"ok",
			"/v1/register",
			http.MethodPost,
			&request.CreateUserRequest{
				Name:     "Organizer",
				Username: "organizer",
				Password: "password",
				Role:     "organizer",
			},
			func() {
				s.usecase.EXPECT().Register(gomock.Any(), gomock.Any()).Return(nil)
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

			s.NoError(s.handler.Register(ctx))
			s.Equal(testCase.ExpectedCode, rec.Code)
		})
	}
}
