package repository

import (
	"github.com/andikabahari/eoplatform/model"
	"gorm.io/gorm"
)

type FeedbackRepository interface {
}

type feedbackRepository struct {
	db *gorm.DB
}

func NewFeedbackRepository(db *gorm.DB) *feedbackRepository {
	return &feedbackRepository{db}
}

func (r *feedbackRepository) Create(feedback *model.Feedback) {
	r.db.Debug().Omit("FromUser").Omit("ToUser").Save(feedback)
}
