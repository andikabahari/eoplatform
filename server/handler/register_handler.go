package handler

import (
	"net/http"

	"github.com/andikabahari/eoplatform/model"
	"github.com/andikabahari/eoplatform/repository"
	"github.com/andikabahari/eoplatform/request"
	"github.com/andikabahari/eoplatform/response"
	s "github.com/andikabahari/eoplatform/server"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type RegisterHandler struct {
	server *s.Server
}

func NewRegisterHandler(server *s.Server) *RegisterHandler {
	return &RegisterHandler{server}
}

func (h *RegisterHandler) Register(c echo.Context) error {
	req := request.CreateUserRequest{}

	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err,
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), h.server.Config.Auth.Cost)
	if err != nil {
		return err
	}

	user := model.User{}
	user.Name = req.Name
	user.Username = req.Username
	user.Password = string(hashedPassword)
	user.Email = req.Email
	user.Address = req.Address

	userRepository := repository.NewUserRepository(h.server.DB)
	userRepository.Create(&user)

	return c.JSON(http.StatusOK, echo.Map{
		"data": response.NewUserResponse(user),
	})
}
