package usecase

import (
	"log"
	"math"
	"net/http"

	"github.com/andikabahari/eoplatform/helper"
	"github.com/andikabahari/eoplatform/model"
	r "github.com/andikabahari/eoplatform/repository"
	"github.com/andikabahari/eoplatform/request"
)

type FeedbackUsecase interface {
	GetFeedbacks(feedbacks *[]model.Feedback, toUserID string)
	CreateFeedback(claims *helper.JWTCustomClaims, feedback *model.Feedback, req *request.CreateFeedbackRequest) helper.APIError
}

type feedbackUsecase struct {
	feedbackRepository r.FeedbackRepository
	userRepository     r.UserRepository
}

func NewFeedbackUsecase(feedbackRepository r.FeedbackRepository, userRepository r.UserRepository) FeedbackUsecase {
	return &feedbackUsecase{feedbackRepository, userRepository}
}

func (u *feedbackUsecase) GetFeedbacks(feedbacks *[]model.Feedback, toUserID string) {
	u.feedbackRepository.Get(feedbacks, toUserID)
}

func (u *feedbackUsecase) CreateFeedback(claims *helper.JWTCustomClaims, feedback *model.Feedback, req *request.CreateFeedbackRequest) helper.APIError {
	feedbacksCount := u.feedbackRepository.GetFeedbacksCount(claims.ID, req.ToUserID)
	ordersCount := u.feedbackRepository.GetOrdersCount(claims.ID, req.ToUserID)

	if feedbacksCount >= ordersCount {
		return helper.NewAPIError(http.StatusForbidden, "forbidden")
	}

	feedback.Description = req.Description
	feedback.Rating = uint(req.Rating)
	feedback.FromUserID = claims.ID
	feedback.ToUserID = req.ToUserID

	score, err := helper.AnalyzeSentiment(req.Description)
	if err != nil {
		log.Printf("Error: %s", err)
		return helper.NewAPIError(http.StatusInternalServerError, "internal server error")
	}

	if score >= 0 {
		feedback.Positive = float64(score)
	} else {
		feedback.Negative = math.Abs(float64(score))
	}

	u.feedbackRepository.Create(feedback)

	u.userRepository.Find(&feedback.FromUser, claims.ID)
	u.userRepository.Find(&feedback.ToUser, req.ToUserID)

	return nil
}
