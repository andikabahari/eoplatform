package response

import "github.com/andikabahari/eoplatform/model"

type FeedbackResponse struct {
	ID          uint          `json:"id"`
	Description string        `json:"description"`
	Rating      uint          `json:"rating"`
	Positive    float64       `json:"positive"`
	Negative    float64       `json:"negative"`
	ToUser      *UserResponse `json:"to_user"`
	FromUser    *UserResponse `json:"from_user"`
}

func NewFeedbackResponse(feedback model.Feedback) *FeedbackResponse {
	res := FeedbackResponse{}
	res.ID = feedback.ID
	res.Description = feedback.Description
	res.Rating = feedback.Rating
	res.Positive = feedback.Positive
	res.Negative = feedback.Negative
	res.ToUser = NewUserResponse(feedback.ToUser)
	res.FromUser = NewUserResponse(feedback.FromUser)

	return &res
}

func NewFeedbacksResponse(feedbacks []model.Feedback) *[]FeedbackResponse {
	res := make([]FeedbackResponse, 0)
	for _, feedback := range feedbacks {
		res = append(res, *NewFeedbackResponse(feedback))
	}

	return &res
}
