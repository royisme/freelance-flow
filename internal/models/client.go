// Package models defines database-backed domain models.
package models

// Client represents a customer's organization or individual.
type Client struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Website       string `json:"website"`
	Avatar        string `json:"avatar"`
	ContactPerson string `json:"contactPerson"`
	Address       string `json:"address"`
	Currency      string `json:"currency"`
	Status        string `json:"status"` // active, inactive
	Notes         string `json:"notes"`
}
