package services

import (
	"database/sql"
	"tally/internal/dto"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// setupFullTestDB creates an in-memory SQLite database with core tables for CRUD tests.
func setupFullTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	queries := []string{
		`CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			uuid TEXT UNIQUE NOT NULL,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			email TEXT,
			avatar_url TEXT,
			created_at TEXT DEFAULT (datetime('now')),
			last_login TEXT,
			settings_json TEXT DEFAULT '{}'
		);`,
		`CREATE TABLE clients (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			name TEXT NOT NULL,
			email TEXT,
			website TEXT,
			avatar TEXT,
			contact_person TEXT,
			address TEXT,
			currency TEXT DEFAULT 'USD',
			status TEXT DEFAULT 'active',
			notes TEXT,
			billing_company TEXT,
			billing_address TEXT,
			billing_city TEXT,
			billing_province TEXT,
			billing_postal_code TEXT,
			FOREIGN KEY(user_id) REFERENCES users(id)
		);`,
		`CREATE TABLE projects (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			client_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			description TEXT,
			hourly_rate REAL,
			currency TEXT,
			status TEXT DEFAULT 'active',
			deadline TEXT,
			tags TEXT,
			service_type TEXT,
			FOREIGN KEY(user_id) REFERENCES users(id),
			FOREIGN KEY(client_id) REFERENCES clients(id)
		);`,
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
			invoiced BOOLEAN DEFAULT 0,
			billing_mode TEXT NOT NULL DEFAULT 'hourly',
			manual_amount REAL,
			FOREIGN KEY(user_id) REFERENCES users(id),
			FOREIGN KEY(project_id) REFERENCES projects(id),
			FOREIGN KEY(invoice_id) REFERENCES invoices(id)
		);`,
		`CREATE TABLE invoices (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			client_id INTEGER NOT NULL,
			number TEXT UNIQUE,
			issue_date TEXT,
			due_date TEXT,
			subtotal REAL,
			tax_rate REAL,
			tax_amount REAL,
			total REAL,
			status TEXT,
			items_json TEXT,
			FOREIGN KEY(user_id) REFERENCES users(id),
			FOREIGN KEY(client_id) REFERENCES clients(id)
		);`,
		`CREATE TABLE user_preferences (
			user_id INTEGER PRIMARY KEY,
			currency TEXT DEFAULT 'USD',
			language TEXT DEFAULT 'en-US',
			theme TEXT DEFAULT 'light',
			timezone TEXT DEFAULT 'UTC',
			date_format TEXT DEFAULT '2006-01-02',
			module_overrides_json TEXT DEFAULT '{}',
			updated_at TEXT DEFAULT (datetime('now')),
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
		);`,
		`CREATE TABLE user_tax_settings (
			user_id INTEGER PRIMARY KEY,
			hst_registered INTEGER DEFAULT 0,
			hst_number TEXT,
			tax_enabled INTEGER DEFAULT 0,
			default_tax_rate REAL DEFAULT 0,
			expected_income TEXT,
			updated_at TEXT DEFAULT (datetime('now')),
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
		);`,
		`CREATE TABLE user_invoice_settings (
			user_id INTEGER PRIMARY KEY,
			sender_name TEXT,
			sender_company TEXT,
			sender_address TEXT,
			sender_phone TEXT,
			sender_email TEXT,
			sender_postal_code TEXT,
			default_terms TEXT DEFAULT 'Due upon receipt',
			default_message_template TEXT DEFAULT 'Thank you for your business.',
			updated_at TEXT DEFAULT (datetime('now')),
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
		);`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			t.Fatalf("failed to create table: %v", err)
		}
	}

	return db
}

func createTestUser(t *testing.T, auth *AuthService, username string) dto.UserOutput {
	t.Helper()
	user, err := auth.Register(dto.RegisterInput{Username: username, Password: "pwd"})
	if err != nil {
		t.Fatalf("failed to register user: %v", err)
	}
	return user
}
