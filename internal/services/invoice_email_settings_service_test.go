package services

import (
	"database/sql"
	"freelance-flow/internal/dto"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestInvoiceEmailSettingsService_GetAndUpdate(t *testing.T) {
	db, err := sql.Open("sqlite3", ":memory:")
	assert.NoError(t, err)
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE users (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT);`)
	assert.NoError(t, err)
	_, err = db.Exec(`CREATE TABLE invoice_email_settings (
		user_id INTEGER PRIMARY KEY,
		provider TEXT,
		from_email TEXT,
		reply_to TEXT,
		subject_template TEXT,
		body_template TEXT,
		signature TEXT,
		resend_api_key TEXT,
		smtp_host TEXT,
		smtp_port INTEGER,
		smtp_username TEXT,
		smtp_password TEXT,
		smtp_use_tls INTEGER,
		updated_at TEXT
	);`)
	assert.NoError(t, err)
	_, err = db.Exec(`INSERT INTO users(id, username) VALUES (1, 'user')`)
	assert.NoError(t, err)

	svc := NewInvoiceEmailSettingsService(db)

	// default
	def := svc.Get(1)
	assert.Equal(t, "mailto", def.Provider)

	// update
	updated, err := svc.Update(1, dto.InvoiceEmailSettings{
		Provider:        "smtp",
		FromEmail:       "from@example.com",
		ReplyTo:         "reply@example.com",
		SubjectTemplate: "Invoice {{number}}",
		BodyTemplate:    "Hello",
		Signature:       "Sig",
		SMTPHost:        "smtp.example.com",
		SMTPPort:        465,
		SMTPUsername:    "u",
		SMTPPassword:    "p",
		SMTPUseTLS:      true,
	})
	assert.NoError(t, err)
	assert.Equal(t, "smtp", updated.Provider)
	assert.Equal(t, "smtp.example.com", updated.SMTPHost)

	fetched := svc.Get(1)
	assert.Equal(t, "smtp", fetched.Provider)
	assert.Equal(t, 465, fetched.SMTPPort)
	assert.Equal(t, "from@example.com", fetched.FromEmail)
}
