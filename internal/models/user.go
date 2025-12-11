// Package models defines database-backed domain models.
package models

import "time"

// User represents an application account owner.
type User struct {
	ID           int       `json:"id"`
	UUID         string    `json:"uuid"` // For future cloud sync
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"` // Never expose hash to frontend
	Email        string    `json:"email"`
	AvatarURL    string    `json:"avatarUrl"`
	CreatedAt    time.Time `json:"createdAt"`
	LastLogin    time.Time `json:"lastLogin"`
	SettingsJSON string    `json:"settingsJson"`
}

// UserSettings defines the structure for the SettingsJSON field
type UserSettings struct {
	Currency       string  `json:"currency"`
	DefaultTaxRate float64 `json:"defaultTaxRate"`
	Language       string  `json:"language"`
	Theme          string  `json:"theme"`
}
