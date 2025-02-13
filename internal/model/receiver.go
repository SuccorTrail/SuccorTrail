package model

import "time"

type Receiver struct {
	ID                  string    `json:"id"`
	Name                string    `json:"name"`
	Phone               string    `json:"phone"`
	Location            string    `json:"location"`
	FamilySize          int       `json:"familySize"`
	DietaryRestrictions []string  `json:"dietaryRestrictions"`
	CreatedAt           time.Time `json:"createdAt"`
}
