package repository

import (
	"github.com/SuccorTrail/SuccorTrail/internal/db"
	"github.com/SuccorTrail/SuccorTrail/internal/model"
)

type MealRepository interface {
	GetAvailableMeals(location string) ([]model.Donation, error)
}

type mealRepository struct{}

func NewMealRepository() MealRepository {
	return &mealRepository{}
}

func (r *mealRepository) GetAvailableMeals(location string) ([]model.Donation, error) {
	rows, err := db.GetDB().Query("SELECT id, name, email, phone, meal_type, quantity, expiry_date, location, notes, created_at, updated_at FROM donations WHERE location = ?", location)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meals []model.Donation
	for rows.Next() {
		var meal model.Donation
		err := rows.Scan(&meal.ID, &meal.Name, &meal.Email, &meal.Phone, &meal.MealType, &meal.Quantity, &meal.ExpiryDate, &meal.Location, &meal.Notes, &meal.CreatedAt, &meal.UpdatedAt)
		if err != nil {
			return nil, err
		}
		meals = append(meals, meal)
	}
	return meals, nil
}
