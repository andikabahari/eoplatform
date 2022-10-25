package handler

import (
	"net/http"

	s "github.com/andikabahari/eoplatform/server"
	"github.com/labstack/echo/v4"
)

type HelloHandler struct {
	server *s.Server
}

func NewHelloHandler(server *s.Server) *HelloHandler {
	return &HelloHandler{server}
}

func (h *HelloHandler) Greeting(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"message": "Welcome!"})
}
