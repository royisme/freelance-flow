package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"freelance-flow/internal/dto"
	"freelance-flow/internal/mapper"
	"freelance-flow/internal/models"
	"log"
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

// GeneratePDF is a placeholder for PDF generation.
func (s *InvoiceService) GeneratePDF(_ int, _ int) string {
	// TODO: implement real PDF generation using user and invoice context.
	return "mock-pdf-base64-from-backend"
}

// SendEmail is a placeholder for email sending.
func (s *InvoiceService) SendEmail(_ int, _ int) bool {
	// TODO: implement real email sending using user and invoice context.
	return true
}
