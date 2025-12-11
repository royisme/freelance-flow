package services

import (
	"database/sql"
	"errors"
	"freelance-flow/internal/dto"
	"freelance-flow/internal/mapper"
	"freelance-flow/internal/models"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Common errors for AuthService
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUsernameExists     = errors.New("username already exists")
)

// AuthService handles user authentication and management.
type AuthService struct {
	db *sql.DB
}

// NewAuthService creates a new AuthService instance.
func NewAuthService(db *sql.DB) *AuthService {
	return &AuthService{db: db}
}

// HashPassword hashes a password using bcrypt.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// VerifyPassword checks if the provided password matches the hash.
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Register creates a new user account.
func (s *AuthService) Register(input dto.RegisterInput) (dto.UserOutput, error) {
	// Validate input
	if input.Username == "" {
		return dto.UserOutput{}, errors.New("username is required")
	}
	if input.Password == "" {
		return dto.UserOutput{}, errors.New("password is required")
	}

	// Check if username already exists
	var existingID int
	err := s.db.QueryRow("SELECT id FROM users WHERE username = ?", input.Username).Scan(&existingID)
	if err == nil {
		return dto.UserOutput{}, ErrUsernameExists
	}
	if err != sql.ErrNoRows {
		log.Println("Error checking username:", err)
		return dto.UserOutput{}, err
	}

	// Hash the password
	passwordHash, err := HashPassword(input.Password)
	if err != nil {
		log.Println("Error hashing password:", err)
		return dto.UserOutput{}, err
	}

	// Generate UUID for future cloud sync
	userUUID := uuid.New().String()

	// Default settings if empty
	settingsJSON := input.SettingsJSON
	if settingsJSON == "" {
		settingsJSON = "{}"
	}

	// Insert user
	stmt, err := s.db.Prepare(`
		INSERT INTO users(uuid, username, password_hash, email, avatar_url, settings_json, created_at, last_login)
		VALUES(?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		log.Println("Error preparing insert:", err)
		return dto.UserOutput{}, err
	}
	defer closeWithLog(stmt, "closing register statement")

	now := time.Now()
	res, err := stmt.Exec(userUUID, input.Username, passwordHash, input.Email, input.AvatarURL, settingsJSON, now.Format(time.RFC3339), now.Format(time.RFC3339))
	if err != nil {
		log.Println("Error inserting user:", err)
		return dto.UserOutput{}, err
	}

	id, _ := res.LastInsertId()
	user := models.User{
		ID:           int(id),
		UUID:         userUUID,
		Username:     input.Username,
		Email:        input.Email,
		AvatarURL:    input.AvatarURL,
		CreatedAt:    now,
		LastLogin:    now,
		SettingsJSON: settingsJSON,
	}

	return mapper.ToUserOutput(user), nil
}

// Login authenticates a user and returns user info.
func (s *AuthService) Login(input dto.LoginInput) (dto.UserOutput, error) {
	if input.Username == "" || input.Password == "" {
		return dto.UserOutput{}, ErrInvalidCredentials
	}

	var user models.User
	var createdAtStr, lastLoginStr sql.NullString
	err := s.db.QueryRow(`
		SELECT id, uuid, username, password_hash, email, avatar_url, created_at, last_login, settings_json
		FROM users WHERE username = ?
	`, input.Username).Scan(
		&user.ID, &user.UUID, &user.Username, &user.PasswordHash,
		&user.Email, &user.AvatarURL, &createdAtStr, &lastLoginStr, &user.SettingsJSON,
	)

	if err == sql.ErrNoRows {
		return dto.UserOutput{}, ErrUserNotFound
	}
	if err != nil {
		log.Println("Error querying user:", err)
		return dto.UserOutput{}, err
	}

	// Verify password
	if !VerifyPassword(input.Password, user.PasswordHash) {
		return dto.UserOutput{}, ErrInvalidCredentials
	}

	// Parse timestamps
	if createdAtStr.Valid {
		user.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr.String)
	}

	// Update last login
	now := time.Now()
	user.LastLogin = now
	_, err = s.db.Exec("UPDATE users SET last_login = ? WHERE id = ?", now.Format(time.RFC3339), user.ID)
	if err != nil {
		log.Println("Error updating last login:", err)
	}

	return mapper.ToUserOutput(user), nil
}

// GetAllUsers returns a list of all users (for user selection screen).
func (s *AuthService) GetAllUsers() []dto.UserListItem {
	rows, err := s.db.Query("SELECT id, username, avatar_url FROM users ORDER BY last_login DESC")
	if err != nil {
		log.Println("Error querying users:", err)
		return []dto.UserListItem{}
	}
	defer closeWithLog(rows, "closing users rows")

	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.ID, &u.Username, &u.AvatarURL)
		if err != nil {
			log.Println("Error scanning user:", err)
			continue
		}
		users = append(users, u)
	}
	return mapper.ToUserListItemList(users)
}

// GetUserByID returns a single user by ID.
func (s *AuthService) GetUserByID(id int) (dto.UserOutput, error) {
	var user models.User
	var createdAtStr, lastLoginStr sql.NullString
	err := s.db.QueryRow(`
		SELECT id, uuid, username, email, avatar_url, created_at, last_login, settings_json
		FROM users WHERE id = ?
	`, id).Scan(
		&user.ID, &user.UUID, &user.Username, &user.Email,
		&user.AvatarURL, &createdAtStr, &lastLoginStr, &user.SettingsJSON,
	)

	if err == sql.ErrNoRows {
		return dto.UserOutput{}, ErrUserNotFound
	}
	if err != nil {
		return dto.UserOutput{}, err
	}

	// Parse timestamps
	if createdAtStr.Valid {
		user.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr.String)
	}
	if lastLoginStr.Valid {
		user.LastLogin, _ = time.Parse(time.RFC3339, lastLoginStr.String)
	}

	return mapper.ToUserOutput(user), nil
}

// closeWithLog logs an error if closing the resource fails.
func closeWithLog(closer interface{ Close() error }, context string) {
	if err := closer.Close(); err != nil {
		log.Printf("Error %s: %v", context, err)
	}
}

// HasUsers checks if any users exist in the database.
func (s *AuthService) HasUsers() bool {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if err != nil {
		log.Println("Error counting users:", err)
		return false
	}
	return count > 0
}
