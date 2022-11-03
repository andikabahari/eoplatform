package test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/andikabahari/eoplatform/server"
	"github.com/andikabahari/eoplatform/server/handler"
	"github.com/andikabahari/eoplatform/test/testhelper"
	"github.com/stretchr/testify/suite"
)

type helloSuite struct {
	suite.Suite
	mock    sqlmock.Sqlmock
	server  *server.Server
	handler *handler.HelloHandler
}

func (s *helloSuite) SetupSuite() {
	os.Setenv("APP_ENV", "production")

	var conn *sql.DB
	conn, s.mock = testhelper.Mock()
	s.server = testhelper.NewServer(conn)
	s.handler = handler.NewHelloHandler(s.server)
}

func TestHelloSuite(t *testing.T) {
	suite.Run(t, new(helloSuite))
}

func (s *helloSuite) TestGreeting() {
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
