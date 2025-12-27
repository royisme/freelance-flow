package models

import "time"

// FinanceAccount represents a bank account or credit card.
type FinanceAccount struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	Name      string    `json:"name"`
	Type      string    `json:"type"` // checking, savings, credit, investment
	Currency  string    `json:"currency"`
	Balance   float64   `json:"balance"`
	BankName  string    `json:"bankName"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
