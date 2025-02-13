package model

import "time"

type Donation struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	MealType   string    `json:"mealType"`
	Quantity   int       `json:"quantity"`
	ExpiryDate time.Time `json:"expiryDate"`
	Location   string    `json:"location"`
	Notes      string    `json:"notes"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
