package handler

import (
	"fmt"
	"net/http"

	"github.com/andikabahari/eoplatform/helper"
	"github.com/andikabahari/eoplatform/model"
	"github.com/andikabahari/eoplatform/repository"
	"github.com/andikabahari/eoplatform/request"
	"github.com/andikabahari/eoplatform/response"
	s "github.com/andikabahari/eoplatform/server"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type AccountHandler struct {
	server *s.Server
}

func NewAccountHandler(server *s.Server) *AccountHandler {
	return &AccountHandler{server}
}

func (h *AccountHandler) GetAccount(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*helper.JWTCustomClaims)

	user := model.User{}
	userRepository := repository.NewUserRepository(h.server.DB)
	userRepository.Find(&user, fmt.Sprintf("%d", claims.ID))

	if user.ID == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "user not found",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": response.NewUserResponse(user),
	})
}

func (h *AccountHandler) UpdateAccount(c echo.Context) error {
	req := request.UpdateUserRequest{}

	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err,
		})
	}

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*helper.JWTCustomClaims)

	user := model.User{}
	userRepository := repository.NewUserRepository(h.server.DB)
	userRepository.Find(&user, fmt.Sprintf("%d", claims.ID))

	if user.ID == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{
			"error": "user not found",
		})
	}

	userRepository.Update(&user, &req)

	return c.JSON(http.StatusOK, echo.Map{
		"data": response.NewUserResponse(user),
	})
}
