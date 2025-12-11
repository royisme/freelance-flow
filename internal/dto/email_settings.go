package dto

// InvoiceEmailSettings represents outbound email configuration for invoices.
type InvoiceEmailSettings struct {
	Provider        string `json:"provider"` // mailto | resend | smtp
	FromEmail       string `json:"from"`
	ReplyTo         string `json:"replyTo"`
	SubjectTemplate string `json:"subjectTemplate"`
	BodyTemplate    string `json:"bodyTemplate"`
	Signature       string `json:"signature"`
	ResendAPIKey    string `json:"resendApiKey"`
	SMTPHost        string `json:"smtpHost"`
	SMTPPort        int    `json:"smtpPort"`
	SMTPUsername    string `json:"smtpUsername"`
	SMTPPassword    string `json:"smtpPassword"`
	SMTPUseTLS      bool   `json:"smtpUseTls"`
}
