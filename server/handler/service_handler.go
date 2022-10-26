package handler

import (
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

type ServiceHandler struct {
	server *s.Server
}

func NewServiceHandler(server *s.Server) *ServiceHandler {
	return &ServiceHandler{server}
}

func (h ServiceHandler) CreateService(c echo.Context) error {
	req := request.CreateServiceRequest{}

	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err,
		})
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*helper.JWTCustomClaims)

	service := model.Service{}
	service.UserID = claims.ID
	service.Name = req.Name
	service.Description = req.Description
	service.Cost = req.Cost

	serviceRepository := repository.NewServiceRepository(h.server.DB)
	serviceRepository.Create(&service)

	return c.JSON(http.StatusOK, echo.Map{
		"data": response.NewServiceResponse(service),
	})
}
