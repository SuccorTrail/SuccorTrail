package model

import "time"

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Don't expose in JSON
	UserType  string    `json:"userType"`
	CreatedAt time.Time `json:"createdAt"`
}
