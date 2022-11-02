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
			"message": "validation error",
			"error":   err,
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), h.server.Config.Auth.Cost)
	if err != nil {
		return err
	}

	userRepository := repository.NewUserRepository(h.server.DB)

	existingUser := model.User{}
	userRepository.FindByUsername(&existingUser, req.Username)
	if existingUser.ID > 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "registration failure",
			"error":   "user with the same username already exists",
		})
	}

	user := model.User{}
	user.Name = req.Name
	user.Username = req.Username
	user.Password = string(hashedPassword)
	user.Role = req.Role

	userRepository.Create(&user)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "registration successful",
		"data":    response.NewUserResponse(user),
	})
}
