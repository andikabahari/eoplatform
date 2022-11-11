package handler

import (
	"net/http"

	"github.com/andikabahari/eoplatform/helper"
	"github.com/andikabahari/eoplatform/model"
	"github.com/andikabahari/eoplatform/request"
	"github.com/andikabahari/eoplatform/response"
	u "github.com/andikabahari/eoplatform/usecase"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type ServiceHandler struct {
	usecase u.ServiceUsecase
}

func NewServiceHandler(usecase u.ServiceUsecase) *ServiceHandler {
	return &ServiceHandler{usecase}
}

func (h *ServiceHandler) GetServices(c echo.Context) error {
	services := make([]model.Service, 0)
	h.usecase.GetServices(&services, c.QueryParam("q"))

	return c.JSON(http.StatusOK, echo.Map{
		"message": "fetch services successful",
		"data":    response.NewServicesResponse(services),
	})
}

func (h *ServiceHandler) FindService(c echo.Context) error {
	service := model.Service{}

	if apiError := h.usecase.FindService(&service, c.Param("id")); apiError != nil {
		code, message := apiError.APIError()
		return c.JSON(code, echo.Map{
			"message": "fetch service failure",
			"error":   message,
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
	h.usecase.CreateService(claims, &service, &req)

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

	if apiError := h.usecase.UpdateService(c, &service, &req); apiError != nil {
		code, message := apiError.APIError()
		return c.JSON(code, echo.Map{
			"message": "update service failure",
			"error":   message,
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "update service successful",
		"data":    response.NewServiceResponse(service),
	})
}

func (h *ServiceHandler) DeleteService(c echo.Context) error {
	service := model.Service{}

	if apiError := h.usecase.DeleteService(c, &service); apiError != nil {
		code, message := apiError.APIError()
		return c.JSON(code, echo.Map{
			"message": "delete service failure",
			"error":   message,
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "delete service successful",
		"data": echo.Map{
			"kind":    "service",
			"id":      c.Param("id"),
			"deleted": true,
		},
	})
}
