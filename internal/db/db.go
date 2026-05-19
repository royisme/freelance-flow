// Package db initializes and migrates the application database.
package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	// sqlite driver used by the application.
	_ "github.com/mattn/go-sqlite3"
)

// Init sets up the SQLite database, creating directories and running migrations.
func Init() *sql.DB {
	// Ensure data directory exists
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Failed to get user home directory:", err)
	}

	appDataDir := filepath.Join(homeDir, ".tally")
	if err := os.MkdirAll(appDataDir, 0700); err != nil {
		log.Fatal("Failed to create app data directory:", err)
	}

	dbPath := filepath.Join(appDataDir, "tally.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	// Try to run golang-migrate first
	if err := RunMigrations(db); err != nil {
		log.Printf("golang-migrate failed, falling back to legacy migrations: %v", err)
		// Fallback to legacy table creation and migrations for existing DBs
		createTables(db)
		runLegacyMigrations(db)
	}

	return db
}

// Open opens the application database without running migrations.
func Open() (*sql.DB, string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, "", err
	}

	appDataDir := filepath.Join(homeDir, ".tally")
	if err := os.MkdirAll(appDataDir, 0700); err != nil {
		return nil, "", err
	}

	dbPath := filepath.Join(appDataDir, "tally.db")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, "", err
	}

	return db, dbPath, nil
}

func createTables(db *sql.DB) {
	queries := []string{
		// Users table - must be created first for foreign key references
		`CREATE TABLE IF NOT EXISTS users (
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
		`CREATE TABLE IF NOT EXISTS clients (
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
		`CREATE TABLE IF NOT EXISTS projects (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			client_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			description TEXT,
			hourly_rate REAL,
			currency TEXT,
			status TEXT DEFAULT 'active',
			deadline TEXT,
			tags TEXT, -- Comma separated
			service_type TEXT,
			FOREIGN KEY(user_id) REFERENCES users(id),
			FOREIGN KEY(client_id) REFERENCES clients(id)
		);`,
		`CREATE TABLE IF NOT EXISTS time_entries (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			project_id INTEGER NOT NULL,
			date TEXT,
			start_time TEXT,
			end_time TEXT,
			duration_seconds INTEGER,
			description TEXT,
			billable BOOLEAN DEFAULT 1,
			invoiced BOOLEAN DEFAULT 0,
			FOREIGN KEY(user_id) REFERENCES users(id),
			FOREIGN KEY(project_id) REFERENCES projects(id)
		);`,
		`CREATE TABLE IF NOT EXISTS invoices (
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
			items_json TEXT, -- Store items as JSON for MVP simplicity
			FOREIGN KEY(user_id) REFERENCES users(id),
			FOREIGN KEY(client_id) REFERENCES clients(id)
		);`,
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			log.Fatalf("Failed to create table: %v\nQuery: %s", err, query)
		}
	}
}

// runLegacyMigrations handles database schema migrations for existing databases.
// This is kept for backward compatibility with databases created before golang-migrate.
func runLegacyMigrations(db *sql.DB) {
	// Migration: Add billable column to time_entries if not exists
	// SQLite doesn't support IF NOT EXISTS for columns, so we check first
	addColumnIfNotExists(db, "time_entries", "billable", "BOOLEAN DEFAULT 1")

	// Migration: Add user_id columns to all entity tables
	addColumnIfNotExists(db, "clients", "user_id", "INTEGER")
	addColumnIfNotExists(db, "projects", "user_id", "INTEGER")
	addColumnIfNotExists(db, "time_entries", "user_id", "INTEGER")
	addColumnIfNotExists(db, "invoices", "user_id", "INTEGER")

	// Migration: Add billing fields to clients
	addColumnIfNotExists(db, "clients", "billing_company", "TEXT")
	addColumnIfNotExists(db, "clients", "billing_address", "TEXT")
	addColumnIfNotExists(db, "clients", "billing_city", "TEXT")
	addColumnIfNotExists(db, "clients", "billing_province", "TEXT")
	addColumnIfNotExists(db, "clients", "billing_postal_code", "TEXT")

	// Migration: Add service_type to projects
	addColumnIfNotExists(db, "projects", "service_type", "TEXT")
}

// addColumnIfNotExists adds a column to a table if it doesn't already exist.
func addColumnIfNotExists(db *sql.DB, table, column, columnDef string) {
	var count int
	query := "SELECT COUNT(*) FROM pragma_table_info('" + table + "') WHERE name='" + column + "'"
	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		log.Printf("Error checking for %s.%s column: %v", table, column, err)
		return
	}
	if count == 0 {
		alterQuery := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", table, column, columnDef)
		_, err = db.Exec(alterQuery)
		if err != nil {
			log.Printf("Error adding %s.%s column: %v", table, column, err)
		} else {
			log.Printf("Migration: Added %s column to %s", column, table)
		}
	}
}
