package services

import (
	"database/sql"
	"freelance-flow/internal/dto"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// setupTestDB creates an in-memory SQLite database for testing.
func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Create users table
	_, err = db.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			uuid TEXT UNIQUE NOT NULL,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			email TEXT,
			avatar_url TEXT,
			created_at TEXT DEFAULT (datetime('now')),
			last_login TEXT,
			settings_json TEXT DEFAULT '{}'
		)
	`)
	if err != nil {
		t.Fatalf("Failed to create users table: %v", err)
	}

	return db
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

// =============================================================================
// Password Hashing Tests
// =============================================================================

func TestHashPassword(t *testing.T) {
	password := "testPassword123!"
	hash, err := HashPassword(password)

	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	if hash == "" {
		t.Error("HashPassword returned empty hash")
	}

	if hash == password {
		t.Error("HashPassword returned plain text password")
	}

	// Verify the hash is different each time (bcrypt uses random salt)
	hash2, _ := HashPassword(password)
	if hash == hash2 {
		t.Error("HashPassword should produce different hashes for same password (random salt)")
	}
}

func TestVerifyPassword_Success(t *testing.T) {
	password := "testPassword123!"
	hash, _ := HashPassword(password)

	if !VerifyPassword(password, hash) {
		t.Error("VerifyPassword should return true for matching password")
	}
}

func TestVerifyPassword_WrongPassword(t *testing.T) {
	password := "testPassword123!"
	wrongPassword := "wrongPassword!"
	hash, _ := HashPassword(password)

	if VerifyPassword(wrongPassword, hash) {
		t.Error("VerifyPassword should return false for wrong password")
	}
}

// =============================================================================
// Register Tests
// =============================================================================

func TestRegister_Success(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close db: %v", err)
		}
	}()

	authService := NewAuthService(db)

	input := dto.RegisterInput{
		Username:  "testuser",
		Password:  "password123",
		Email:     "test@example.com",
		AvatarURL: "https://example.com/avatar.png",
	}

	output, err := authService.Register(input)

	if err != nil {
		t.Fatalf("Register returned error: %v", err)
	}

	if output.ID == 0 {
		t.Error("Register should return a valid ID")
	}

	if output.Username != input.Username {
		t.Errorf("Username mismatch: expected %s, got %s", input.Username, output.Username)
	}

	if output.Email != input.Email {
		t.Errorf("Email mismatch: expected %s, got %s", input.Email, output.Email)
	}

	if output.UUID == "" {
		t.Error("UUID should not be empty")
	}
}

func TestRegister_DuplicateUsername(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close db: %v", err)
		}
	}()

	authService := NewAuthService(db)

	input := dto.RegisterInput{
		Username: "testuser",
		Password: "password123",
	}

	// First registration should succeed
	_, err := authService.Register(input)
	if err != nil {
		t.Fatalf("First Register returned error: %v", err)
	}

	// Second registration with same username should fail
	_, err = authService.Register(input)
	if err != ErrUsernameExists {
		t.Errorf("Expected ErrUsernameExists, got: %v", err)
	}
}

func TestRegister_EmptyUsername(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close db: %v", err)
		}
	}()

	authService := NewAuthService(db)

	input := dto.RegisterInput{
		Username: "",
		Password: "password123",
	}

	_, err := authService.Register(input)
	if err == nil {
		t.Error("Register should fail with empty username")
	}
}

func TestRegister_EmptyPassword(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close db: %v", err)
		}
	}()

	authService := NewAuthService(db)

	input := dto.RegisterInput{
		Username: "testuser",
		Password: "",
	}

	_, err := authService.Register(input)
	if err == nil {
		t.Error("Register should fail with empty password")
	}
}

// =============================================================================
// Login Tests
// =============================================================================

func TestLogin_Success(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close db: %v", err)
		}
	}()

	authService := NewAuthService(db)

	// Register a user first
	registerInput := dto.RegisterInput{
		Username:  "testuser",
		Password:  "password123",
		Email:     "test@example.com",
		AvatarURL: "https://example.com/avatar.png",
	}
	_, err := authService.Register(registerInput)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	// Now login
	loginInput := dto.LoginInput{
		Username: "testuser",
		Password: "password123",
	}

	output, err := authService.Login(loginInput)

	if err != nil {
		t.Fatalf("Login returned error: %v", err)
	}

	if output.Username != registerInput.Username {
		t.Errorf("Username mismatch: expected %s, got %s", registerInput.Username, output.Username)
	}

	if output.ID == 0 {
		t.Error("Login should return a valid user ID")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close db: %v", err)
		}
	}()

	authService := NewAuthService(db)

	// Register a user first
	registerInput := dto.RegisterInput{
		Username: "testuser",
		Password: "password123",
	}
	_, err := authService.Register(registerInput)
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	// Login with wrong password
	loginInput := dto.LoginInput{
		Username: "testuser",
		Password: "wrongpassword",
	}

	_, err = authService.Login(loginInput)

	if err != ErrInvalidCredentials {
		t.Errorf("Expected ErrInvalidCredentials, got: %v", err)
	}
}

func TestLogin_UserNotFound(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close db: %v", err)
		}
	}()

	authService := NewAuthService(db)

	loginInput := dto.LoginInput{
		Username: "nonexistentuser",
		Password: "password123",
	}

	_, err := authService.Login(loginInput)

	if err != ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got: %v", err)
	}
}

// =============================================================================
// GetAllUsers Tests
// =============================================================================

func TestGetAllUsers_Empty(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close db: %v", err)
		}
	}()

	authService := NewAuthService(db)

	users := authService.GetAllUsers()

	if len(users) != 0 {
		t.Errorf("Expected 0 users, got %d", len(users))
	}
}

func TestGetAllUsers_MultipleUsers(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close db: %v", err)
		}
	}()

	authService := NewAuthService(db)

	// Register multiple users
	_, _ = authService.Register(dto.RegisterInput{Username: "user1", Password: "pass1"})
	_, _ = authService.Register(dto.RegisterInput{Username: "user2", Password: "pass2"})
	_, _ = authService.Register(dto.RegisterInput{Username: "user3", Password: "pass3"})

	users := authService.GetAllUsers()

	if len(users) != 3 {
		t.Errorf("Expected 3 users, got %d", len(users))
	}
}

// =============================================================================
// HasUsers Tests
// =============================================================================

func TestHasUsers_False(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close db: %v", err)
		}
	}()

	authService := NewAuthService(db)

	if authService.HasUsers() {
		t.Error("HasUsers should return false when no users exist")
	}
}

func TestHasUsers_True(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close db: %v", err)
		}
	}()

	authService := NewAuthService(db)

	// Register a user
	_, _ = authService.Register(dto.RegisterInput{Username: "testuser", Password: "pass"})

	if !authService.HasUsers() {
		t.Error("HasUsers should return true when users exist")
	}
}

// =============================================================================
// GetUserByID Tests
// =============================================================================

func TestGetUserByID_Success(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close db: %v", err)
		}
	}()

	authService := NewAuthService(db)

	// Register a user
	registered, _ := authService.Register(dto.RegisterInput{
		Username: "testuser",
		Password: "pass",
		Email:    "test@example.com",
	})

	// Get user by ID
	user, err := authService.GetUserByID(registered.ID)

	if err != nil {
		t.Fatalf("GetUserByID returned error: %v", err)
	}

	if user.Username != "testuser" {
		t.Errorf("Username mismatch: expected testuser, got %s", user.Username)
	}
}

func TestGetUserByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close db: %v", err)
		}
	}()

	authService := NewAuthService(db)

	_, err := authService.GetUserByID(999)

	if err != ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got: %v", err)
	}
}
