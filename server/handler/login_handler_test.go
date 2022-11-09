package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/andikabahari/eoplatform/config"
	"github.com/andikabahari/eoplatform/model"
	"github.com/andikabahari/eoplatform/repository/mock_repository"
	"github.com/andikabahari/eoplatform/request"
	"github.com/andikabahari/eoplatform/server"
	"github.com/andikabahari/eoplatform/testhelper"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type loginHandlerSuite struct {
	suite.Suite

	ctrl           *gomock.Controller
	userRepository *mock_repository.MockUserRepository

	server  *server.Server
	handler *LoginHandler
}

func (s *loginHandlerSuite) SetupSuite() {
	os.Setenv("APP_ENV", "production")

	s.ctrl = gomock.NewController(s.T())
	s.userRepository = mock_repository.NewMockUserRepository(s.ctrl)

	conn, _ := testhelper.Mock()
	s.server = testhelper.NewServer(conn)
	s.handler = NewLoginHandler(s.server, s.userRepository)
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
				s.userRepository.EXPECT().FindByUsername(
					gomock.Eq(&model.User{}),
					gomock.Eq(""),
				)
			},
			http.StatusNotFound,
		},
		{
			"bad request",
			"/v1/login",
			http.MethodPost,
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
			"/v1/login",
			http.MethodPost,
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
