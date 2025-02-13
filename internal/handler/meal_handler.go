package handler

import (
	"encoding/json"
	"net/http"

	"github.com/SuccorTrail/SuccorTrail/internal/repository"
	"github.com/sirupsen/logrus"
)

type MealHandler struct {
	Repo repository.MealRepository
}

func (h *MealHandler) GetAvailableMeals(w http.ResponseWriter, r *http.Request) {
	location := r.URL.Query().Get("location")
	if location == "" {
		http.Error(w, "Location is required", http.StatusBadRequest)
		return
	}

	meals, err := h.Repo.GetAvailableMeals(location)
	if err != nil {
		logrus.WithError(err).Error("Error fetching meals")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(meals)
}

func (h *MealHandler) VerifyMeal(w http.ResponseWriter, r *http.Request) {
	var data struct {
		DonationID string `json:"donationId"`
		ReceiverID string `json:"receiverId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		logrus.WithError(err).Error("Error decoding request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if data.DonationID == "" || data.ReceiverID == "" {
		http.Error(w, "DonationID and ReceiverID are required", http.StatusBadRequest)
		return
	}

	// Verify the meal (implementation depends on business logic)
	// For now, just return a success message
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Meal verified successfully"})
}
