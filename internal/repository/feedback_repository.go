package repository

import (
	"github.com/SuccorTrail/SuccorTrail/internal/db"
	"github.com/SuccorTrail/SuccorTrail/internal/model"
)

type FeedbackRepository interface {
	Create(feedback *model.Feedback) error
}

type feedbackRepository struct{}

func NewFeedbackRepository() FeedbackRepository {
	return &feedbackRepository{}
}

func (r *feedbackRepository) Create(feedback *model.Feedback) error {
	_, err := db.GetDB().Exec(
		"INSERT INTO feedback (id, donation_id, quality, comments, created_at) VALUES (?, ?, ?, ?, ?)",
		feedback.ID, feedback.DonationID, feedback.Quality, feedback.Comments, feedback.CreatedAt)
	return err
}
