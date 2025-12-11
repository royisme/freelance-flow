package services

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"freelance-flow/internal/dto"
	"freelance-flow/internal/mapper"
	"freelance-flow/internal/models"
	"log"
	"os"
	"strings"

	"github.com/go-pdf/fpdf"
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
	rows, err := s.db.Query("SELECT id, client_id, number, issue_date, due_date, subtotal, tax_rate, tax_amount, total, status, items_json FROM invoices WHERE user_id = ?", userID)
	if err != nil {
		log.Println("Error querying invoices:", err)
		return []dto.InvoiceOutput{}
	}
	defer closeWithLog(rows, "closing invoice rows")

	var invoices []models.Invoice
	for rows.Next() {
		var i models.Invoice
		var itemsJSON string
		err := rows.Scan(&i.ID, &i.ClientID, &i.Number, &i.IssueDate, &i.DueDate, &i.Subtotal, &i.TaxRate, &i.TaxAmount, &i.Total, &i.Status, &itemsJSON)
		if err != nil {
			log.Println("Error scanning invoice:", err)
			continue
		}
		if itemsJSON != "" {
			if err := json.Unmarshal([]byte(itemsJSON), &i.Items); err != nil {
				log.Printf("Error unmarshalling items for invoice %d: %v", i.ID, err)
				i.Items = []models.InvoiceItem{}
			}
		} else {
			i.Items = []models.InvoiceItem{}
		}
		invoices = append(invoices, i)
	}
	return mapper.ToInvoiceOutputList(invoices)
}

// Get returns a single invoice by ID for a specific user.
func (s *InvoiceService) Get(userID int, id int) (dto.InvoiceOutput, error) {
	row := s.db.QueryRow("SELECT id, client_id, number, issue_date, due_date, subtotal, tax_rate, tax_amount, total, status, items_json FROM invoices WHERE id = ? AND user_id = ?", id, userID)
	var i models.Invoice
	var itemsJSON string
	err := row.Scan(&i.ID, &i.ClientID, &i.Number, &i.IssueDate, &i.DueDate, &i.Subtotal, &i.TaxRate, &i.TaxAmount, &i.Total, &i.Status, &itemsJSON)
	if err != nil {
		return dto.InvoiceOutput{}, err
	}
	if itemsJSON != "" {
		if err := json.Unmarshal([]byte(itemsJSON), &i.Items); err != nil {
			log.Printf("Error unmarshalling items for invoice %d: %v", i.ID, err)
			return dto.InvoiceOutput{}, fmt.Errorf("failed to unmarshal items: %w", err)
		}
	} else {
		i.Items = []models.InvoiceItem{}
	}
	return mapper.ToInvoiceOutput(i), nil
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

// Delete removes an invoice by ID for a specific user.
func (s *InvoiceService) Delete(userID int, id int) {
	_, err := s.db.Exec("DELETE FROM invoices WHERE id=? AND user_id=?", id, userID)
	if err != nil {
		log.Println("Error deleting invoice:", err)
	}
}

// SendEmail sends invoice via configured provider (mailto/resend/smtp placeholder).
func (s *InvoiceService) SendEmail(userID int, invoiceID int) bool {
	invoice, err := s.Get(userID, invoiceID)
	if err != nil {
		log.Println("SendEmail: invoice not found:", err)
		return false
	}
	client, err := s.getClient(userID, invoice.ClientID)
	if err != nil {
		log.Println("SendEmail: client not found:", err)
		return false
	}

	settingsSvc := NewInvoiceEmailSettingsService(s.db)
	emailSettings := settingsSvc.Get(userID)
	if emailSettings.Provider == "" {
		emailSettings.Provider = "mailto"
	}

	pdfBase64, err := s.GeneratePDF(userID, invoiceID, "")
	if err != nil {
		log.Println("SendEmail: generate pdf failed:", err)
		return false
	}
	pdfBytes, err := base64.StdEncoding.DecodeString(pdfBase64)
	if err != nil {
		log.Println("SendEmail: decode pdf failed:", err)
		return false
	}

	// mailto just indicate success (front-end opens client)
	if emailSettings.Provider == "mailto" {
		return true
	}

	// Resend provider
	if emailSettings.Provider == "resend" {
		apiKey := emailSettings.ResendAPIKey
		if apiKey == "" {
			log.Println("SendEmail: resend api key missing")
			return false
		}
		subject := applyTemplate(emailSettings.SubjectTemplate, invoice)
		if subject == "" {
			subject = fmt.Sprintf("Invoice %s", invoice.Number)
		}
		body := applyTemplate(emailSettings.BodyTemplate, invoice)
		if body == "" {
			body = "Please see attached invoice."
		}

		if os.Getenv("RESEND_DRY_RUN") == "1" {
			log.Println("SendEmail: RESEND_DRY_RUN enabled, skipping network call")
			return true
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
			return false
		}
		return true
	}

	// TODO: SMTP provider can be added similarly.
	return false
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

	pdf, err := s.renderPDF(invoice, client, settings, finalMessage)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return "", fmt.Errorf("failed to render pdf: %w", err)
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

// getClient fetches client detail for invoices.
func (s *InvoiceService) getClient(userID int, clientID int) (models.Client, error) {
	row := s.db.QueryRow("SELECT id, name, email, website, avatar, contact_person, address, currency, status, notes FROM clients WHERE id = ? AND user_id = ?", clientID, userID)
	var c models.Client
	err := row.Scan(&c.ID, &c.Name, &c.Email, &c.Website, &c.Avatar, &c.ContactPerson, &c.Address, &c.Currency, &c.Status, &c.Notes)
	if err != nil {
		return models.Client{}, err
	}
	return c, nil
}

// getUserSettings loads settings_json for a user.
func (s *InvoiceService) getUserSettings(userID int) (models.UserSettings, error) {
	raw := "{}"
	if err := s.db.QueryRow("SELECT settings_json FROM users WHERE id = ?", userID).Scan(&raw); err != nil {
		return defaultUserSettings(), err
	}
	settings := defaultUserSettings()
	if err := json.Unmarshal([]byte(raw), &settings); err != nil {
		return defaultUserSettings(), err
	}
	return normalizeUserSettings(settings), nil
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

// ensureInvoiceRecalcForPDF recalculates invoice totals/items from linked time entries if present.
func (s *InvoiceService) ensureInvoiceRecalcForPDF(userID int, invoice dto.InvoiceOutput, _ models.UserSettings) dto.InvoiceOutput {
	updated, err := s.recalculateInvoiceFromTimeEntries(userID, invoice.ID, invoice.TaxRate)
	if err != nil {
		log.Println("Recalc before PDF failed, using stored invoice:", err)
		return invoice
	}
	return updated
}

// renderPDF draws the invoice PDF with the expected layout.
func (s *InvoiceService) renderPDF(invoice dto.InvoiceOutput, client models.Client, settings models.UserSettings, message string) (*fpdf.Fpdf, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 20, 15)
	pdf.AddPage()

	headerFill := func() {
		pdf.SetFillColor(51, 51, 51)
		pdf.Rect(0, 0, 210, 35, "F")
		pdf.SetTextColor(255, 255, 255)
		pdf.SetFont("Helvetica", "B", 28)
		pdf.SetXY(150, 12)
		pdf.CellFormat(50, 10, "INVOICE", "", 0, "R", false, 0, "")

		pdf.SetFont("Helvetica", "", 11)
		pdf.SetXY(15, 12)
		senderName := settings.SenderName
		if senderName == "" {
			senderName = settings.SenderCompany
		}
		if senderName == "" {
			senderName = "Sender"
		}
		pdf.CellFormat(120, 6, senderName, "", 1, "L", false, 0, "")
		if settings.SenderAddress != "" {
			pdf.CellFormat(120, 6, settings.SenderAddress, "", 1, "L", false, 0, "")
		}
		if settings.SenderPostalCode != "" {
			pdf.CellFormat(120, 6, settings.SenderPostalCode, "", 1, "L", false, 0, "")
		}
		if settings.SenderPhone != "" {
			pdf.CellFormat(120, 6, settings.SenderPhone, "", 1, "L", false, 0, "")
		}
		if settings.SenderEmail != "" {
			pdf.CellFormat(120, 6, settings.SenderEmail, "", 1, "L", false, 0, "")
		}
		pdf.SetTextColor(0, 0, 0)
	}

	sectionInvoiceInfo := func() {
		pdf.Ln(10)
		pdf.SetFont("Helvetica", "B", 11)
		pdf.CellFormat(90, 6, "INVOICE TO "+strings.ToUpper(client.Name), "", 0, "L", false, 0, "")
		pdf.CellFormat(90, 6, fmt.Sprintf("INVOICE# %s", invoice.Number), "", 1, "R", false, 0, "")

		pdf.SetFont("Helvetica", "", 10)
		pdf.CellFormat(90, 6, client.Address, "", 0, "L", false, 0, "")
		rightY := pdf.GetY()
		pdf.SetXY(105, rightY)
		pdf.CellFormat(90, 6, fmt.Sprintf("DATE %s", invoice.IssueDate), "", 1, "R", false, 0, "")

		pdf.CellFormat(90, 6, fmt.Sprintf("%s %s", client.ContactPerson, client.Email), "", 0, "L", false, 0, "")
		pdf.SetXY(105, pdf.GetY()-6)
		due := invoice.DueDate
		if due == "" {
			due = "DUE DATE"
		}
		pdf.CellFormat(90, 6, fmt.Sprintf("DUE DATE %s", due), "", 1, "R", false, 0, "")

		pdf.SetXY(105, pdf.GetY()-6)
		pdf.CellFormat(90, 6, fmt.Sprintf("TERMS %s", settings.InvoiceTerms), "", 1, "R", false, 0, "")
		pdf.Ln(6)
	}

	tableItems := func() {
		pdf.SetFillColor(12, 168, 67)
		pdf.SetTextColor(255, 255, 255)
		pdf.SetFont("Helvetica", "B", 10)
		pdf.CellFormat(100, 8, "DESCRIPTION", "1", 0, "L", true, 0, "")
		pdf.CellFormat(30, 8, "QTY", "1", 0, "C", true, 0, "")
		pdf.CellFormat(30, 8, "RATE", "1", 0, "C", true, 0, "")
		pdf.CellFormat(30, 8, "AMOUNT", "1", 1, "C", true, 0, "")

		pdf.SetTextColor(0, 0, 0)
		pdf.SetFont("Helvetica", "", 10)
		description := "Invoice Summary"
		if len(invoice.Items) > 0 && invoice.Items[0].Description != "" {
			description = invoice.Items[0].Description
		}
		qty := 0.0
		rate := 0.0
		for _, item := range invoice.Items {
			qty += item.Quantity
			rate = item.UnitPrice
		}
		if qty == 0 && rate > 0 {
			qty = 1
		}
		amount := invoice.Total

		pdf.CellFormat(100, 10, description, "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("%.2f", qty), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("%.2f %s", rate, settings.Currency), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("%.2f %s", amount, settings.Currency), "1", 1, "C", false, 0, "")
	}

	totalSection := func() {
		pdf.Ln(4)
		startX := 110.0
		pdf.SetXY(startX, pdf.GetY())
		pdf.SetFont("Helvetica", "", 10)
		rows := []struct {
			label string
			value string
		}{
			{"SUBTOTAL", fmt.Sprintf("%.2f %s", invoice.Subtotal, settings.Currency)},
			{"DISCOUNT", "0"},
			{"TAX", fmt.Sprintf("%.2f %s", invoice.TaxAmount, settings.Currency)},
			{"TOTAL", fmt.Sprintf("%.2f %s", invoice.Total, settings.Currency)},
			{"BALANCE DUE", fmt.Sprintf("%.2f %s", invoice.Total, settings.Currency)},
		}
		for _, row := range rows {
			pdf.CellFormat(40, 8, row.label, "", 0, "L", false, 0, "")
			pdf.CellFormat(40, 8, row.value, "", 1, "R", false, 0, "")
		}
	}

	messageSection := func() {
		pdf.Ln(6)
		pdf.SetFont("Helvetica", "B", 10)
		pdf.CellFormat(60, 6, "MESSAGE", "", 1, "L", false, 0, "")
		pdf.SetFont("Helvetica", "", 10)
		pdf.MultiCell(120, 6, message, "1", "L", false)
	}

	headerFill()
	sectionInvoiceInfo()
	tableItems()
	messageSection()
	totalSection()

	return pdf, nil
}

// recalculateInvoiceFromTimeEntries recalculates subtotal/tax/total and items_json based on linked time entries.
func (s *InvoiceService) recalculateInvoiceFromTimeEntries(userID int, invoiceID int, taxRate float64) (dto.InvoiceOutput, error) {
	type entryRow struct {
		ProjectID int
		Project   string
		Hourly    float64
		Currency  string
		Hours     float64
	}

	rows, err := s.db.Query(`
SELECT p.id, p.name, COALESCE(p.hourly_rate, 0), COALESCE(p.currency, ''), te.duration_seconds
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
		var name, currency string
		var rate float64
		var seconds int
		if err := rows.Scan(&pid, &name, &rate, &currency, &seconds); err != nil {
			log.Println("Error scanning time entry for recalc:", err)
			continue
		}
		r := projectHours[pid]
		if r.ProjectID == 0 {
			r.ProjectID = pid
			r.Project = name
			r.Hourly = rate
			r.Currency = currency
		}
		r.Hours += float64(seconds) / 3600
		projectHours[pid] = r
	}

	var items []models.InvoiceItem
	var subtotal float64
	for _, r := range projectHours {
		amount := r.Hours * r.Hourly
		item := models.InvoiceItem{
			Description: r.Project,
			Quantity:    r.Hours,
			UnitPrice:   r.Hourly,
			Amount:      amount,
		}
		items = append(items, item)
		subtotal += amount
	}

	taxAmount := subtotal * taxRate
	total := subtotal + taxAmount

	itemsBytes, _ := json.Marshal(items)
	itemsJSON := string(itemsBytes)

	_, err = s.db.Exec(`
UPDATE invoices
SET subtotal = ?, tax_amount = ?, total = ?, items_json = ?
WHERE id = ? AND user_id = ?`,
		subtotal, taxAmount, total, itemsJSON, invoiceID, userID)
	if err != nil {
		return dto.InvoiceOutput{}, fmt.Errorf("failed to update invoice totals: %w", err)
	}

	// Return refreshed invoice dto
	return s.Get(userID, invoiceID)
}
