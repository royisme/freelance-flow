package services

import (
	"bytes"
	"crypto/tls"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"
	"tally/internal/dto"
	"tally/internal/mapper"
	"tally/internal/models"
	"tally/internal/pdf"
	"tally/internal/utils"

	"github.com/resend/resend-go/v3"
)

// InvoiceService handles all invoice-related operations.
type InvoiceService struct {
	db *sql.DB
}

// NewInvoiceService creates a new InvoiceService instance.
func NewInvoiceService(db *sql.DB) *InvoiceService {
	return &InvoiceService{db: db}
}

// List returns all invoices for a specific user as DTOs.
func (s *InvoiceService) List(userID int) []dto.InvoiceOutput {
	rows, err := s.db.Query(`
		SELECT i.id, i.client_id, 
		(SELECT project_id FROM time_entries WHERE invoice_id = i.id LIMIT 1) as project_id,
		i.number, i.issue_date, i.due_date, i.subtotal, i.tax_rate, i.tax_amount, i.total, i.status, i.items_json 
		FROM invoices i WHERE i.user_id = ?`, userID)
	if err != nil {
		log.Println("Error querying invoices:", err)
		return []dto.InvoiceOutput{}
	}
	defer closeWithLog(rows, "closing invoice rows")

	var invoices []dto.InvoiceOutput // Direct to DTO since we changed the struct
	// Wait, scan needs to map to variables.
	// To avoid changing `models.Invoice` (which maps to DB table directly), let's scan into variables directly here.

	for rows.Next() {
		var id, clientId int
		var projectId sql.NullInt64
		var number, issueDate, dueDate, status, itemsJSON string
		var subtotal, taxRate, taxAmount, total float64

		err := rows.Scan(&id, &clientId, &projectId, &number, &issueDate, &dueDate, &subtotal, &taxRate, &taxAmount, &total, &status, &itemsJSON)
		if err != nil {
			log.Println("Error scanning invoice:", err)
			continue
		}

		var items []dto.InvoiceItemOutput
		if itemsJSON != "" {
			var entityItems []models.InvoiceItem
			if err := json.Unmarshal([]byte(itemsJSON), &entityItems); err == nil {
				items = mapper.ToInvoiceItemOutputList(entityItems)
			}
		} else {
			items = []dto.InvoiceItemOutput{}
		}

		invoices = append(invoices, dto.InvoiceOutput{
			ID:        id,
			ClientID:  clientId,
			ProjectID: int(projectId.Int64),
			Number:    number,
			IssueDate: issueDate,
			DueDate:   dueDate,
			Subtotal:  subtotal,
			TaxRate:   taxRate,
			TaxAmount: taxAmount,
			Total:     total,
			Status:    status,
			Items:     items,
		})
	}
	return invoices
}

// Get returns a single invoice by ID for a specific user.
func (s *InvoiceService) Get(userID int, id int) (dto.InvoiceOutput, error) {
	row := s.db.QueryRow(`
		SELECT i.id, i.client_id, 
		(SELECT project_id FROM time_entries WHERE invoice_id = i.id LIMIT 1) as project_id,
		i.number, i.issue_date, i.due_date, i.subtotal, i.tax_rate, i.tax_amount, i.total, i.status, i.items_json 
		FROM invoices i WHERE i.id = ? AND i.user_id = ?`, id, userID)

	var invId, clientId int
	var projectId sql.NullInt64
	var number, issueDate, dueDate, status, itemsJSON string
	var subtotal, taxRate, taxAmount, total float64

	err := row.Scan(&invId, &clientId, &projectId, &number, &issueDate, &dueDate, &subtotal, &taxRate, &taxAmount, &total, &status, &itemsJSON)
	if err != nil {
		return dto.InvoiceOutput{}, err
	}

	var items []dto.InvoiceItemOutput
	if itemsJSON != "" {
		var entityItems []models.InvoiceItem
		if err := json.Unmarshal([]byte(itemsJSON), &entityItems); err != nil {
			log.Printf("Error unmarshalling items for invoice %d: %v", invId, err)
			return dto.InvoiceOutput{}, fmt.Errorf("failed to unmarshal items: %w", err)
		}
		items = mapper.ToInvoiceItemOutputList(entityItems)
	} else {
		items = []dto.InvoiceItemOutput{}
	}

	return dto.InvoiceOutput{
		ID:        invId,
		ClientID:  clientId,
		ProjectID: int(projectId.Int64),
		Number:    number,
		IssueDate: issueDate,
		DueDate:   dueDate,
		Subtotal:  subtotal,
		TaxRate:   taxRate,
		TaxAmount: taxAmount,
		Total:     total,
		Status:    status,
		Items:     items,
	}, nil
}

// Create adds a new invoice for a specific user and returns the created invoice as DTO.
func (s *InvoiceService) Create(userID int, input dto.CreateInvoiceInput) dto.InvoiceOutput {
	entity := mapper.ToInvoiceEntity(input)
	itemsBytes, _ := json.Marshal(entity.Items)
	itemsJSON := string(itemsBytes)

	stmt, err := s.db.Prepare("INSERT INTO invoices(user_id, client_id, number, issue_date, due_date, subtotal, tax_rate, tax_amount, total, status, items_json) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("Error preparing invoice insert:", err)
		return dto.InvoiceOutput{}
	}
	defer closeWithLog(stmt, "closing invoice insert statement")

	res, err := stmt.Exec(userID, entity.ClientID, entity.Number, entity.IssueDate, entity.DueDate, entity.Subtotal, entity.TaxRate, entity.TaxAmount, entity.Total, entity.Status, itemsJSON)
	if err != nil {
		log.Println("Error inserting invoice:", err)
		return dto.InvoiceOutput{}
	}

	id, _ := res.LastInsertId()
	entity.ID = int(id)
	return mapper.ToInvoiceOutput(entity)
}

// Update modifies an existing invoice for a specific user and returns the updated invoice as DTO.
func (s *InvoiceService) Update(userID int, input dto.UpdateInvoiceInput) dto.InvoiceOutput {
	// Convert items to JSON
	items := mapper.ToInvoiceItemEntityList(input.Items)
	itemsBytes, _ := json.Marshal(items)
	itemsJSON := string(itemsBytes)

	stmt, err := s.db.Prepare("UPDATE invoices SET client_id=?, number=?, issue_date=?, due_date=?, subtotal=?, tax_rate=?, tax_amount=?, total=?, status=?, items_json=? WHERE id=? AND user_id=?")
	if err != nil {
		log.Println("Error preparing invoice update:", err)
		return dto.InvoiceOutput{}
	}
	defer closeWithLog(stmt, "closing invoice update statement")

	_, err = stmt.Exec(input.ClientID, input.Number, input.IssueDate, input.DueDate, input.Subtotal, input.TaxRate, input.TaxAmount, input.Total, input.Status, itemsJSON, input.ID, userID)
	if err != nil {
		log.Println("Error updating invoice:", err)
		return dto.InvoiceOutput{}
	}

	output, _ := s.Get(userID, input.ID)
	return output
}

// UpdateStatus updates the status of an invoice.
func (s *InvoiceService) UpdateStatus(userID int, invoiceID int, status string) error {
	stmt, err := s.db.Prepare("UPDATE invoices SET status=? WHERE id=? AND user_id=?")
	if err != nil {
		log.Println("Error preparing invoice status update:", err)
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer closeWithLog(stmt, "closing invoice status update statement")

	res, err := stmt.Exec(status, invoiceID, userID)
	if err != nil {
		log.Println("Error updating invoice status:", err)
		return fmt.Errorf("failed to update status: %w", err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("invoice not found or not owned by user")
	}

	return nil
}

// Delete removes an invoice by ID for a specific user and unlinks associated time entries.
func (s *InvoiceService) Delete(userID int, id int) {
	if _, err := s.db.Exec("UPDATE time_entries SET invoice_id = NULL, invoiced = 0 WHERE user_id = ? AND invoice_id = ?", userID, id); err != nil {
		log.Println("Error unlinking time entries from invoice:", err)
	}
	if _, err := s.db.Exec("DELETE FROM invoices WHERE id=? AND user_id=?", id, userID); err != nil {
		log.Println("Error deleting invoice:", err)
	}
}

// SendEmail sends invoice via configured provider (mailto/resend/smtp placeholder).
func (s *InvoiceService) SendEmail(userID int, invoiceID int) error {
	invoice, err := s.Get(userID, invoiceID)
	if err != nil {
		log.Println("SendEmail: invoice not found:", err)
		return fmt.Errorf("invoice not found: %w", err)
	}
	client, err := s.getClient(userID, invoice.ClientID)
	if err != nil {
		log.Println("SendEmail: client not found:", err)
		return fmt.Errorf("client not found: %w", err)
	}

	settingsSvc := NewInvoiceEmailSettingsService(s.db)
	emailSettings := settingsSvc.Get(userID)
	if emailSettings.Provider == "" {
		emailSettings.Provider = "mailto"
	}

	pdfBase64, err := s.GeneratePDF(userID, invoiceID, "")
	if err != nil {
		log.Println("SendEmail: generate pdf failed:", err)
		return fmt.Errorf("failed to generate PDF: %w", err)
	}
	pdfBytes, err := base64.StdEncoding.DecodeString(pdfBase64)
	if err != nil {
		log.Println("SendEmail: decode pdf failed:", err)
		return fmt.Errorf("failed to decode PDF: %w", err)
	}

	// mailto just indicate success (front-end opens client)
	if emailSettings.Provider == "mailto" {
		return nil
	}

	// Resend provider
	if emailSettings.Provider == "resend" {
		apiKey := emailSettings.ResendAPIKey
		if apiKey == "" {
			return fmt.Errorf("resend API key is missing")
		}
		subject := utils.ApplyTemplate(emailSettings.SubjectTemplate, invoice)
		if subject == "" {
			subject = fmt.Sprintf("Invoice %s", invoice.Number)
		}
		body := utils.ApplyTemplate(emailSettings.BodyTemplate, invoice)
		if body == "" {
			body = "Please see attached invoice."
		}

		if os.Getenv("RESEND_DRY_RUN") == "1" {
			log.Println("SendEmail: RESEND_DRY_RUN enabled, skipping network call")
			return nil
		}

		clientResend := resend.NewClient(apiKey)
		_, err := clientResend.Emails.Send(&resend.SendEmailRequest{
			From:    emailSettings.FromEmail,
			To:      []string{client.Email},
			Subject: subject,
			Html:    body,
			Attachments: []*resend.Attachment{
				{
					Filename: fmt.Sprintf("INV-%s.pdf", invoice.Number),
					Content:  pdfBytes,
				},
			},
		})
		if err != nil {
			log.Println("SendEmail: resend send failed:", err)
			return fmt.Errorf("resend failed: %v", err)
		}
		// Auto-update status to 'sent' if currently 'draft'
		if invoice.Status == "draft" {
			if err := s.UpdateStatus(userID, invoiceID, "sent"); err != nil {
				log.Println("SendEmail: failed to update status after resend:", err)
				// Don't return error - email was sent successfully
			}
		}
		return nil
	}

	// SMTP provider
	if emailSettings.Provider == "smtp" {
		// Validate SMTP settings
		if emailSettings.SMTPHost == "" {
			return fmt.Errorf("smtp host is missing")
		}
		if emailSettings.SMTPUsername == "" {
			return fmt.Errorf("smtp username is missing")
		}
		if emailSettings.SMTPPassword == "" {
			return fmt.Errorf("smtp password is missing")
		}
		if emailSettings.FromEmail == "" {
			return fmt.Errorf("smtp from email is missing")
		}

		subject := utils.ApplyTemplate(emailSettings.SubjectTemplate, invoice)
		if subject == "" {
			subject = fmt.Sprintf("Invoice %s", invoice.Number)
		}
		body := utils.ApplyTemplate(emailSettings.BodyTemplate, invoice)
		if body == "" {
			body = "Please see attached invoice."
		}

		// Add signature if provided
		if emailSettings.Signature != "" {
			body += "\n\n" + emailSettings.Signature
		}

		if os.Getenv("SMTP_DRY_RUN") == "1" {
			log.Println("SendEmail: SMTP_DRY_RUN enabled, skipping network call")
			return nil
		}

		// Send email via SMTP
		if err := s.sendViaSMTP(emailSettings, client.Email, subject, body, pdfBytes, invoice.Number); err != nil {
			log.Println("SendEmail: SMTP send failed:", err)
			return fmt.Errorf("smtp failed: %w", err)
		}
		// Auto-update status to 'sent' if currently 'draft'
		if invoice.Status == "draft" {
			if err := s.UpdateStatus(userID, invoiceID, "sent"); err != nil {
				log.Println("SendEmail: failed to update status after smtp:", err)
				// Don't return error - email was sent successfully
			}
		}
		return nil
	}

	return fmt.Errorf("unknown email provider: %s", emailSettings.Provider)
}

// SetTimeEntries associates time entries with an invoice and recalculates totals.
func (s *InvoiceService) SetTimeEntries(userID int, input dto.SetInvoiceTimeEntriesInput) (dto.InvoiceOutput, error) {
	// Ensure invoice belongs to user
	inv, err := s.Get(userID, input.InvoiceID)
	if err != nil {
		return dto.InvoiceOutput{}, fmt.Errorf("invoice not found: %w", err)
	}

	// Clear existing links
	if _, err := s.db.Exec("UPDATE time_entries SET invoice_id = NULL, invoiced = 0 WHERE user_id = ? AND invoice_id = ?", userID, input.InvoiceID); err != nil {
		return dto.InvoiceOutput{}, fmt.Errorf("failed to clear previous links: %w", err)
	}

	// Apply new links
	if len(input.TimeEntryIDs) > 0 {
		for _, id := range input.TimeEntryIDs {
			if _, err := s.db.Exec(
				"UPDATE time_entries SET invoice_id = ?, invoiced = 1 WHERE user_id = ? AND id = ?",
				input.InvoiceID, userID, id,
			); err != nil {
				return dto.InvoiceOutput{}, fmt.Errorf("failed to link time entry %d: %w", id, err)
			}
		}
	}

	// Recalculate totals from linked entries
	updated, err := s.recalculateInvoiceFromTimeEntries(userID, input.InvoiceID, inv.TaxRate)
	if err != nil {
		return dto.InvoiceOutput{}, err
	}
	return updated, nil
}

// GeneratePDF builds a PDF based on invoice, client, settings, and linked time entries.
func (s *InvoiceService) GeneratePDF(userID int, invoiceID int, message string) (string, error) {
	invoice, err := s.Get(userID, invoiceID)
	if err != nil {
		return "", fmt.Errorf("invoice not found: %w", err)
	}

	client, err := s.getClient(userID, invoice.ClientID)
	if err != nil {
		return "", fmt.Errorf("client not found: %w", err)
	}

	settings, err := s.getUserSettings(userID)
	if err != nil {
		log.Println("Falling back to default settings due to error:", err)
	}

	finalMessage := strings.TrimSpace(message)
	if finalMessage == "" {
		finalMessage = s.buildDefaultMessage(userID, invoice.ID, settings)
	}

	// Ensure we use latest totals and items derived from linked time entries
	invoice = s.ensureInvoiceRecalcForPDF(userID, invoice, settings)

	// Use PDF Generator
	generator := pdf.NewGenerator(pdf.GetTemplatesDir())
	return generator.GeneratePDF(invoice, client, settings, finalMessage)
}

// ensureInvoiceRecalcForPDF recalculates invoice totals/items from linked time entries if present.
func (s *InvoiceService) ensureInvoiceRecalcForPDF(userID int, invoice dto.InvoiceOutput, _ models.UserSettings) dto.InvoiceOutput {
	updated, err := s.recalculateInvoiceFromTimeEntries(userID, invoice.ID, invoice.TaxRate)
	if err != nil {
		log.Println("Recalc before PDF failed, using stored invoice:", err)
		return invoice
	}
	return updated
}

// GetDefaultMessage exposes the default MESSAGE generation for frontend preview.
func (s *InvoiceService) GetDefaultMessage(userID int, invoiceID int) (string, error) {
	settings, err := s.getUserSettings(userID)
	if err != nil {
		log.Println("Falling back to default settings due to error:", err)
	}
	if _, err := s.Get(userID, invoiceID); err != nil {
		return "", fmt.Errorf("invoice not found: %w", err)
	}
	return s.buildDefaultMessage(userID, invoiceID, settings), nil
}

// getClient fetches client detail for invoices.
func (s *InvoiceService) getClient(userID int, clientID int) (models.Client, error) {
	row := s.db.QueryRow("SELECT id, name, email, website, avatar, contact_person, address, currency, status, notes, billing_company, billing_address, billing_city, billing_province, billing_postal_code FROM clients WHERE id = ? AND user_id = ?", clientID, userID)
	var c models.Client
	err := row.Scan(&c.ID, &c.Name, &c.Email, &c.Website, &c.Avatar, &c.ContactPerson, &c.Address, &c.Currency, &c.Status, &c.Notes, &c.BillingCompany, &c.BillingAddress, &c.BillingCity, &c.BillingProvince, &c.BillingPostalCode)
	if err != nil {
		return models.Client{}, err
	}
	return c, nil
}

// getUserSettings loads settings from the new dedicated tables (preferences, tax, invoice).
func (s *InvoiceService) getUserSettings(userID int) (models.UserSettings, error) {
	prefsSvc := NewUserPreferencesService(s.db)
	taxSvc := NewUserTaxSettingsService(s.db)
	invSvc := NewUserInvoiceSettingsService(s.db)

	// Fetch all settings (ignoring errors as services return defaults on error)
	prefs, _ := prefsSvc.Get(userID)
	tax, _ := taxSvc.Get(userID)
	inv, _ := invSvc.Get(userID)

	// Aggregate into legacy model
	return models.UserSettings{
		// Preferences
		Currency:        prefs.Currency,
		Language:        prefs.Language,
		Theme:           prefs.Theme,
		Timezone:        prefs.Timezone,
		DateFormat:      prefs.DateFormat,
		ModuleOverrides: prefs.ModuleOverrides,

		// Tax
		HstRegistered:  tax.HstRegistered,
		HstNumber:      tax.HstNumber,
		TaxEnabled:     tax.TaxEnabled,
		DefaultTaxRate: tax.DefaultTaxRate,
		ExpectedIncome: tax.ExpectedIncome,

		// Invoice & Sender
		SenderName:             inv.SenderName,
		SenderCompany:          inv.SenderCompany,
		SenderAddress:          inv.SenderAddress,
		SenderPhone:            inv.SenderPhone,
		SenderEmail:            inv.SenderEmail,
		SenderPostalCode:       inv.SenderPostalCode,
		InvoiceTerms:           inv.DefaultTerms,
		DefaultMessageTemplate: inv.DefaultMessageTemplate,
	}, nil
}

// buildDefaultMessage generates a fallback MESSAGE block based on time entries.
func (s *InvoiceService) buildDefaultMessage(userID int, invoiceID int, settings models.UserSettings) string {
	query := `
SELECT date, start_time, end_time, duration_seconds
FROM time_entries
WHERE user_id = ?
  AND invoice_id = ?
ORDER BY date ASC, start_time ASC`

	rows, err := s.db.Query(query, userID, invoiceID)
	if err != nil {
		log.Println("Error querying time entries for message:", err)
		return settings.DefaultMessageTemplate
	}
	defer closeWithLog(rows, "closing message time entries rows")

	var lines []string
	for rows.Next() {
		var date, start, end string
		var duration int
		if err := rows.Scan(&date, &start, &end, &duration); err != nil {
			log.Println("Error scanning time entry for message:", err)
			continue
		}
		hours := float64(duration) / 3600
		if start != "" && end != "" {
			lines = append(lines, fmt.Sprintf("%s %s-%s %.1f", date, start, end, hours))
		} else {
			lines = append(lines, fmt.Sprintf("%s %.1f hours", date, hours))
		}
	}

	if len(lines) == 0 {
		return settings.DefaultMessageTemplate
	}

	return strings.Join(lines, "\n")
}

// recalculateInvoiceFromTimeEntries recalculates subtotal/tax/total and items_json based on linked time entries.
func (s *InvoiceService) recalculateInvoiceFromTimeEntries(userID int, invoiceID int, taxRate float64) (dto.InvoiceOutput, error) {
	type entryRow struct {
		ProjectID   int
		Project     string
		Hourly      float64
		Currency    string
		ServiceType string
		Hours       float64
	}

	rows, err := s.db.Query(`
SELECT p.id, p.name, COALESCE(p.hourly_rate, 0), COALESCE(p.currency, ''), COALESCE(p.service_type, ''), te.duration_seconds
FROM time_entries te
JOIN projects p ON te.project_id = p.id
WHERE te.user_id = ? AND te.invoice_id = ?`, userID, invoiceID)
	if err != nil {
		return dto.InvoiceOutput{}, fmt.Errorf("failed to load time entries for invoice: %w", err)
	}
	defer closeWithLog(rows, "closing recalc time entries rows")

	projectHours := map[int]entryRow{}
	for rows.Next() {
		var pid int
		var name, currency, serviceType string
		var rate float64
		var seconds int
		if err := rows.Scan(&pid, &name, &rate, &currency, &serviceType, &seconds); err != nil {
			log.Println("Error scanning time entry for recalc:", err)
			continue
		}
		r := projectHours[pid]
		if r.ProjectID == 0 {
			r.ProjectID = pid
			r.Project = name
			r.Hourly = rate
			r.Currency = currency
			r.ServiceType = serviceType
		}
		r.Hours += float64(seconds) / 3600
		projectHours[pid] = r
	}

	var items []models.InvoiceItem
	var subtotal float64
	for _, r := range projectHours {
		amount := r.Hours * r.Hourly
		description := utils.FormatServiceType(r.ServiceType)
		if description == "" {
			description = r.Project
		}
		item := models.InvoiceItem{
			Description: description,
			Quantity:    r.Hours,
			UnitPrice:   r.Hourly,
			Amount:      amount,
		}
		items = append(items, item)
		subtotal += amount
	}

	// Update invoice record
	itemsBytes, _ := json.Marshal(items)
	itemsJSON := string(itemsBytes)
	taxAmount := subtotal * taxRate
	total := subtotal + taxAmount

	if _, err := s.db.Exec(`
UPDATE invoices
SET subtotal=?, tax_amount=?, total=?, items_json=?
WHERE id=? AND user_id=?`, subtotal, taxAmount, total, itemsJSON, invoiceID, userID); err != nil {
		return dto.InvoiceOutput{}, fmt.Errorf("failed to update invoice totals: %w", err)
	}

	// Fetch updated
	return s.Get(userID, invoiceID)
}

func (s *InvoiceService) sendViaSMTP(settings dto.InvoiceEmailSettings, toEmail, subject, body string, pdfBytes []byte, invoiceNumber string) error {
	// Setup auth
	auth := smtp.PlainAuth("", settings.SMTPUsername, settings.SMTPPassword, settings.SMTPHost)

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = settings.FromEmail
	headers["To"] = toEmail
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	boundary := "f46d043c8aa2b9211f43924705572551" // Random boundary
	headers["Content-Type"] = "multipart/mixed; boundary=" + boundary

	// Body buffer
	var msg bytes.Buffer
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")

	// Message body
	msg.WriteString("--" + boundary + "\r\n")
	msg.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(body)
	msg.WriteString("\r\n")

	// Attachment
	msg.WriteString("--" + boundary + "\r\n")
	msg.WriteString("Content-Type: application/pdf; name=\"invoice-" + invoiceNumber + ".pdf\"\r\n")
	msg.WriteString("Content-Transfer-Encoding: base64\r\n")
	msg.WriteString("Content-Disposition: attachment; filename=\"invoice-" + invoiceNumber + ".pdf\"\r\n")
	msg.WriteString("\r\n")

	b64 := base64.StdEncoding.EncodeToString(pdfBytes)
	// Split base64 lines (max 76 chars)
	for i := 0; i < len(b64); i += 76 {
		end := i + 76
		if end > len(b64) {
			end = len(b64)
		}
		msg.WriteString(b64[i:end] + "\r\n")
	}
	msg.WriteString("\r\n")
	msg.WriteString("--" + boundary + "--")

	// SMTP Connection (support TLS/StartTLS)
	addr := fmt.Sprintf("%s:%d", settings.SMTPHost, settings.SMTPPort)

	// If port 465, use implicit TLS
	if settings.SMTPPort == 465 {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: false,
			MinVersion:         tls.VersionTLS12,
			ServerName:         settings.SMTPHost,
		}
		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			log.Printf("SMTP: Auth failed: %v", err)
			return fmt.Errorf("authentication failed: %w", err)
		}

		client, err := smtp.NewClient(conn, settings.SMTPHost)
		if err != nil {
			log.Printf("SMTP: NewClient failed: %v", err)
			return fmt.Errorf("smtp client creation failed: %w", err)
		}
		defer func() {
			if err := client.Quit(); err != nil {
				log.Printf("SMTP: Quit failed: %v", err)
			}
		}()

		if err := client.Auth(auth); err != nil {
			log.Printf("SMTP: Auth failed: %v", err)
			return fmt.Errorf("authentication failed: %w", err)
		}

		if err := client.Mail(settings.FromEmail); err != nil {
			log.Printf("SMTP: Mail command failed: %v", err)
			return fmt.Errorf("mail command failed: %w", err)
		}

		if err := client.Rcpt(toEmail); err != nil {
			log.Printf("SMTP: Rcpt command failed: %v", err)
			return fmt.Errorf("rcpt command failed (check recipient address): %w", err)
		}

		w, err := client.Data()
		if err != nil {
			log.Printf("SMTP: Data command failed: %v", err)
			return fmt.Errorf("data command failed: %w", err)
		}

		if _, err := w.Write(msg.Bytes()); err != nil {
			log.Printf("SMTP: Write failed: %v", err)
			if closeErr := w.Close(); closeErr != nil {
				log.Printf("SMTP: Close writer after write failure failed: %v", closeErr)
			}
			return fmt.Errorf("message write failed: %w", err)
		}

		if err := w.Close(); err != nil {
			log.Printf("SMTP: Close writer failed: %v", err)
			return fmt.Errorf("message close failed: %w", err)
		}

		return nil
	}

	// Use plain SMTP without TLS
	client, err := smtp.Dial(addr)
	if err != nil {
		log.Printf("SMTP: Dial failed: %v", err)
		return fmt.Errorf("smtp dial failed: %w", err)
	}
	defer func() {
		if err := client.Quit(); err != nil {
			log.Printf("SMTP: Quit failed: %v", err)
		}
	}()

	if err := client.Auth(auth); err != nil {
		log.Printf("SMTP: Auth failed: %v", err)
		return fmt.Errorf("authentication failed: %w", err)
	}

	if err := client.Mail(settings.FromEmail); err != nil {
		log.Printf("SMTP: Mail command failed: %v", err)
		return fmt.Errorf("mail command failed: %w", err)
	}

	if err := client.Rcpt(toEmail); err != nil {
		log.Printf("SMTP: Rcpt command failed: %v", err)
		return fmt.Errorf("rcpt command failed: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		log.Printf("SMTP: Data command failed: %v", err)
		return fmt.Errorf("data command failed: %w", err)
	}

	if _, err := w.Write(msg.Bytes()); err != nil {
		log.Printf("SMTP: Write failed: %v", err)
		if closeErr := w.Close(); closeErr != nil {
			log.Printf("SMTP: Close writer after write failure failed: %v", closeErr)
		}
		return fmt.Errorf("message write failed: %w", err)
	}

	if err := w.Close(); err != nil {
		log.Printf("SMTP: Close writer failed: %v", err)
		return fmt.Errorf("message close failed: %w", err)
	}

	return nil
}
