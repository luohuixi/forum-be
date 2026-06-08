package dao

import "time"

type FeedbackModel struct {
	ID        uint32 `gorm:"primaryKey"`
	UserID    uint32
	Category  string
	Content   string
	Contact   string
	ImgURL    string `gorm:"column:img_url"`
	CreatedAt time.Time
}

func (FeedbackModel) TableName() string {
	return "feedbacks"
}

func (d *Dao) CreateFeedback(feedback *FeedbackModel) error {
	return d.DB.Create(feedback).Error
}
