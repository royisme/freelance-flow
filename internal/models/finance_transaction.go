package models

import "time"

// FinanceTransaction represents a single financial transaction.
type FinanceTransaction struct {
	ID          int              `json:"id"`
	UserID      int              `json:"userId"`
	AccountID   int              `json:"accountId"`
	CategoryID  *int             `json:"categoryId"` // Nullable
	Category    *FinanceCategory `json:"category,omitempty"` // For joining
	Date        time.Time        `json:"date"`
	Description string           `json:"description"`
	Amount      float64          `json:"amount"`
	Status      string           `json:"status"` // pending, cleared, reconciled
	ReferenceID string           `json:"referenceId"`
	CreatedAt   time.Time        `json:"createdAt"`
	UpdatedAt   time.Time        `json:"updatedAt"`
}
