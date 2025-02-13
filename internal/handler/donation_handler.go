package handler

import (
	"encoding/json"
	"net/http"

	"github.com/SuccorTrail/SuccorTrail/internal/model"
	"github.com/SuccorTrail/SuccorTrail/internal/repository"
	"github.com/SuccorTrail/SuccorTrail/internal/util"
	"github.com/sirupsen/logrus"
)

type DonationHandler struct {
	Repo repository.DonationRepository
}

func (h *DonationHandler) CreateDonation(w http.ResponseWriter, r *http.Request) {
	var donation model.Donation
	if err := json.NewDecoder(r.Body).Decode(&donation); err != nil {
		logrus.WithError(err).Error("Error decoding donation")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if donation.Name == "" || donation.Email == "" || donation.Phone == "" ||
		donation.MealType == "" || donation.Quantity == 0 || donation.Location == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	if err := h.Repo.Create(&donation); err != nil {
		logrus.WithError(err).Error("Error creating donation")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(donation)
}

func (h *DonationHandler) RenderDonationForm(w http.ResponseWriter, r *http.Request) {
	util.RenderTemplate(w, "donor.html", nil)
}
