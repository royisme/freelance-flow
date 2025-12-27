package dto

import "time"

// CreateAccountInput represents the input to create a new account.
type CreateAccountInput struct {
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Currency string  `json:"currency"`
	Balance  float64 `json:"balance"`
	BankName string  `json:"bankName"`
}

// UpdateAccountInput represents the input to update an account.
type UpdateAccountInput struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Currency string  `json:"currency"`
	Balance  float64 `json:"balance"`
	BankName string  `json:"bankName"`
}

// AccountOutput represents the output for an account.
type AccountOutput struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Currency  string    `json:"currency"`
	Balance   float64   `json:"balance"`
	BankName  string    `json:"bankName"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CreateCategoryInput represents the input to create a new category.
type CreateCategoryInput struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Color string `json:"color"`
	Icon  string `json:"icon"`
}

// UpdateCategoryInput represents the input to update a category.
type UpdateCategoryInput struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Color string `json:"color"`
	Icon  string `json:"icon"`
}

// CategoryOutput represents the output for a category.
type CategoryOutput struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Color string `json:"color"`
	Icon  string `json:"icon"`
}

// TransactionOutput represents the output for a transaction.
type TransactionOutput struct {
	ID           int             `json:"id"`
	AccountID    int             `json:"accountId"`
	CategoryID   *int            `json:"categoryId"`
	CategoryName string          `json:"categoryName,omitempty"`
	CategoryColor string         `json:"categoryColor,omitempty"`
	Date         time.Time       `json:"date"`
	Description  string          `json:"description"`
	Amount       float64         `json:"amount"`
	Status       string          `json:"status"`
	ReferenceID  string          `json:"referenceId"`
}

// ImportTransactionsInput represents the input to import transactions.
type ImportTransactionsInput struct {
	AccountID   int    `json:"accountId"`
	BankType    string `json:"bankType"` // CIBC, RBC, TD
	FileContent string `json:"fileContent"` // Base64 or raw string
}

// TransactionFilter represents filtering options for transactions.
type TransactionFilter struct {
	StartDate string `json:"startDate,omitempty"`
	EndDate   string `json:"endDate,omitempty"`
	AccountID int    `json:"accountId,omitempty"`
}

// FinanceSummary represents the dashboard summary.
type FinanceSummary struct {
	TotalBalance float64 `json:"totalBalance"`
	TotalIncome  float64 `json:"totalIncome"`
	TotalExpense float64 `json:"totalExpense"`
	CashFlow     float64 `json:"cashFlow"`
}
