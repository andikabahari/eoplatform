package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/andikabahari/eoplatform/server"
	"github.com/andikabahari/eoplatform/testhelper"
	"github.com/stretchr/testify/suite"
)

type helloHandlerSuite struct {
	suite.Suite
	server  *server.Server
	handler *HelloHandler
}

func (s *helloHandlerSuite) SetupSuite() {
	os.Setenv("APP_ENV", "production")

	conn, _ := testhelper.Mock()
	s.server = testhelper.NewServer(conn)
	s.handler = NewHelloHandler(s.server)
}

func TestHelloHandlerSuite(t *testing.T) {
	suite.Run(t, new(helloHandlerSuite))
}

func (s *helloHandlerSuite) TestGreeting() {
	testCases := []struct {
		Name         string
		Endpoint     string
		Method       string
		Body         any
		ExpectedCode int
	}{
		{
			"ok",
			"/hello",
			http.MethodGet,
			nil,
			http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		s.T().Run(testCase.Name, func(t *testing.T) {
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

			s.NoError(s.handler.Greeting(ctx))
			s.Equal(testCase.ExpectedCode, rec.Code)
		})
	}
}
