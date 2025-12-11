// Package models defines database-backed domain models.
package models

// InvoiceItem describes a single line item within an invoice.
type InvoiceItem struct {
	ID          int     `json:"id"`
	Description string  `json:"description"`
	Quantity    float64 `json:"quantity"`
	UnitPrice   float64 `json:"unitPrice"`
	Amount      float64 `json:"amount"`
}

// Invoice captures billing information for a client.
type Invoice struct {
	ID        int           `json:"id"`
	ClientID  int           `json:"clientId"`
	Number    string        `json:"number"`
	IssueDate string        `json:"issueDate"`
	DueDate   string        `json:"dueDate"`
	Subtotal  float64       `json:"subtotal"`
	TaxRate   float64       `json:"taxRate"`
	TaxAmount float64       `json:"taxAmount"`
	Total     float64       `json:"total"`
	Status    string        `json:"status"` // draft, sent, paid, overdue
	Items     []InvoiceItem `json:"items"`
}
