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

type FeedbackHandler struct {
	server *s.Server
}

func NewFeedbackHandler(server *s.Server) *FeedbackHandler {
	return &FeedbackHandler{server}
}

func (h *FeedbackHandler) CreateFeedback(c echo.Context) error {
	req := request.CreateFeedbackRequest{}

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

	feedback := model.Feedback{}
	feedback.Description = req.Description
	feedback.Rating = req.Rating
	feedback.FromUserID = claims.ID
	feedback.ToUserID = req.ToUserID

	score, err := helper.AnalyzeSentiment(req.Description)
	if err != nil {
		return err
	}

	if score >= 0 {
		feedback.Positive = float64(score)
	} else {
		feedback.Negative = float64(score)
	}

	feedbackRepository := repository.NewFeedbackRepository(h.server.DB)
	feedbackRepository.Create(&feedback)

	userRepository := repository.NewUserRepository(h.server.DB)
	userRepository.Find(&feedback.FromUser, claims.ID)
	userRepository.Find(&feedback.ToUser, req.ToUserID)

	return c.JSON(http.StatusOK, echo.Map{
		"data": response.NewFeedbackResponse(feedback),
	})
}
