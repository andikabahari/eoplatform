package usecase

import (
	"log"
	"net/http"

	"github.com/andikabahari/eoplatform/helper"
	"github.com/andikabahari/eoplatform/model"
	r "github.com/andikabahari/eoplatform/repository"
	"github.com/andikabahari/eoplatform/request"

	"golang.org/x/crypto/bcrypt"
)

type LoginUsecase interface {
	Login(token *string, req *request.LoginRequest) helper.APIError
}

type loginUsecase struct {
	userRepository r.UserRepository
}

func NewLoginUsecase(userRepository r.UserRepository) LoginUsecase {
	return &loginUsecase{userRepository}
}

func (u *loginUsecase) Login(token *string, req *request.LoginRequest) helper.APIError {
	user := model.User{}

	u.userRepository.FindByUsername(&user, req.Username)

	if user.ID == 0 {
		return helper.NewAPIError(http.StatusNotFound, "user not found")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return helper.NewAPIError(http.StatusBadRequest, "invalid password")
	}

	*token, err = helper.CreateToken(user.ID, user.Role)
	if err != nil {
		log.Printf("Error: %s", err)
		return helper.NewAPIError(http.StatusInternalServerError, "internal server error")
	}

	return nil
}
