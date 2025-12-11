package services

import (
	"database/sql"
	"freelance-flow/internal/dto"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

// TestMultiUserFlow_Integration simulates a complete user lifecycle to verify system integrity.
// Flow:
// 1. Initial State: No users/data.
// 2. Register User A.
// 3. User A creates a client and project.
// 4. Register User B.
// 5. Verify User B sees NO clients/projects (Isolation).
// 6. User B creates their own client.
// 7. Switch back to User A (Login).
// 8. Verify User A strictly sees their own client/project, not User B's.
func TestMultiUserFlow_Integration(t *testing.T) {
	// 1. Setup in-memory DB
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	defer func() {
		if cerr := db.Close(); cerr != nil {
			t.Errorf("failed to close db: %v", cerr)
		}
	}()

	// Create tables with full schema
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
		_, err = db.Exec(query)
		assert.NoError(t, err)
	}

	// Initialize Services
	authService := NewAuthService(db)
	clientService := NewClientService(db)
	projectService := NewProjectService(db)

	// --- Step 2: Register User A ---
	userAInput := dto.RegisterInput{
		Username: "UserA",
		Password: "password123",
		Email:    "a@example.com",
	}
	userA, err := authService.Register(userAInput)
	assert.NoError(t, err)
	assert.NotZero(t, userA.ID)
	t.Logf("Registered User A: ID=%d", userA.ID)

	// --- Step 3: User A creates data ---
	// Create Client for A
	clientAInput := dto.CreateClientInput{
		Name:  "UserA Client",
		Email: "clienta@test.com",
	}
	// Create returns only ClientOutput (errors logged)
	clientA := clientService.Create(userA.ID, clientAInput)
	assert.NotZero(t, clientA.ID)
	assert.Equal(t, "UserA Client", clientA.Name)

	// Create Project for A
	projectAInput := dto.CreateProjectInput{
		Name:     "UserA Project",
		ClientID: clientA.ID,
		Status:   "Active",
	}
	projectA := projectService.Create(userA.ID, projectAInput)
	assert.NotZero(t, projectA.ID)
	assert.Equal(t, "UserA Project", projectA.Name)

	// --- Step 4: Register User B ---
	userBInput := dto.RegisterInput{
		Username: "UserB",
		Password: "password456",
		Email:    "b@example.com",
	}
	userB, err := authService.Register(userBInput)
	assert.NoError(t, err)
	assert.NotZero(t, userB.ID)
	assert.NotEqual(t, userA.ID, userB.ID)
	t.Logf("Registered User B: ID=%d", userB.ID)

	// --- Step 5: Verify Isolation (User B reads) ---
	// B tries to list clients
	clientsB := clientService.List(userB.ID)
	assert.Equal(t, 0, len(clientsB), "User B should see 0 clients initially")

	// B tries to list projects
	projectsB := projectService.List(userB.ID)
	assert.Equal(t, 0, len(projectsB), "User B should see 0 projects initially")

	// B tries to Get User A's client directly
	_, err = clientService.Get(userB.ID, clientA.ID)
	// Expect error because row is filtered by user_id
	assert.Error(t, err, "User B should not be able to get User A's client")

	// --- Step 6: User B creates their own data ---
	clientBInput := dto.CreateClientInput{
		Name: "UserB Client",
	}
	clientB := clientService.Create(userB.ID, clientBInput)
	assert.NotZero(t, clientB.ID)

	// Verify B sees their client
	clientsB = clientService.List(userB.ID)
	assert.Equal(t, 1, len(clientsB))
	assert.Equal(t, "UserB Client", clientsB[0].Name)

	// --- Step 7: Switch back to User A (Simulate Login) ---
	loginA, err := authService.Login(dto.LoginInput{Username: "UserA", Password: "password123"})
	assert.NoError(t, err)
	assert.Equal(t, userA.ID, loginA.ID)

	// --- Step 8: Verify User A sees only their data ---
	clientsA := clientService.List(userA.ID)
	assert.Equal(t, 1, len(clientsA))
	assert.Equal(t, "UserA Client", clientsA[0].Name)
	assert.NotEqual(t, clientB.ID, clientsA[0].ID)

	t.Log("Integration Flow Verified Successfully")
}

func TestAuth_LoginFlow_Integration(t *testing.T) {
	// Setup DB
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	defer func() {
		if cerr := db.Close(); cerr != nil {
			t.Errorf("failed to close db: %v", cerr)
		}
	}()

	// Create users table
	_, err = db.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			uuid TEXT NOT NULL UNIQUE,
			username TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			email TEXT,
			avatar_url TEXT,
			settings_json TEXT DEFAULT '{}',
			created_at DATETIME,
			last_login DATETIME
		);
	`)
	assert.NoError(t, err)

	authService := NewAuthService(db)

	// 1. Register
	regInput := dto.RegisterInput{
		Username:     "TestLogin",
		Password:     "CorrectHorseBatteryStaple",
		SettingsJSON: `{"theme":"dark"}`,
	}
	user, err := authService.Register(regInput)
	assert.NoError(t, err)
	assert.Equal(t, `{"theme":"dark"}`, user.SettingsJSON)

	// 2. Login Success
	loginInput := dto.LoginInput{
		Username: "TestLogin",
		Password: "CorrectHorseBatteryStaple",
	}
	loggedInUser, err := authService.Login(loginInput)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, loggedInUser.ID)
	assert.Equal(t, user.Username, loggedInUser.Username)

	// Check LastLogin updated
	assert.NotEmpty(t, loggedInUser.LastLogin)
	// Optionally parse time if needed, but NotEmpty is sufficient for basic verification

	// 3. Login Failure (Wrong Password)
	badPassInput := dto.LoginInput{
		Username: "TestLogin",
		Password: "WrongPassword",
	}
	_, err = authService.Login(badPassInput)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid") // Expect "invalid username or password"
}
