package handler

import (
	"encoding/json"
	"net/http"

	"github.com/SuccorTrail/SuccorTrail/internal/model"
	"github.com/SuccorTrail/SuccorTrail/internal/repository"
	"github.com/sirupsen/logrus"
)

type FeedbackHandler struct {
	Repo repository.FeedbackRepository
}

func (h *FeedbackHandler) SubmitFeedback(w http.ResponseWriter, r *http.Request) {
	var feedback model.Feedback
	if err := json.NewDecoder(r.Body).Decode(&feedback); err != nil {
		logrus.WithError(err).Error("Error decoding feedback")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if feedback.DonationID == "" || feedback.Quality == 0 || feedback.Comments == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	if err := h.Repo.Create(&feedback); err != nil {
		logrus.WithError(err).Error("Error creating feedback")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(feedback)
}
