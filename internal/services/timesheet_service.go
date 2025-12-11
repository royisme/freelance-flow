package services

import (
	"database/sql"
	"freelance-flow/internal/dto"
	"freelance-flow/internal/mapper"
	"freelance-flow/internal/models"
	"log"
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
	query := "SELECT id, project_id, invoice_id, date, start_time, end_time, duration_seconds, description, billable, invoiced FROM time_entries WHERE user_id = ?"
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
		err := rows.Scan(&t.ID, &t.ProjectID, &t.InvoiceID, &t.Date, &t.StartTime, &t.EndTime, &t.DurationSeconds, &t.Description, &t.Billable, &t.Invoiced)
		if err != nil {
			log.Println("Error scanning time entry:", err)
			continue
		}
		entries = append(entries, t)
	}
	return mapper.ToTimeEntryOutputList(entries)
}

// Get returns a single time entry by ID for a specific user.
func (s *TimesheetService) Get(userID int, id int) (dto.TimeEntryOutput, error) {
	row := s.db.QueryRow("SELECT id, project_id, invoice_id, date, start_time, end_time, duration_seconds, description, billable, invoiced FROM time_entries WHERE id = ? AND user_id = ?", id, userID)
	var t models.TimeEntry
	err := row.Scan(&t.ID, &t.ProjectID, &t.InvoiceID, &t.Date, &t.StartTime, &t.EndTime, &t.DurationSeconds, &t.Description, &t.Billable, &t.Invoiced)
	if err != nil {
		return dto.TimeEntryOutput{}, err
	}
	return mapper.ToTimeEntryOutput(t), nil
}

// Create adds a new time entry for a specific user and returns the created entry as DTO.
func (s *TimesheetService) Create(userID int, input dto.CreateTimeEntryInput) dto.TimeEntryOutput {
	entity := mapper.ToTimeEntryEntity(input)

	stmt, err := s.db.Prepare("INSERT INTO time_entries(user_id, project_id, invoice_id, date, start_time, end_time, duration_seconds, description, billable, invoiced) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("Error preparing time entry insert:", err)
		return dto.TimeEntryOutput{}
	}
	defer closeWithLog(stmt, "closing time entry insert statement")

	res, err := stmt.Exec(userID, entity.ProjectID, entity.InvoiceID, entity.Date, entity.StartTime, entity.EndTime, entity.DurationSeconds, entity.Description, entity.Billable, entity.Invoiced)
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
	stmt, err := s.db.Prepare("UPDATE time_entries SET project_id=?, invoice_id=?, date=?, start_time=?, end_time=?, duration_seconds=?, description=?, billable=?, invoiced=? WHERE id=? AND user_id=?")
	if err != nil {
		log.Println("Error preparing time entry update:", err)
		return dto.TimeEntryOutput{}
	}
	defer closeWithLog(stmt, "closing time entry update statement")

	_, err = stmt.Exec(input.ProjectID, input.InvoiceID, input.Date, input.StartTime, input.EndTime, input.DurationSeconds, input.Description, input.Billable, input.Invoiced, input.ID, userID)
	if err != nil {
		log.Println("Error updating time entry:", err)
		return dto.TimeEntryOutput{}
	}

	output, _ := s.Get(userID, input.ID)
	return output
}

// Delete removes a time entry by ID for a specific user.
func (s *TimesheetService) Delete(userID int, id int) {
	_, err := s.db.Exec("DELETE FROM time_entries WHERE id=? AND user_id=?", id, userID)
	if err != nil {
		log.Println("Error deleting time entry:", err)
	}
}
