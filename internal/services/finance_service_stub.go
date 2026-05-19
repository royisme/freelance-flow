//go:build FINANCE_MODULE_DISABLED

package services

import (
	"database/sql"
	"tally/internal/dto"
)

// FinanceService handles all finance-related operations (disabled stub).
type FinanceService struct {
	db *sql.DB
}

// NewFinanceService creates a new FinanceService instance (disabled stub).
func NewFinanceService(db *sql.DB) *FinanceService {
	return &FinanceService{db: db}
}

// GetAccounts returns empty accounts when finance module is disabled.
func (s *FinanceService) GetAccounts(userID int) []dto.AccountOutput {
	return []dto.AccountOutput{}
}

// CreateAccount returns empty account when finance module is disabled.
func (s *FinanceService) CreateAccount(userID int, input dto.CreateAccountInput) dto.AccountOutput {
	return dto.AccountOutput{}
}

// UpdateAccount returns empty account when finance module is disabled.
func (s *FinanceService) UpdateAccount(userID int, input dto.UpdateAccountInput) dto.AccountOutput {
	return dto.AccountOutput{}
}

// DeleteAccount does nothing when finance module is disabled.
func (s *FinanceService) DeleteAccount(userID int, id int) {
}

// GetCategories returns empty categories when finance module is disabled.
func (s *FinanceService) GetCategories(userID int) []dto.CategoryOutput {
	return []dto.CategoryOutput{}
}

// CreateCategory returns empty category when finance module is disabled.
func (s *FinanceService) CreateCategory(userID int, input dto.CreateCategoryInput) dto.CategoryOutput {
	return dto.CategoryOutput{}
}

// UpdateCategory returns empty category when finance module is disabled.
func (s *FinanceService) UpdateCategory(userID int, input dto.UpdateCategoryInput) dto.CategoryOutput {
	return dto.CategoryOutput{}
}

// DeleteCategory does nothing when finance module is disabled.
func (s *FinanceService) DeleteCategory(userID int, id int) {
}

// GetTransactions returns empty transactions when finance module is disabled.
func (s *FinanceService) GetTransactions(userID int, filter dto.TransactionFilter) []dto.TransactionOutput {
	return []dto.TransactionOutput{}
}

// UpdateTransaction does nothing when finance module is disabled.
func (s *FinanceService) UpdateTransaction(userID int, transactionID int, categoryID *int) {
}

// DeleteTransaction does nothing when finance module is disabled.
func (s *FinanceService) DeleteTransaction(userID int, id int) {
}

// ImportTransactions returns 0 and error when finance module is disabled.
func (s *FinanceService) ImportTransactions(userID int, input dto.ImportTransactionsInput) (int, error) {
	return 0, nil
}

// GetSummary returns empty summary when finance module is disabled.
func (s *FinanceService) GetSummary(userID int) dto.FinanceSummary {
	return dto.FinanceSummary{}
}
