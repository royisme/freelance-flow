package services

import (
	"database/sql"
	"freelance-flow/internal/dto"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

// Ensures SetTimeEntries recalculates totals from linked time entries.
func TestInvoiceService_SetTimeEntries_RecalculateTotals(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	defer db.Close()

	schema := []string{
		`CREATE TABLE users (id INTEGER PRIMARY KEY AUTOINCREMENT, uuid TEXT, username TEXT, password_hash TEXT, settings_json TEXT DEFAULT '{}');`,
		`CREATE TABLE clients (id INTEGER PRIMARY KEY AUTOINCREMENT, user_id INTEGER, name TEXT);`,
		`CREATE TABLE projects (id INTEGER PRIMARY KEY AUTOINCREMENT, user_id INTEGER, client_id INTEGER, name TEXT, hourly_rate REAL, currency TEXT);`,
		`CREATE TABLE invoices (id INTEGER PRIMARY KEY AUTOINCREMENT, user_id INTEGER, client_id INTEGER, number TEXT, issue_date TEXT, due_date TEXT, subtotal REAL, tax_rate REAL, tax_amount REAL, total REAL, status TEXT, items_json TEXT);`,
		`CREATE TABLE time_entries (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			project_id INTEGER NOT NULL,
			invoice_id INTEGER,
			date TEXT,
			start_time TEXT,
			end_time TEXT,
			duration_seconds INTEGER,
			description TEXT,
			billable BOOLEAN DEFAULT 1,
			invoiced BOOLEAN DEFAULT 0
		);`,
	}
	for _, q := range schema {
		_, err := db.Exec(q)
		assert.NoError(t, err)
	}

	invoiceSvc := NewInvoiceService(db)

	// seed data
	_, err = db.Exec(`INSERT INTO users(id, username, password_hash) VALUES (1, 'user', 'pw')`)
	assert.NoError(t, err)
	_, err = db.Exec(`INSERT INTO clients(id, user_id, name) VALUES (1, 1, 'Client A')`)
	assert.NoError(t, err)
	_, err = db.Exec(`INSERT INTO projects(id, user_id, client_id, name, hourly_rate, currency) VALUES (1, 1, 1, 'Project A', 100, 'USD')`)
	assert.NoError(t, err)
	_, err = db.Exec(`INSERT INTO invoices(id, user_id, client_id, number, issue_date, due_date, tax_rate, subtotal, tax_amount, total, status, items_json) VALUES (1, 1, 1, 'INV-1', '2023-01-01', '2023-01-15', 0.1, 0, 0, 0, 'draft', '[]')`)
	assert.NoError(t, err)

	// two time entries: 1h and 2h
	_, err = db.Exec(`INSERT INTO time_entries(user_id, project_id, date, duration_seconds, billable, invoiced) VALUES (1, 1, '2023-01-01', 3600, 1, 0)`)
	assert.NoError(t, err)
	_, err = db.Exec(`INSERT INTO time_entries(user_id, project_id, date, duration_seconds, billable, invoiced) VALUES (1, 1, '2023-01-02', 7200, 1, 0)`)
	assert.NoError(t, err)

	out, err := invoiceSvc.SetTimeEntries(1, dto.SetInvoiceTimeEntriesInput{
		InvoiceID:    1,
		TimeEntryIDs: []int{1, 2},
	})
	assert.NoError(t, err)

	// subtotal = (3 hours * 100), tax 10%, total = 330
	assert.InDelta(t, 300.0, out.Subtotal, 0.001)
	assert.InDelta(t, 30.0, out.TaxAmount, 0.001)
	assert.InDelta(t, 330.0, out.Total, 0.001)
	assert.Len(t, out.Items, 1)
	assert.InDelta(t, 300.0, out.Items[0].Amount, 0.001)

	// ensure time entries flagged and linked
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM time_entries WHERE invoice_id = 1 AND invoiced = 1").Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 2, count)
}
