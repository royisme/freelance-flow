package services

import (
	"database/sql"
	"encoding/json"
	"errors"
	"freelance-flow/internal/dto"
	"freelance-flow/internal/mapper"
	"freelance-flow/internal/models"
	"log"
	"strings"
	"time"
)

// InvoiceEmailSettingsService manages per-user invoice email configuration.
type InvoiceEmailSettingsService struct {
	db *sql.DB
}

// NewInvoiceEmailSettingsService creates a new instance.
func NewInvoiceEmailSettingsService(db *sql.DB) *InvoiceEmailSettingsService {
	return &InvoiceEmailSettingsService{db: db}
}

// Get returns settings for a user, falling back to defaults.
func (s *InvoiceEmailSettingsService) Get(userID int) dto.InvoiceEmailSettings {
	var stored models.InvoiceEmailSettings
	err := s.db.QueryRow(`
SELECT user_id, provider, from_email, reply_to, subject_template, body_template, signature,
       resend_api_key, smtp_host, smtp_port, smtp_username, smtp_password, smtp_use_tls
FROM invoice_email_settings
WHERE user_id = ?`, userID).Scan(
		&stored.UserID,
		&stored.Provider,
		&stored.FromEmail,
		&stored.ReplyTo,
		&stored.SubjectTemplate,
		&stored.BodyTemplate,
		&stored.Signature,
		&stored.ResendAPIKey,
		&stored.SMTPHost,
		&stored.SMTPPort,
		&stored.SMTPUsername,
		&stored.SMTPPassword,
		&stored.SMTPUseTLS,
	)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Println("Error fetching invoice email settings:", err)
		}
		return mapper.ToInvoiceEmailSettingsDTO(defaultInvoiceEmailSettings())
	}
	return mapper.ToInvoiceEmailSettingsDTO(normalizeInvoiceEmailSettings(stored))
}

// Update upserts settings for a user.
func (s *InvoiceEmailSettingsService) Update(userID int, input dto.InvoiceEmailSettings) (dto.InvoiceEmailSettings, error) {
	model := mapper.ToInvoiceEmailSettingsModel(input, userID)
	model = normalizeInvoiceEmailSettings(model)

	// Store sensitive fields as-is; if future encryption needed, apply here.
	_, err := s.db.Exec(`
INSERT INTO invoice_email_settings (
    user_id, provider, from_email, reply_to, subject_template, body_template, signature,
    resend_api_key, smtp_host, smtp_port, smtp_username, smtp_password, smtp_use_tls, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(user_id) DO UPDATE SET
    provider = excluded.provider,
    from_email = excluded.from_email,
    reply_to = excluded.reply_to,
    subject_template = excluded.subject_template,
    body_template = excluded.body_template,
    signature = excluded.signature,
    resend_api_key = excluded.resend_api_key,
    smtp_host = excluded.smtp_host,
    smtp_port = excluded.smtp_port,
    smtp_username = excluded.smtp_username,
    smtp_password = excluded.smtp_password,
    smtp_use_tls = excluded.smtp_use_tls,
    updated_at = excluded.updated_at
`, model.UserID, model.Provider, model.FromEmail, model.ReplyTo, model.SubjectTemplate, model.BodyTemplate, model.Signature,
		model.ResendAPIKey, model.SMTPHost, model.SMTPPort, model.SMTPUsername, model.SMTPPassword, model.SMTPUseTLS, time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		log.Println("Error upserting invoice email settings:", err)
		return mapper.ToInvoiceEmailSettingsDTO(model), err
	}

	return mapper.ToInvoiceEmailSettingsDTO(model), nil
}

func normalizeInvoiceEmailSettings(m models.InvoiceEmailSettings) models.InvoiceEmailSettings {
	trim := func(v string) string { return strings.TrimSpace(v) }
	if m.Provider == "" {
		m.Provider = "mailto"
	}
	m.FromEmail = trim(m.FromEmail)
	m.ReplyTo = trim(m.ReplyTo)
	m.SubjectTemplate = trim(m.SubjectTemplate)
	m.BodyTemplate = trim(m.BodyTemplate)
	m.Signature = trim(m.Signature)
	m.ResendAPIKey = trim(m.ResendAPIKey)
	m.SMTPHost = trim(m.SMTPHost)
	m.SMTPUsername = trim(m.SMTPUsername)
	m.SMTPPassword = trim(m.SMTPPassword)
	if m.SMTPPort == 0 {
		m.SMTPPort = 587
	}
	return m
}

func defaultInvoiceEmailSettings() models.InvoiceEmailSettings {
	return models.InvoiceEmailSettings{
		Provider:        "mailto",
		SubjectTemplate: "Invoice {{number}}",
		BodyTemplate:    "Please find attached invoice {{number}}.",
		Signature:       "",
		SMTPUseTLS:      true,
		SMTPPort:        587,
	}
}

// ExportSettings returns settings as JSON (for future backup/diagnostics).
func (s *InvoiceEmailSettingsService) ExportSettings(userID int) string {
	settings := s.Get(userID)
	b, err := json.Marshal(settings)
	if err != nil {
		return "{}"
	}
	return string(b)
}
