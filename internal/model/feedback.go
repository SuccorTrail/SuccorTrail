package model

import "time"

type Feedback struct {
	ID         string    `json:"id"`
	DonationID string    `json:"donationId"`
	Quality    int       `json:"quality"`
	Comments   string    `json:"comments"`
	CreatedAt  time.Time `json:"createdAt"`
}
