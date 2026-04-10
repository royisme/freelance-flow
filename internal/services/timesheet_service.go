package services

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"tally/internal/dto"
	"tally/internal/mapper"
	"tally/internal/models"
)

// TimesheetService handles all time entry-related operations.
type TimesheetService struct {
	db *sql.DB
}

// NewTimesheetService creates a new TimesheetService instance.
func NewTimesheetService(db *sql.DB) *TimesheetService {
	return &TimesheetService{db: db}
}

// List returns all time entries for a specific user, optionally filtered by project ID.
func (s *TimesheetService) List(userID int, projectID int) []dto.TimeEntryOutput {
	query := "SELECT id, project_id, invoice_id, date, start_time, end_time, duration_seconds, description, billable, invoiced, billing_mode, manual_amount FROM time_entries WHERE user_id = ?"
	args := []interface{}{userID}
	if projectID > 0 {
		query += " AND project_id = ?"
		args = append(args, projectID)
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		log.Println("Error querying time entries:", err)
		return []dto.TimeEntryOutput{}
	}
	defer closeWithLog(rows, "closing time entry rows")

	var entries []models.TimeEntry
	for rows.Next() {
		var t models.TimeEntry
		var invoiceID sql.NullInt64
		var manualAmount sql.NullFloat64
		err := rows.Scan(&t.ID, &t.ProjectID, &invoiceID, &t.Date, &t.StartTime, &t.EndTime, &t.DurationSeconds, &t.Description, &t.Billable, &t.Invoiced, &t.BillingMode, &manualAmount)
		if err != nil {
			log.Println("Error scanning time entry:", err)
			continue
		}
		if invoiceID.Valid {
			t.InvoiceID = int(invoiceID.Int64)
		} else {
			t.InvoiceID = 0
		}
		if manualAmount.Valid {
			amount := manualAmount.Float64
			t.ManualAmount = &amount
		}
		entries = append(entries, t)
	}
	return mapper.ToTimeEntryOutputList(entries)
}

// Get returns a single time entry by ID for a specific user.
func (s *TimesheetService) Get(userID int, id int) (dto.TimeEntryOutput, error) {
	row := s.db.QueryRow("SELECT id, project_id, invoice_id, date, start_time, end_time, duration_seconds, description, billable, invoiced, billing_mode, manual_amount FROM time_entries WHERE id = ? AND user_id = ?", id, userID)
	var t models.TimeEntry
	var invoiceID sql.NullInt64
	var manualAmount sql.NullFloat64
	err := row.Scan(&t.ID, &t.ProjectID, &invoiceID, &t.Date, &t.StartTime, &t.EndTime, &t.DurationSeconds, &t.Description, &t.Billable, &t.Invoiced, &t.BillingMode, &manualAmount)
	if err != nil {
		return dto.TimeEntryOutput{}, err
	}
	if invoiceID.Valid {
		t.InvoiceID = int(invoiceID.Int64)
	} else {
		t.InvoiceID = 0
	}
	if manualAmount.Valid {
		amount := manualAmount.Float64
		t.ManualAmount = &amount
	}
	return mapper.ToTimeEntryOutput(t), nil
}

// Create adds a new time entry for a specific user and returns the created entry as DTO.
func (s *TimesheetService) Create(userID int, input dto.CreateTimeEntryInput) dto.TimeEntryOutput {
	if err := validateTimeEntryInput(input.BillingMode, input.DurationSeconds, input.ManualAmount); err != nil {
		log.Println("Error validating time entry create:", err)
		return dto.TimeEntryOutput{}
	}
	entity := mapper.ToTimeEntryEntity(input)

	stmt, err := s.db.Prepare("INSERT INTO time_entries(user_id, project_id, invoice_id, date, start_time, end_time, duration_seconds, description, billable, invoiced, billing_mode, manual_amount) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("Error preparing time entry insert:", err)
		return dto.TimeEntryOutput{}
	}
	defer closeWithLog(stmt, "closing time entry insert statement")

	res, err := stmt.Exec(userID, entity.ProjectID, entity.InvoiceID, entity.Date, entity.StartTime, entity.EndTime, entity.DurationSeconds, entity.Description, entity.Billable, entity.Invoiced, entity.BillingMode, entity.ManualAmount)
	if err != nil {
		log.Println("Error inserting time entry:", err)
		return dto.TimeEntryOutput{}
	}

	id, _ := res.LastInsertId()
	entity.ID = int(id)
	return mapper.ToTimeEntryOutput(entity)
}

// Update modifies an existing time entry for a specific user and returns the updated entry as DTO.
func (s *TimesheetService) Update(userID int, input dto.UpdateTimeEntryInput) dto.TimeEntryOutput {
	if err := validateTimeEntryInput(input.BillingMode, input.DurationSeconds, input.ManualAmount); err != nil {
		log.Println("Error validating time entry update:", err)
		return dto.TimeEntryOutput{}
	}
	stmt, err := s.db.Prepare("UPDATE time_entries SET project_id=?, invoice_id=?, date=?, start_time=?, end_time=?, duration_seconds=?, description=?, billable=?, invoiced=?, billing_mode=?, manual_amount=? WHERE id=? AND user_id=?")
	if err != nil {
		log.Println("Error preparing time entry update:", err)
		return dto.TimeEntryOutput{}
	}
	defer closeWithLog(stmt, "closing time entry update statement")

	_, err = stmt.Exec(input.ProjectID, input.InvoiceID, input.Date, input.StartTime, input.EndTime, input.DurationSeconds, input.Description, input.Billable, input.Invoiced, input.BillingMode, input.ManualAmount, input.ID, userID)
	if err != nil {
		log.Println("Error updating time entry:", err)
		return dto.TimeEntryOutput{}
	}

	output, _ := s.Get(userID, input.ID)
	return output
}

func validateTimeEntryInput(billingMode string, durationSeconds int, manualAmount *float64) error {
	mode := billingMode
	if mode == "" {
		mode = "hourly"
	}

	switch mode {
	case "hourly":
		if durationSeconds <= 0 {
			return errors.New("hourly billing requires durationSeconds > 0")
		}
		if manualAmount != nil {
			return errors.New("hourly billing cannot include manualAmount")
		}
	case "fixed":
		if manualAmount == nil {
			return errors.New("fixed billing requires manualAmount")
		}
		if *manualAmount < 0 {
			return errors.New("fixed billing manualAmount must be >= 0")
		}
		if durationSeconds < 0 {
			return errors.New("fixed billing durationSeconds must be >= 0")
		}
	default:
		return fmt.Errorf("unsupported billing mode: %s", mode)
	}

	return nil
}

// Delete removes a time entry by ID for a specific user.
func (s *TimesheetService) Delete(userID int, id int) {
	_, err := s.db.Exec("DELETE FROM time_entries WHERE id=? AND user_id=?", id, userID)
	if err != nil {
		log.Println("Error deleting time entry:", err)
	}
}
