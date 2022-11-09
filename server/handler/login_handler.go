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
	server         *s.Server
	userRepository repository.UserRepository
}

func NewLoginHandler(
	server *s.Server,
	userRepository repository.UserRepository,
) *LoginHandler {
	return &LoginHandler{
		server,
		userRepository,
	}
}

func (h *LoginHandler) Login(c echo.Context) error {
	req := request.LoginRequest{}

	if err := c.Bind(&req); err != nil {
		return err
	}

	user := model.User{}

	h.userRepository.FindByUsername(&user, req.Username)

	if user.ID == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "login failure",
			"error":   "user not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "login failure",
			"error":   "invalid password",
		})
	}

	token, err := helper.CreateToken(user.ID, user.Role)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "login successful",
		"data":    token,
	})
}
