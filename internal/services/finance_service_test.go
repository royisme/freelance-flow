package services

import (
	"database/sql"
	"testing"
	"time"

	"tally/internal/dto"
)

func setupFinanceTestDB(t *testing.T) *sql.DB {
	db := setupFullTestDB(t)

	queries := []string{
		`CREATE TABLE IF NOT EXISTS finance_accounts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			type TEXT NOT NULL,
			currency TEXT DEFAULT 'CAD',
			balance REAL DEFAULT 0,
			bank_name TEXT,
			created_at TEXT DEFAULT (datetime('now')),
			updated_at TEXT DEFAULT (datetime('now')),
			FOREIGN KEY(user_id) REFERENCES users(id)
		);`,
		`CREATE TABLE IF NOT EXISTS finance_categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			type TEXT NOT NULL,
			color TEXT,
			icon TEXT,
			created_at TEXT DEFAULT (datetime('now')),
			updated_at TEXT DEFAULT (datetime('now')),
			FOREIGN KEY(user_id) REFERENCES users(id)
		);`,
		`CREATE TABLE IF NOT EXISTS finance_transactions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			account_id INTEGER NOT NULL,
			category_id INTEGER,
			date TEXT NOT NULL,
			description TEXT NOT NULL,
			amount REAL NOT NULL,
			status TEXT DEFAULT 'pending',
			reference_id TEXT,
			created_at TEXT DEFAULT (datetime('now')),
			updated_at TEXT DEFAULT (datetime('now')),
			FOREIGN KEY(user_id) REFERENCES users(id),
			FOREIGN KEY(account_id) REFERENCES finance_accounts(id),
			FOREIGN KEY(category_id) REFERENCES finance_categories(id)
		);`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			t.Fatalf("failed to create finance table: %v", err)
		}
	}

	return db
}

func TestFinanceService_Accounts(t *testing.T) {
	db := setupFinanceTestDB(t)
	defer db.Close()

	authService := NewAuthService(db)
	user := createTestUser(t, authService, "finance_user")

	service := NewFinanceService(db)

	// Test CreateAccount
	input := dto.CreateAccountInput{
		Name:     "TD Checking",
		Type:     "checking",
		Currency: "CAD",
		Balance:  1000.50,
		BankName: "TD",
	}

	acc := service.CreateAccount(user.ID, input)
	if acc.ID == 0 {
		t.Error("expected account ID to be set")
	}
	if acc.Name != input.Name {
		t.Errorf("expected name %s, got %s", input.Name, acc.Name)
	}

	// Test GetAccounts
	accounts := service.GetAccounts(user.ID)
	if len(accounts) != 1 {
		t.Errorf("expected 1 account, got %d", len(accounts))
	}
	if accounts[0].ID != acc.ID {
		t.Errorf("expected account ID %d, got %d", acc.ID, accounts[0].ID)
	}

	// Test UpdateAccount
	updateInput := dto.UpdateAccountInput{
		ID:       acc.ID,
		Name:     "TD Checking Updated",
		Type:     "checking",
		Currency: "CAD",
		Balance:  2000.00,
		BankName: "TD",
	}
	updatedAcc := service.UpdateAccount(user.ID, updateInput)
	if updatedAcc.Name != "TD Checking Updated" {
		t.Errorf("expected updated name, got %s", updatedAcc.Name)
	}
	if updatedAcc.Balance != 2000.00 {
		t.Errorf("expected updated balance, got %f", updatedAcc.Balance)
	}

	// Test DeleteAccount
	service.DeleteAccount(user.ID, acc.ID)
	accounts = service.GetAccounts(user.ID)
	if len(accounts) != 0 {
		t.Errorf("expected 0 accounts after delete, got %d", len(accounts))
	}
}

func TestFinanceService_Categories(t *testing.T) {
	db := setupFinanceTestDB(t)
	defer db.Close()

	authService := NewAuthService(db)
	user := createTestUser(t, authService, "cat_user")
	service := NewFinanceService(db)

	input := dto.CreateCategoryInput{
		Name:  "Groceries",
		Type:  "expense",
		Color: "#FF0000",
		Icon:  "cart",
	}

	cat := service.CreateCategory(user.ID, input)
	if cat.ID == 0 {
		t.Error("expected category ID")
	}

	cats := service.GetCategories(user.ID)
	if len(cats) != 1 {
		t.Errorf("expected 1 category, got %d", len(cats))
	}

	service.DeleteCategory(user.ID, cat.ID)
	cats = service.GetCategories(user.ID)
	if len(cats) != 0 {
		t.Error("expected 0 categories")
	}
}

func TestFinanceService_ImportTransactions(t *testing.T) {
	db := setupFinanceTestDB(t)
	defer db.Close()

	authService := NewAuthService(db)
	user := createTestUser(t, authService, "import_user")
	service := NewFinanceService(db)

	// Create account
	accInput := dto.CreateAccountInput{Name: "TD", Type: "checking", Currency: "CAD"}
	acc := service.CreateAccount(user.ID, accInput)

	// CSV Content (TD style)
	// Date, Description, Withdrawal, Deposit, Balance
	csvContent := `01/15/2024,Walmart,-50.00,,1000.00
01/16/2024,Salary,,2000.00,3000.00`

	importInput := dto.ImportTransactionsInput{
		AccountID:   acc.ID,
		BankType:    "TD",
		FileContent: csvContent,
	}

	count, err := service.ImportTransactions(user.ID, importInput)
	if err != nil {
		t.Fatalf("import failed: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 imported transactions, got %d", count)
	}

	// Verify transactions
	filter := dto.TransactionFilter{AccountID: acc.ID}
	txs := service.GetTransactions(user.ID, filter)
	if len(txs) != 2 {
		t.Errorf("expected 2 transactions in DB, got %d", len(txs))
	}

	// Check details
	// Order is DESC by date. So Salary (16th) first, then Walmart (15th).
	if txs[0].Description != "Salary" {
		t.Errorf("expected first tx to be Salary, got %s", txs[0].Description)
	}
	if txs[0].Amount != 2000.00 {
		t.Errorf("expected Salary amount 2000, got %f", txs[0].Amount)
	}

	if txs[1].Description != "Walmart" {
		t.Errorf("expected second tx to be Walmart, got %s", txs[1].Description)
	}
	// TD parseRow logic: Withdrawal column is index 2.
	// row: 01/15/2024,Walmart,-50.00,,1000.00
	// row[2] = "-50.00". parseAmount("-50.00") -> -50.
	// Logic: if withdrawal > 0 { Amount = -withdrawal } else { Amount = withdrawal }
	// -50 is NOT > 0. So Amount = -50.
	// So it is correctly negative.
	if txs[1].Amount != -50.00 {
		t.Errorf("expected Walmart amount -50.00, got %f", txs[1].Amount)
	}
}

func TestFinanceService_Summary(t *testing.T) {
	db := setupFinanceTestDB(t)
	defer db.Close()

	authService := NewAuthService(db)
	user := createTestUser(t, authService, "summary_user")
	service := NewFinanceService(db)

	// Create account with balance 5000
	service.CreateAccount(user.ID, dto.CreateAccountInput{Name: "Acc1", Balance: 5000})

	// Add transactions manually (SQL) because Import is complex and we verified it.
	// Income: 2000
	// Expense: -500
	_, err := db.Exec("INSERT INTO finance_transactions(user_id, account_id, date, description, amount) VALUES(?, ?, ?, ?, ?)",
		user.ID, 1, time.Now().Format("2006-01-02"), "Income", 2000.00)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec("INSERT INTO finance_transactions(user_id, account_id, date, description, amount) VALUES(?, ?, ?, ?, ?)",
		user.ID, 1, time.Now().Format("2006-01-02"), "Expense", -500.00)
	if err != nil {
		t.Fatal(err)
	}

	summary := service.GetSummary(user.ID)
	if summary.TotalBalance != 5000 {
		t.Errorf("expected total balance 5000, got %f", summary.TotalBalance)
	}
	if summary.TotalIncome != 2000 {
		t.Errorf("expected total income 2000, got %f", summary.TotalIncome)
	}
	// Logic says if negative, flip it.
	// TotalExpense should be 500 (positive representation of negative sum)
	// Code: if summary.TotalExpense < 0 { summary.TotalExpense = -summary.TotalExpense }
	// SQL sum is -500. So result is 500.
	if summary.TotalExpense != 500 {
		t.Errorf("expected total expense 500, got %f", summary.TotalExpense)
	}
	// CashFlow = Income + Expense (signed) = 2000 + (-500) = 1500.
	if summary.CashFlow != 1500 {
		t.Errorf("expected cash flow 1500, got %f", summary.CashFlow)
	}
}
