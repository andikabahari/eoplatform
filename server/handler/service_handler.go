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

func (h *ServiceHandler) GetServices(c echo.Context) error {
	services := make([]model.Service, 0)
	serviceRepository := repository.NewServiceRepository(h.server.DB)
	serviceRepository.Get(&services, c.QueryParam("q"))

	return c.JSON(http.StatusOK, echo.Map{
		"message": "fetch services successful",
		"data":    response.NewServicesResponse(services),
	})
}

func (h *ServiceHandler) FindService(c echo.Context) error {
	service := model.Service{}
	serviceRepository := repository.NewServiceRepository(h.server.DB)
	serviceRepository.Find(&service, c.Param("id"))

	if service.ID == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "fetch service failure",
			"error":   "service not found",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "fetch service successful",
		"data":    response.NewServiceResponse(service),
	})
}

func (h *ServiceHandler) CreateService(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*helper.JWTCustomClaims)

	if claims.Role != "organizer" {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "create service failure",
			"error":   "unauthorized",
		})
	}

	req := request.CreateServiceRequest{}

	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "validation error",
			"error":   err,
		})
	}

	service := model.Service{}
	service.UserID = claims.ID
	service.Name = req.Name
	service.Cost = req.Cost
	service.Phone = req.Phone
	service.Email = req.Email
	service.IsPublished = req.IsPublished
	service.Description = req.Description

	serviceRepository := repository.NewServiceRepository(h.server.DB)
	serviceRepository.Create(&service)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "create service successful",
		"data":    response.NewServiceResponse(service),
	})
}

func (h *ServiceHandler) UpdateService(c echo.Context) error {
	req := request.UpdateServiceRequest{}

	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "validation error",
			"error":   err,
		})
	}

	service := model.Service{}
	serviceRepository := repository.NewServiceRepository(h.server.DB)
	serviceRepository.Find(&service, c.Param("id"))

	if service.ID == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "update service failure",
			"error":   "service not found",
		})
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*helper.JWTCustomClaims)

	if service.UserID != claims.ID {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "update service failure",
			"error":   "unauthorized",
		})
	}

	count := 0
	query := "SELECT COUNT(1) FROM order_services WHERE service_id = ?"
	h.server.DB.Raw(query, service.ID).Scan(&count)
	if count > 0 {
		return c.JSON(http.StatusForbidden, echo.Map{
			"message": "update service failure",
			"error":   "forbidden",
		})
	}

	serviceRepository.Update(&service, &req)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "update service successful",
		"data":    response.NewServiceResponse(service),
	})
}

func (h *ServiceHandler) DeleteService(c echo.Context) error {
	service := model.Service{}
	serviceRepository := repository.NewServiceRepository(h.server.DB)
	serviceRepository.Find(&service, c.Param("id"))

	if service.ID == 0 {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "delete service failure",
			"error":   "service not found",
		})
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*helper.JWTCustomClaims)

	if service.UserID != claims.ID {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "delete service failure",
			"error":   "unauthorized",
		})
	}

	serviceRepository.Delete(&service)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "delete service successful",
		"data": echo.Map{
			"kind":    "service",
			"id":      c.Param("id"),
			"deleted": true,
		},
	})
}
