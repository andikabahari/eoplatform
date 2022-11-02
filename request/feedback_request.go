package request

import validation "github.com/go-ozzo/ozzo-validation"

type CreateFeedbackRequest struct {
	Description string `json:"description"`
	Rating      int    `json:"rating"`
	ToUserID    uint   `json:"to_user_id"`
}

func (r CreateFeedbackRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Description, validation.Required, validation.Length(1, 250)),
		validation.Field(&r.Rating, validation.Required, validation.Min(0), validation.Max(5)),
		validation.Field(&r.ToUserID, validation.Required),
	)
}
