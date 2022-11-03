package response

import "github.com/andikabahari/eoplatform/model"

type FeedbackResponse struct {
	ID          uint          `json:"id"`
	Description string        `json:"description"`
	Sentiment   string        `json:"sentiment"`
	Rating      uint          `json:"rating"`
	ToUser      *UserResponse `json:"to_user"`
	FromUser    *UserResponse `json:"from_user"`
}

func NewFeedbackResponse(feedback model.Feedback) *FeedbackResponse {
	res := FeedbackResponse{}
	res.ID = feedback.ID
	res.Description = feedback.Description
	res.Rating = feedback.Rating
	if feedback.Positive > 0 {
		res.Sentiment = "positive"
	} else if feedback.Negative > 0 {
		res.Sentiment = "negative"
	}
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
