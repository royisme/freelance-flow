// Package models defines database-backed domain models.
package models

// Project represents a client engagement tracked for billing.
type Project struct {
	ID          int      `json:"id"`
	ClientID    int      `json:"clientId"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	HourlyRate  float64  `json:"hourlyRate"`
	Currency    string   `json:"currency"`
	Status      string   `json:"status"` // active, archived, completed
	Deadline    string   `json:"deadline"`
	Tags        []string `json:"tags"` // Handled as pipe-delimited string in DB for simplicity
}
