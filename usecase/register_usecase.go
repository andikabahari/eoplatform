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

type RegisterUsecase interface {
	Register(user *model.User, req *request.CreateUserRequest) helper.APIError
}

type registerUsecase struct {
	userRepository r.UserRepository
}

func NewRegisterUsecase(userRepository r.UserRepository) RegisterUsecase {
	return &registerUsecase{userRepository}
}

func (u *registerUsecase) Register(user *model.User, req *request.CreateUserRequest) helper.APIError {
	authConfig := config.LoadAuthConfig()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), authConfig.Cost)
	if err != nil {
		log.Printf("Error: %s", err)
		return helper.NewAPIError(http.StatusInternalServerError, "internal server error")
	}

	existingUser := model.User{}
	u.userRepository.FindByUsername(&existingUser, req.Username)
	if existingUser.ID > 0 {
		return helper.NewAPIError(http.StatusBadRequest, "user with the same username already exists")
	}

	*user = model.User{}
	user.Name = req.Name
	user.Username = req.Username
	user.Password = string(hashedPassword)
	user.Role = req.Role

	u.userRepository.Create(user)

	return nil
}
