package usecase

import (
	"net/http"

	"github.com/andikabahari/eoplatform/helper"
	"github.com/andikabahari/eoplatform/model"
	r "github.com/andikabahari/eoplatform/repository"
	"github.com/andikabahari/eoplatform/request"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type ServiceUsecase interface {
	GetServices(services *[]model.Service, keyword string)
	FindService(service *model.Service, id string) helper.APIError
	CreateService(claims *helper.JWTCustomClaims, service *model.Service, req *request.CreateServiceRequest)
	UpdateService(ctx echo.Context, service *model.Service, req *request.UpdateServiceRequest) helper.APIError
	DeleteService(ctx echo.Context, service *model.Service) helper.APIError
}

type serviceUsecase struct {
	serviceRepository r.ServiceRepository
}

func NewServiceUsecase(serviceRepository r.ServiceRepository) ServiceUsecase {
	return &serviceUsecase{serviceRepository}
}

func (u *serviceUsecase) GetServices(services *[]model.Service, keyword string) {
	u.serviceRepository.Get(services, keyword)
}

func (u *serviceUsecase) FindService(service *model.Service, id string) helper.APIError {
	u.serviceRepository.Find(service, id)

	if service.ID == 0 {
		return helper.NewAPIError(http.StatusNotFound, "service not found")
	}

	return nil
}

func (u *serviceUsecase) CreateService(claims *helper.JWTCustomClaims, service *model.Service, req *request.CreateServiceRequest) {
	service.UserID = claims.ID
	service.Name = req.Name
	service.Cost = req.Cost
	service.Phone = req.Phone
	service.Email = req.Email
	service.Description = req.Description

	u.serviceRepository.Create(service)
}

func (u *serviceUsecase) UpdateService(ctx echo.Context, service *model.Service, req *request.UpdateServiceRequest) helper.APIError {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(*helper.JWTCustomClaims)

	u.serviceRepository.Find(service, ctx.Param("id"))

	if service.ID == 0 {
		return helper.NewAPIError(http.StatusNotFound, "service not found")
	}

	if service.UserID != claims.ID {
		return helper.NewAPIError(http.StatusUnauthorized, "unauthorized")
	}

	u.serviceRepository.Update(service, req)

	return nil
}

func (u *serviceUsecase) DeleteService(ctx echo.Context, service *model.Service) helper.APIError {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(*helper.JWTCustomClaims)

	u.serviceRepository.Find(service, ctx.Param("id"))

	if service.ID == 0 {
		return helper.NewAPIError(http.StatusNotFound, "service not found")
	}

	if service.UserID != claims.ID {
		return helper.NewAPIError(http.StatusUnauthorized, "unauthorized")
	}

	u.serviceRepository.Delete(service)

	return nil
}
