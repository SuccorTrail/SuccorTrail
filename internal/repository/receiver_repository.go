package repository

import (
	"strings"
	"github.com/SuccorTrail/SuccorTrail/internal/db"
	"github.com/SuccorTrail/SuccorTrail/internal/model"
)

type ReceiverRepository interface {
	Create(receiver *model.Receiver) error
}

type receiverRepository struct{}

func NewReceiverRepository() ReceiverRepository {
	return &receiverRepository{}
}

func (r *receiverRepository) Create(receiver *model.Receiver) error {
	// Convert dietary restrictions to comma-separated string
	dietaryRestrictionsStr := strings.Join(receiver.DietaryRestrictions, ",")

	_, err := db.GetDB().Exec(
		"INSERT INTO receivers (id, name, phone, location, family_size, dietary_restrictions, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		receiver.ID, receiver.Name, receiver.Phone, receiver.Location, receiver.FamilySize, dietaryRestrictionsStr, receiver.CreatedAt)
	return err
}
