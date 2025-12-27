package services

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"tally/internal/dto"
	"tally/internal/models"
)

// FinanceService handles all finance-related operations.
type FinanceService struct {
	db *sql.DB
}

// NewFinanceService creates a new FinanceService instance.
func NewFinanceService(db *sql.DB) *FinanceService {
	return &FinanceService{db: db}
}

// --- Accounts ---

// GetAccounts returns all accounts for a user.
func (s *FinanceService) GetAccounts(userID int) []dto.AccountOutput {
	rows, err := s.db.Query("SELECT id, name, type, currency, balance, bank_name, updated_at FROM finance_accounts WHERE user_id = ?", userID)
	if err != nil {
		log.Println("Error querying finance accounts:", err)
		return []dto.AccountOutput{}
	}
	defer rows.Close()

	var accounts []dto.AccountOutput
	for rows.Next() {
		var acc dto.AccountOutput
		var bankName sql.NullString
		var updatedAtStr string
		err := rows.Scan(&acc.ID, &acc.Name, &acc.Type, &acc.Currency, &acc.Balance, &bankName, &updatedAtStr)
		if err != nil {
			log.Println("Error scanning finance account:", err)
			continue
		}
		acc.BankName = bankName.String
		if t, err := time.Parse("2006-01-02 15:04:05", updatedAtStr); err == nil {
			acc.UpdatedAt = t
		}
		accounts = append(accounts, acc)
	}
	return accounts
}

// CreateAccount creates a new account.
func (s *FinanceService) CreateAccount(userID int, input dto.CreateAccountInput) dto.AccountOutput {
	stmt, err := s.db.Prepare("INSERT INTO finance_accounts(user_id, name, type, currency, balance, bank_name) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("Error preparing insert account:", err)
		return dto.AccountOutput{}
	}
	defer stmt.Close()

	res, err := stmt.Exec(userID, input.Name, input.Type, input.Currency, input.Balance, input.BankName)
	if err != nil {
		log.Println("Error inserting finance account:", err)
		return dto.AccountOutput{}
	}

	id, _ := res.LastInsertId()
	return dto.AccountOutput{
		ID:       int(id),
		Name:     input.Name,
		Type:     input.Type,
		Currency: input.Currency,
		Balance:  input.Balance,
		BankName: input.BankName,
	}
}

// UpdateAccount updates an existing account.
func (s *FinanceService) UpdateAccount(userID int, input dto.UpdateAccountInput) dto.AccountOutput {
	stmt, err := s.db.Prepare("UPDATE finance_accounts SET name=?, type=?, currency=?, balance=?, bank_name=?, updated_at=datetime('now') WHERE id=? AND user_id=?")
	if err != nil {
		log.Println("Error preparing update account:", err)
		return dto.AccountOutput{}
	}
	defer stmt.Close()

	_, err = stmt.Exec(input.Name, input.Type, input.Currency, input.Balance, input.BankName, input.ID, userID)
	if err != nil {
		log.Println("Error updating finance account:", err)
		return dto.AccountOutput{}
	}

	return dto.AccountOutput{
		ID:       input.ID,
		Name:     input.Name,
		Type:     input.Type,
		Currency: input.Currency,
		Balance:  input.Balance,
		BankName: input.BankName,
	}
}

// DeleteAccount deletes an account.
func (s *FinanceService) DeleteAccount(userID int, id int) {
	// First delete transactions associated with this account to avoid constraints if any
	// (Though foreign key cascade might handle this, SQLite needs PRAGMA foreign_keys=ON)
	_, _ = s.db.Exec("DELETE FROM finance_transactions WHERE account_id=? AND user_id=?", id, userID)

	_, err := s.db.Exec("DELETE FROM finance_accounts WHERE id=? AND user_id=?", id, userID)
	if err != nil {
		log.Println("Error deleting finance account:", err)
	}
}

// --- Categories ---

// GetCategories returns all categories for a user.
func (s *FinanceService) GetCategories(userID int) []dto.CategoryOutput {
	rows, err := s.db.Query("SELECT id, name, type, color, icon FROM finance_categories WHERE user_id = ?", userID)
	if err != nil {
		log.Println("Error querying finance categories:", err)
		return []dto.CategoryOutput{}
	}
	defer rows.Close()

	var categories []dto.CategoryOutput
	for rows.Next() {
		var c dto.CategoryOutput
		var color, icon sql.NullString
		err := rows.Scan(&c.ID, &c.Name, &c.Type, &color, &icon)
		if err != nil {
			log.Println("Error scanning finance category:", err)
			continue
		}
		c.Color = color.String
		c.Icon = icon.String
		categories = append(categories, c)
	}
	return categories
}

// CreateCategory creates a new category.
func (s *FinanceService) CreateCategory(userID int, input dto.CreateCategoryInput) dto.CategoryOutput {
	stmt, err := s.db.Prepare("INSERT INTO finance_categories(user_id, name, type, color, icon) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("Error preparing insert category:", err)
		return dto.CategoryOutput{}
	}
	defer stmt.Close()

	res, err := stmt.Exec(userID, input.Name, input.Type, input.Color, input.Icon)
	if err != nil {
		log.Println("Error inserting finance category:", err)
		return dto.CategoryOutput{}
	}

	id, _ := res.LastInsertId()
	return dto.CategoryOutput{
		ID:    int(id),
		Name:  input.Name,
		Type:  input.Type,
		Color: input.Color,
		Icon:  input.Icon,
	}
}

// UpdateCategory updates an existing category.
func (s *FinanceService) UpdateCategory(userID int, input dto.UpdateCategoryInput) dto.CategoryOutput {
	stmt, err := s.db.Prepare("UPDATE finance_categories SET name=?, type=?, color=?, icon=?, updated_at=datetime('now') WHERE id=? AND user_id=?")
	if err != nil {
		log.Println("Error preparing update category:", err)
		return dto.CategoryOutput{}
	}
	defer stmt.Close()

	_, err = stmt.Exec(input.Name, input.Type, input.Color, input.Icon, input.ID, userID)
	if err != nil {
		log.Println("Error updating finance category:", err)
		return dto.CategoryOutput{}
	}

	return dto.CategoryOutput{
		ID:    input.ID,
		Name:  input.Name,
		Type:  input.Type,
		Color: input.Color,
		Icon:  input.Icon,
	}
}

// DeleteCategory deletes a category.
func (s *FinanceService) DeleteCategory(userID int, id int) {
	// Set transactions with this category to null
	_, _ = s.db.Exec("UPDATE finance_transactions SET category_id=NULL WHERE category_id=? AND user_id=?", id, userID)

	_, err := s.db.Exec("DELETE FROM finance_categories WHERE id=? AND user_id=?", id, userID)
	if err != nil {
		log.Println("Error deleting finance category:", err)
	}
}

// --- Transactions ---

// GetTransactions returns filtered transactions.
func (s *FinanceService) GetTransactions(userID int, filter dto.TransactionFilter) []dto.TransactionOutput {
	query := `
		SELECT t.id, t.account_id, t.category_id, c.name, c.color, t.date, t.description, t.amount, t.status, t.reference_id
		FROM finance_transactions t
		LEFT JOIN finance_categories c ON t.category_id = c.id
		WHERE t.user_id = ?
	`
	args := []interface{}{userID}

	if filter.AccountID > 0 {
		query += " AND t.account_id = ?"
		args = append(args, filter.AccountID)
	}

	if filter.StartDate != "" {
		query += " AND t.date >= ?"
		args = append(args, filter.StartDate)
	}

	if filter.EndDate != "" {
		query += " AND t.date <= ?"
		args = append(args, filter.EndDate)
	}

	query += " ORDER BY t.date DESC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		log.Println("Error querying transactions:", err)
		return []dto.TransactionOutput{}
	}
	defer rows.Close()

	var transactions []dto.TransactionOutput
	for rows.Next() {
		var t dto.TransactionOutput
		var catID sql.NullInt64
		var catName, catColor, refID sql.NullString
		var dateStr string

		err := rows.Scan(&t.ID, &t.AccountID, &catID, &catName, &catColor, &dateStr, &t.Description, &t.Amount, &t.Status, &refID)
		if err != nil {
			log.Println("Error scanning transaction:", err)
			continue
		}

		if catID.Valid {
			id := int(catID.Int64)
			t.CategoryID = &id
			t.CategoryName = catName.String
			t.CategoryColor = catColor.String
		}
		t.ReferenceID = refID.String
		if date, err := time.Parse("2006-01-02", dateStr); err == nil {
			t.Date = date
		} else if date, err := time.Parse("2006-01-02 15:04:05", dateStr); err == nil {
			t.Date = date
		} else if date, err := time.Parse(time.RFC3339, dateStr); err == nil {
			t.Date = date
		}

		transactions = append(transactions, t)
	}
	return transactions
}

// UpdateTransaction updates a transaction (mainly for categorization).
func (s *FinanceService) UpdateTransaction(userID int, transactionID int, categoryID *int) {
	query := "UPDATE finance_transactions SET category_id=?, updated_at=datetime('now') WHERE id=? AND user_id=?"
	_, err := s.db.Exec(query, categoryID, transactionID, userID)
	if err != nil {
		log.Println("Error updating transaction:", err)
	}
}

// DeleteTransaction deletes a transaction.
func (s *FinanceService) DeleteTransaction(userID int, id int) {
	_, err := s.db.Exec("DELETE FROM finance_transactions WHERE id=? AND user_id=?", id, userID)
	if err != nil {
		log.Println("Error deleting transaction:", err)
	}
}

// ImportTransactions parses and saves transactions from CSV.
func (s *FinanceService) ImportTransactions(userID int, input dto.ImportTransactionsInput) (int, error) {
	reader := csv.NewReader(strings.NewReader(input.FileContent))
	records, err := reader.ReadAll()
	if err != nil {
		return 0, fmt.Errorf("failed to parse CSV: %v", err)
	}

	if len(records) < 1 {
		return 0, fmt.Errorf("empty CSV file")
	}

	var transactions []models.FinanceTransaction
	bankType := strings.ToUpper(input.BankType)

	// Simple logic: Skip header usually.
	// But we need to identify columns.

	for i, row := range records {
		// Heuristic: Skip header row if it contains "Date" or "Description"
		if i == 0 {
			rowStr := strings.Join(row, ",")
			if strings.Contains(strings.ToLower(rowStr), "date") || strings.Contains(strings.ToLower(rowStr), "description") {
				continue
			}
		}

		t, err := s.parseRow(row, bankType)
		if err != nil {
			// Log error but continue? or fail?
			// For now, log and skip
			log.Printf("Skipping row %d: %v", i, err)
			continue
		}
		t.UserID = userID
		t.AccountID = input.AccountID
		transactions = append(transactions, t)
	}

	count := 0
	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO finance_transactions(user_id, account_id, date, description, amount, status, reference_id) VALUES(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	for _, t := range transactions {
		// Deduplication check: Same Date, Description, Amount
		// Note: This is a weak check. ReferenceID is better if available.
		// For MVP, we just insert. Or we can check existence.
		// Let's implement a simple existence check.
		var exists int
		checkQuery := "SELECT COUNT(*) FROM finance_transactions WHERE user_id=? AND account_id=? AND date=? AND description=? AND amount=?"
		err := tx.QueryRow(checkQuery, userID, input.AccountID, t.Date.Format("2006-01-02"), t.Description, t.Amount).Scan(&exists)
		if err == nil && exists > 0 {
			continue
		}

		_, err = stmt.Exec(userID, input.AccountID, t.Date.Format("2006-01-02"), t.Description, t.Amount, "imported", t.ReferenceID)
		if err != nil {
			log.Printf("Failed to insert transaction: %v", err)
			continue
		}
		count++
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return count, nil
}

func (s *FinanceService) parseRow(row []string, bankType string) (models.FinanceTransaction, error) {
	var t models.FinanceTransaction
	var err error

	switch bankType {
	case "CIBC":
		// CIBC Export: Date, Description, Debit, Credit, CardName, CardNumber
		// Example: 2024-01-01,AMAZON.CA,-10.00,,...
		// OR: 2024-01-01,DEPOSIT,,,1000.00,...
		if len(row) < 4 {
			return t, fmt.Errorf("not enough columns for CIBC")
		}

		t.Date, err = parseDate(row[0])
		if err != nil {
			return t, err
		}
		t.Description = row[1]

		debit := parseAmount(row[2])
		credit := parseAmount(row[3])

		if credit > 0 {
			t.Amount = credit // Income
		} else if debit != 0 {
			// CIBC debits might be positive or negative in CSV?
			// Usually in CIBC CSV: "Debit" column contains the amount spent.
			// "Credit" column contains amount received.
			// Let's assume Debit is positive number for expense -> we make it negative.
			if debit > 0 {
				t.Amount = -debit
			} else {
				t.Amount = debit // If already negative
			}
		}

	case "RBC":
		// RBC Export: Account Type, Account Number, Transaction Date, Cheque Number, Description 1, Description 2, CAD$, USD$
		// Index: 0=Type, 1=Num, 2=Date, 3=Cheque, 4=Desc1, 5=Desc2, 6=CAD, 7=USD
		if len(row) < 7 {
			return t, fmt.Errorf("not enough columns for RBC")
		}

		t.Date, err = parseDate(row[2])
		if err != nil {
			return t, err
		}
		t.Description = row[4] + " " + row[5]
		t.Amount = parseAmount(row[6]) // RBC uses negative for expense

	case "TD":
		// TD Export: Date, Description, Withdrawal, Deposit, Balance
		if len(row) < 4 {
			return t, fmt.Errorf("not enough columns for TD")
		}

		t.Date, err = parseDate(row[0])
		if err != nil {
			return t, err
		}
		t.Description = row[1]

		withdrawal := parseAmount(row[2])
		deposit := parseAmount(row[3])

		if deposit > 0 {
			t.Amount = deposit
		} else {
			if withdrawal > 0 {
				t.Amount = -withdrawal
			} else {
				t.Amount = withdrawal
			}
		}

	default:
		// Generic fallback: Date, Description, Amount
		if len(row) < 3 {
			return t, fmt.Errorf("not enough columns for Generic")
		}
		t.Date, err = parseDate(row[0])
		if err != nil {
			return t, err
		}
		t.Description = row[1]
		t.Amount = parseAmount(row[2])
	}

	return t, nil
}

func parseDate(s string) (time.Time, error) {
	formats := []string{"2006-01-02", "01/02/2006", "1/2/2006", "2006/01/02"}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid date format: %s", s)
}

func parseAmount(s string) float64 {
	if s == "" {
		return 0
	}
	// Remove currency symbols and commas
	clean := strings.ReplaceAll(s, "$", "")
	clean = strings.ReplaceAll(clean, ",", "")
	clean = strings.TrimSpace(clean)
	val, _ := strconv.ParseFloat(clean, 64)
	return val
}

// GetSummary returns finance summary.
func (s *FinanceService) GetSummary(userID int) dto.FinanceSummary {
	var summary dto.FinanceSummary

	// Total Balance from Accounts
	row := s.db.QueryRow("SELECT COALESCE(SUM(balance), 0) FROM finance_accounts WHERE user_id = ?", userID)
	_ = row.Scan(&summary.TotalBalance)

	// Income and Expense from Transactions (last 30 days? or all time?)
	// Let's do All Time for now, or maybe month to date?
	// User request didn't specify, but dashboard usually implies current context.
	// Let's do ALL TIME for simplicty of the SQL, or maybe filtered by current month.
	// Let's do All Time to match the Total Balance logic.

	row = s.db.QueryRow("SELECT COALESCE(SUM(amount), 0) FROM finance_transactions WHERE user_id = ? AND amount > 0", userID)
	_ = row.Scan(&summary.TotalIncome)

	row = s.db.QueryRow("SELECT COALESCE(SUM(amount), 0) FROM finance_transactions WHERE user_id = ? AND amount < 0", userID)
	_ = row.Scan(&summary.TotalExpense)

	// Expense is usually shown as positive number in UI if labeled "Expense", but here we return raw sums.
	// Actually, let's keep it signed. expense will be negative.

	summary.CashFlow = summary.TotalIncome + summary.TotalExpense

	// If expense is negative, we might want to flip it for display if the UI expects positive "Total Expense"
	// But let's leave it to frontend.
	// Wait, the UI mock shows "Total Expense" separately. Usually people expect to see "$500" not "$-500".
	if summary.TotalExpense < 0 {
		summary.TotalExpense = -summary.TotalExpense
	}

	return summary
}
