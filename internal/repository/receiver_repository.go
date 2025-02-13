package repository

import (
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
	_, err := db.GetDB().Exec(
		"INSERT INTO receivers (id, name, phone, location, family_size, dietary_restrictions, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		receiver.ID, receiver.Name, receiver.Phone, receiver.Location, receiver.FamilySize, receiver.DietaryRestrictions, receiver.CreatedAt)
	return err
}
