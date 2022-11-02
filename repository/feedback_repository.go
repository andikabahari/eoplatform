package repository

import (
	"github.com/andikabahari/eoplatform/model"
	"gorm.io/gorm"
)

type FeedbackRepository interface {
	Get(feedbacks *model.Feedback, toUserID string)
	Create(feedback *model.Feedback)
}

type feedbackRepository struct {
	db *gorm.DB
}

func NewFeedbackRepository(db *gorm.DB) *feedbackRepository {
	return &feedbackRepository{db}
}

func (r *feedbackRepository) Get(feedbacks *[]model.Feedback, toUserID string) {
	if toUserID != "" {
		r.db.Debug().Preload("FromUser").Preload("ToUser").Where("to_user_id = ?", toUserID).Find(feedbacks)
	} else {
		r.db.Debug().Preload("FromUser").Preload("ToUser").Find(feedbacks)
	}
}

func (r *feedbackRepository) Create(feedback *model.Feedback) {
	r.db.Debug().Omit("FromUser").Omit("ToUser").Save(feedback)
}
