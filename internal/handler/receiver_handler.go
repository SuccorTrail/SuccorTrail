package handler

import (
	"encoding/json"
	"net/http"

	"github.com/SuccorTrail/SuccorTrail/internal/model"
	"github.com/SuccorTrail/SuccorTrail/internal/repository"
	"github.com/SuccorTrail/SuccorTrail/internal/util"
	"github.com/sirupsen/logrus"
)

type ReceiverHandler struct {
	Repo repository.ReceiverRepository
}

func (h *ReceiverHandler) CreateReceiver(w http.ResponseWriter, r *http.Request) {
	var receiver model.Receiver
	if err := json.NewDecoder(r.Body).Decode(&receiver); err != nil {
		logrus.WithError(err).Error("Error decoding receiver")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if receiver.Name == "" || receiver.Phone == "" || receiver.Location == "" || receiver.FamilySize <= 0 {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	if err := h.Repo.Create(&receiver); err != nil {
		logrus.WithError(err).Error("Error creating receiver")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(receiver)
}

func (h *ReceiverHandler) RenderReceiverForm(w http.ResponseWriter, r *http.Request) {
	util.RenderTemplate(w, "receiver.html", nil)
}
