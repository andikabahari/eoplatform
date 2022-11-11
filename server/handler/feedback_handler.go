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

type FeedbackHandler struct {
	usecase u.FeedbackUsecase
}

func NewFeedbackHandler(usecase u.FeedbackUsecase) *FeedbackHandler {
	return &FeedbackHandler{usecase}
}

func (h *FeedbackHandler) GetFeedbacks(c echo.Context) error {
	feedbacks := make([]model.Feedback, 0)
	h.usecase.GetFeedbacks(&feedbacks, c.QueryParam("to_user_id"))

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

	feedback := model.Feedback{}
	h.usecase.CreateFeedback(claims, &feedback, &req)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "create feedback successful",
		"data":    response.NewFeedbackResponse(feedback),
	})
}
