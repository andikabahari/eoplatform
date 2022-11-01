package handler

import (
	"net/http"

	"github.com/andikabahari/eoplatform/helper"
	"github.com/andikabahari/eoplatform/model"
	"github.com/andikabahari/eoplatform/repository"
	"github.com/andikabahari/eoplatform/request"
	s "github.com/andikabahari/eoplatform/server"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type LoginHandler struct {
	server *s.Server
}

func NewLoginHandler(server *s.Server) *LoginHandler {
	return &LoginHandler{server}
}

func (h *LoginHandler) Login(c echo.Context) error {
	req := request.LoginRequest{}

	if err := c.Bind(&req); err != nil {
		return err
	}

	user := model.User{}

	userRepository := repository.NewUserRepository(h.server.DB)
	userRepository.FindByUsername(&user, req.Username)

	if user.ID == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "user not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "invalid password",
		})
	}

	token, err := helper.CreateToken(user.ID, user.Role)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": token,
	})
}
