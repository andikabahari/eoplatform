package usecase

import (
	"log"
	"net/http"

	"github.com/andikabahari/eoplatform/config"
	"github.com/andikabahari/eoplatform/helper"
	"github.com/andikabahari/eoplatform/model"
	r "github.com/andikabahari/eoplatform/repository"
	"github.com/andikabahari/eoplatform/request"
	"golang.org/x/crypto/bcrypt"
)

type AccountUsecase interface {
	GetAccount(claims *helper.JWTCustomClaims, user *model.User) helper.APIError
	UpdateAccount(claims *helper.JWTCustomClaims, user *model.User, req *request.UpdateUserRequest) helper.APIError
	ResetPassword(claims *helper.JWTCustomClaims, user *model.User, req *request.UpdateUserPasswordRequest) helper.APIError
}

type accountUsecase struct {
	userRepository r.UserRepository
}

func NewAccountUsecase(userRepository r.UserRepository) AccountUsecase {
	return &accountUsecase{userRepository}
}

func (u *accountUsecase) GetAccount(claims *helper.JWTCustomClaims, user *model.User) helper.APIError {
	u.userRepository.Find(user, claims.ID)

	if user.ID == 0 {
		return helper.NewAPIError(http.StatusNotFound, "user not found")
	}

	return nil
}

func (u *accountUsecase) UpdateAccount(claims *helper.JWTCustomClaims, user *model.User, req *request.UpdateUserRequest) helper.APIError {
	u.userRepository.Find(user, claims.ID)

	if user.ID == 0 {
		return helper.NewAPIError(http.StatusNotFound, "user not found")
	}

	u.userRepository.Update(user, req)

	return nil
}

func (u *accountUsecase) ResetPassword(claims *helper.JWTCustomClaims, user *model.User, req *request.UpdateUserPasswordRequest) helper.APIError {
	u.userRepository.Find(user, claims.ID)

	if user.ID == 0 {
		return helper.NewAPIError(http.StatusNotFound, "user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return helper.NewAPIError(http.StatusUnauthorized, "unauthorized")
	}

	authConfig := config.LoadAuthConfig()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), authConfig.Cost)
	if err != nil {
		log.Printf("Error: %s", err)
		return helper.NewAPIError(http.StatusInternalServerError, "internal server error")
	}

	u.userRepository.ResetPassword(user, string(hashedPassword))

	return nil
}
