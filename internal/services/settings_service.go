package services

import (
	"database/sql"
	"encoding/json"
	"freelance-flow/internal/dto"
	"freelance-flow/internal/mapper"
	"freelance-flow/internal/models"
	"log"
	"strings"
)

// SettingsService manages user-level preferences stored in settings_json.
type SettingsService struct {
	db *sql.DB
}

// NewSettingsService creates a new SettingsService instance.
func NewSettingsService(db *sql.DB) *SettingsService {
	return &SettingsService{db: db}
}

// Get returns normalized user settings.
func (s *SettingsService) Get(userID int) dto.UserSettings {
	raw := "{}"
	err := s.db.QueryRow("SELECT settings_json FROM users WHERE id = ?", userID).Scan(&raw)
	if err != nil {
		log.Println("Error fetching user settings:", err)
		return mapper.ToUserSettingsDTO(defaultUserSettings())
	}
	settings := defaultUserSettings()
	if err := json.Unmarshal([]byte(raw), &settings); err != nil {
		log.Println("Error parsing settings_json, using defaults:", err)
	}
	normalized := normalizeUserSettings(settings)
	return mapper.ToUserSettingsDTO(normalized)
}

// Update persists settings_json and returns normalized result.
func (s *SettingsService) Update(userID int, input dto.UserSettings) (dto.UserSettings, error) {
	modelSettings := mapper.ToUserSettingsModel(input)
	normalized := normalizeUserSettings(modelSettings)
	payload, err := json.Marshal(normalized)
	if err != nil {
		log.Println("Error marshaling settings:", err)
		return mapper.ToUserSettingsDTO(normalized), err
	}

	_, err = s.db.Exec("UPDATE users SET settings_json = ? WHERE id = ?", string(payload), userID)
	if err != nil {
		log.Println("Error updating settings_json:", err)
		return mapper.ToUserSettingsDTO(normalized), err
	}

	return mapper.ToUserSettingsDTO(normalized), nil
}

// normalizeUserSettings fills defaults and trims spaces.
func normalizeUserSettings(settings models.UserSettings) models.UserSettings {
	trim := func(value string) string {
		return strings.TrimSpace(value)
	}

	if settings.Currency == "" {
		settings.Currency = "USD"
	}
	if settings.DefaultTaxRate < 0 {
		settings.DefaultTaxRate = 0
	}
	if settings.DateFormat == "" {
		settings.DateFormat = "2006-01-02"
	}
	if settings.Timezone == "" {
		settings.Timezone = "UTC"
	}
	if settings.Language == "" {
		settings.Language = "en-US"
	}
	if settings.Theme == "" {
		settings.Theme = "light"
	}
	if settings.InvoiceTerms == "" {
		settings.InvoiceTerms = "Due upon receipt"
	}

	settings.SenderName = trim(settings.SenderName)
	settings.SenderCompany = trim(settings.SenderCompany)
	settings.SenderAddress = trim(settings.SenderAddress)
	settings.SenderPhone = trim(settings.SenderPhone)
	settings.SenderEmail = trim(settings.SenderEmail)
	settings.SenderPostalCode = trim(settings.SenderPostalCode)
	settings.DefaultMessageTemplate = trim(settings.DefaultMessageTemplate)

	return settings
}

// defaultUserSettings returns baseline defaults.
func defaultUserSettings() models.UserSettings {
	return models.UserSettings{
		Currency:               "USD",
		DefaultTaxRate:         0,
		Language:               "en-US",
		Theme:                  "light",
		DateFormat:             "2006-01-02",
		Timezone:               "UTC",
		InvoiceTerms:           "Due upon receipt",
		DefaultMessageTemplate: "",
	}
}
