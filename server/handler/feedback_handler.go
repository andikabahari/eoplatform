package handler

import (
	"math"
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

type FeedbackHandler struct {
	server             *s.Server
	feedbackRepository *repository.FeedbackRepository
	userRepository     *repository.UserRepository
}

func NewFeedbackHandler(server *s.Server) *FeedbackHandler {
	return &FeedbackHandler{
		server,
		repository.NewFeedbackRepository(server.DB),
		repository.NewUserRepository(server.DB),
	}
}

func (h *FeedbackHandler) GetFeedbacks(c echo.Context) error {
	feedbacks := make([]model.Feedback, 0)
	h.feedbackRepository.Get(&feedbacks, c.QueryParam("to_user_id"))

	return c.JSON(http.StatusOK, echo.Map{
		"message": "fetch feedbacks successful",
		"data":    response.NewFeedbacksResponse(feedbacks),
	})
}

func (h *FeedbackHandler) CreateFeedback(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*helper.JWTCustomClaims)

	if claims.Role != "customer" {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "create feedback failure",
			"error":   "unauthorized",
		})
	}

	req := request.CreateFeedbackRequest{}

	if err := c.Bind(&req); err != nil {
		return err
	}

	if err := req.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "validation error",
			"error":   err,
		})
	}

	feedbacksCount := h.feedbackRepository.GetFeedbacksCount(claims.ID, req.ToUserID)
	ordersCount := h.feedbackRepository.GetOrdersCount(claims.ID, req.ToUserID)

	if feedbacksCount >= ordersCount {
		return c.JSON(http.StatusForbidden, echo.Map{
			"message": "create feedback failure",
			"error":   "forbidden",
		})
	}

	feedback := model.Feedback{}
	feedback.Description = req.Description
	feedback.Rating = uint(req.Rating)
	feedback.FromUserID = claims.ID
	feedback.ToUserID = req.ToUserID

	score, err := helper.AnalyzeSentiment(req.Description)
	if err != nil {
		return err
	}

	if score >= 0 {
		feedback.Positive = float64(score)
	} else {
		feedback.Negative = math.Abs(float64(score))
	}

	h.feedbackRepository.Create(&feedback)

	h.userRepository.Find(&feedback.FromUser, claims.ID)
	h.userRepository.Find(&feedback.ToUser, req.ToUserID)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "create feedback successful",
		"data":    response.NewFeedbackResponse(feedback),
	})
}
