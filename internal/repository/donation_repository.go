package repository

import (
	"github.com/SuccorTrail/SuccorTrail/internal/db"
	"github.com/SuccorTrail/SuccorTrail/internal/model"
)

type DonationRepository interface {
	Create(donation *model.Donation) error
}

type donationRepository struct{}

func NewDonationRepository() DonationRepository {
	return &donationRepository{}
}

func (r *donationRepository) Create(donation *model.Donation) error {
	_, err := db.GetDB().Exec(
		"INSERT INTO donations (id, name, email, phone, meal_type, quantity, expiry_date, location, notes, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		donation.ID, donation.Name, donation.Email, donation.Phone, donation.MealType, donation.Quantity, donation.ExpiryDate, donation.Location, donation.Notes, donation.CreatedAt, donation.UpdatedAt)
	return err
}
