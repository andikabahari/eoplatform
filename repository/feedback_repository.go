package repository

import (
	"database/sql"

	"github.com/andikabahari/eoplatform/model"
	"gorm.io/gorm"
)

type IFeedbackRepository interface {
	Get(feedbacks *model.Feedback, toUserID string)
	Create(feedback *model.Feedback)
	GetFeedbacksCount(fromUserID, toUserID any) int
	GetOrdersCount(fromUserID, toUserID any) int
}

type FeedbackRepository struct {
	db *gorm.DB
}

func NewFeedbackRepository(db *gorm.DB) *FeedbackRepository {
	return &FeedbackRepository{db}
}

func (r *FeedbackRepository) Get(feedbacks *[]model.Feedback, toUserID string) {
	if toUserID != "" {
		r.db.Debug().Preload("FromUser").Preload("ToUser").Where("to_user_id = ?", toUserID).Find(feedbacks)
	} else {
		r.db.Debug().Preload("FromUser").Preload("ToUser").Find(feedbacks)
	}
}

func (r *FeedbackRepository) Create(feedback *model.Feedback) {
	r.db.Debug().Omit("FromUser").Omit("ToUser").Save(feedback)
}

func (r *FeedbackRepository) GetFeedbacksCount(fromUserID, toUserID any) int {
	feedbacksCount := 0

	query := "SELECT COUNT(1) FROM feedbacks " +
		"WHERE from_user_id=@FromUserID AND to_user_id=@ToUserID"

	r.db.Debug().Raw(query,
		sql.Named("FromUserID", fromUserID),
		sql.Named("ToUserID", toUserID),
	).Scan(&feedbacksCount)

	return feedbacksCount
}

func (r *FeedbackRepository) GetOrdersCount(fromUserID, toUserID any) int {
	ordersCount := 0

	query := "SELECT COUNT(1) FROM(" +
		"SELECT DISTINCT o.id FROM orders o " +
		"JOIN order_services os ON os.order_id=o.id " +
		"JOIN services s ON s.id=os.service_id " +
		"WHERE o.user_id=@FromUserID AND s.user_id=@ToUserID AND is_completed>0" +
		") AS t"

	r.db.Debug().Raw(query,
		sql.Named("FromUserID", fromUserID),
		sql.Named("ToUserID", toUserID),
	).Scan(&ordersCount)

	return ordersCount
}
