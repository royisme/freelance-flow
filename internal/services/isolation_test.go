package services

import (
	"database/sql"
	"freelance-flow/internal/dto"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// setupFullTestDB creates an in-memory SQLite database with ALL tables.
func setupFullTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
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
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			t.Fatalf("Failed to create table: %v", err)
		}
	}

	return db
}

func TestDataIsolation(t *testing.T) {
	db := setupFullTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close db: %v", err)
		}
	}()

	authService := NewAuthService(db)
	clientService := NewClientService(db)
	projectService := NewProjectService(db)
	timesheetService := NewTimesheetService(db)
	invoiceService := NewInvoiceService(db)

	// 1. Create two users
	userA, _ := authService.Register(dto.RegisterInput{Username: "UserA", Password: "pwd"})
	userB, _ := authService.Register(dto.RegisterInput{Username: "UserB", Password: "pwd"})

	// 2. User A creates a Client
	clientInput := dto.CreateClientInput{
		Name: "Client A",
	}
	createdClient := clientService.Create(userA.ID, clientInput)

	// Verify User A can see it
	clientsA := clientService.List(userA.ID)
	if len(clientsA) != 1 {
		t.Errorf("User A should see 1 client, saw %d", len(clientsA))
	}

	// Verify User B sees NOTHING
	clientsB := clientService.List(userB.ID)
	if len(clientsB) != 0 {
		t.Errorf("User B should see 0 clients, saw %d", len(clientsB))
	}

	// Verify User B cannot Get it by ID
	_, err := clientService.Get(userB.ID, createdClient.ID)
	if err == nil {
		t.Error("User B should not be able to get User A's client")
	}

	// 3. User A creates a Project for that Client
	projectInput := dto.CreateProjectInput{
		ClientID:   createdClient.ID,
		Name:       "Project A",
		HourlyRate: 100,
	}
	createdProject := projectService.Create(userA.ID, projectInput)

	// Verify User A can see it
	projectsA := projectService.List(userA.ID)
	if len(projectsA) != 1 {
		t.Errorf("User A should see 1 project, saw %d", len(projectsA))
	}

	// Verify User B sees NOTHING
	projectsB := projectService.List(userB.ID)
	if len(projectsB) != 0 {
		t.Errorf("User B should see 0 projects, saw %d", len(projectsB))
	}

	// Verify User B cannot Get it by ID
	_, err = projectService.Get(userB.ID, createdProject.ID)
	if err == nil {
		t.Error("User B should not be able to get User A's project")
	}

	// 4. User A creates a Time Entry
	timeEntryInput := dto.CreateTimeEntryInput{
		ProjectID:       createdProject.ID,
		Date:            time.Now().Format("2006-01-02"),
		DurationSeconds: 3600,
	}
	createdTimeEntry := timesheetService.Create(userA.ID, timeEntryInput)

	// Verify User A can see it
	entriesA := timesheetService.List(userA.ID, 0)
	if len(entriesA) != 1 {
		t.Errorf("User A should see 1 time entry, saw %d", len(entriesA))
	}

	// Verify User B sees NOTHING
	entriesB := timesheetService.List(userB.ID, 0)
	if len(entriesB) != 0 {
		t.Errorf("User B should see 0 time entries, saw %d", len(entriesB))
	}

	// Verify User B cannot Get it
	_, err = timesheetService.Get(userB.ID, createdTimeEntry.ID)
	if err == nil {
		t.Error("User B should not be able to get User A's time entry")
	}

	// 5. User A creates an Invoice
	invoiceInput := dto.CreateInvoiceInput{
		ClientID:  createdClient.ID,
		Number:    "INV-001",
		IssueDate: "2023-01-01",
		Total:     100,
	}
	createdInvoice := invoiceService.Create(userA.ID, invoiceInput)

	// Verify User A info
	invoicesA := invoiceService.List(userA.ID)
	if len(invoicesA) != 1 {
		t.Errorf("User A should see 1 invoice, saw %d", len(invoicesA))
	}

	// Verify User B sees NOTHING
	invoicesB := invoiceService.List(userB.ID)
	if len(invoicesB) != 0 {
		t.Errorf("User B should see 0 invoices, saw %d", len(invoicesB))
	}

	// Verify User B cannot Get it
	_, err = invoiceService.Get(userB.ID, createdInvoice.ID)
	if err == nil {
		t.Error("User B should not be able to get User A's invoice")
	}

	// 6. Cross-User Update Attempt (User B tries to update User A's client)
	updateInput := dto.UpdateClientInput{
		ID:   createdClient.ID,
		Name: "Hacked Client Name",
	}
	updatedClient := clientService.Update(userB.ID, updateInput)

	// Should return empty/zero value (ClientService.Update returns empty DTO on failure, checks via Get)
	if updatedClient.ID != 0 {
		t.Error("User B should not be able to update User A's client")
	}

	// Verify name didn't change
	refreshedClient, _ := clientService.Get(userA.ID, createdClient.ID)
	if refreshedClient.Name == "Hacked Client Name" {
		t.Error("User A's client was modified by User B!")
	}
}
