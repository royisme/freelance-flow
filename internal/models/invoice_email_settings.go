package models

import "time"

// InvoiceEmailSettings holds outbound email configuration for invoices.
type InvoiceEmailSettings struct {
	UserID          int
	Provider        string
	FromEmail       string
	ReplyTo         string
	SubjectTemplate string
	BodyTemplate    string
	Signature       string
	ResendAPIKey    string
	SMTPHost        string
	SMTPPort        int
	SMTPUsername    string
	SMTPPassword    string
	SMTPUseTLS      bool
	UpdatedAt       time.Time
}
