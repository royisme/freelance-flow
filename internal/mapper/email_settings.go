package mapper

import (
	"freelance-flow/internal/dto"
	"freelance-flow/internal/models"
	"strings"
	"time"
)

// ToInvoiceEmailSettingsDTO converts model to dto.
func ToInvoiceEmailSettingsDTO(m models.InvoiceEmailSettings) dto.InvoiceEmailSettings {
	return dto.InvoiceEmailSettings{
		Provider:        m.Provider,
		FromEmail:       m.FromEmail,
		ReplyTo:         m.ReplyTo,
		SubjectTemplate: m.SubjectTemplate,
		BodyTemplate:    m.BodyTemplate,
		Signature:       m.Signature,
		ResendAPIKey:    m.ResendAPIKey,
		SMTPHost:        m.SMTPHost,
		SMTPPort:        m.SMTPPort,
		SMTPUsername:    m.SMTPUsername,
		SMTPPassword:    m.SMTPPassword,
		SMTPUseTLS:      m.SMTPUseTLS,
	}
}

// ToInvoiceEmailSettingsModel converts dto to model.
func ToInvoiceEmailSettingsModel(d dto.InvoiceEmailSettings, userID int) models.InvoiceEmailSettings {
	return models.InvoiceEmailSettings{
		UserID:          userID,
		Provider:        strings.TrimSpace(d.Provider),
		FromEmail:       strings.TrimSpace(d.FromEmail),
		ReplyTo:         strings.TrimSpace(d.ReplyTo),
		SubjectTemplate: strings.TrimSpace(d.SubjectTemplate),
		BodyTemplate:    strings.TrimSpace(d.BodyTemplate),
		Signature:       strings.TrimSpace(d.Signature),
		ResendAPIKey:    strings.TrimSpace(d.ResendAPIKey),
		SMTPHost:        strings.TrimSpace(d.SMTPHost),
		SMTPPort:        d.SMTPPort,
		SMTPUsername:    strings.TrimSpace(d.SMTPUsername),
		SMTPPassword:    strings.TrimSpace(d.SMTPPassword),
		SMTPUseTLS:      d.SMTPUseTLS,
		UpdatedAt:       time.Now(),
	}
}
