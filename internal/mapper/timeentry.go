package mapper

import (
	"tally/internal/dto"
	"tally/internal/models"
)

// ToTimeEntryOutput converts a TimeEntry entity to TimeEntryOutput DTO.
func ToTimeEntryOutput(e models.TimeEntry) dto.TimeEntryOutput {
	return dto.TimeEntryOutput{
		ID:              e.ID,
		ProjectID:       e.ProjectID,
		InvoiceID:       e.InvoiceID,
		Date:            e.Date,
		StartTime:       e.StartTime,
		EndTime:         e.EndTime,
		DurationSeconds: e.DurationSeconds,
		Description:     e.Description,
		Billable:        e.Billable,
		Invoiced:        e.Invoiced,
		BillingMode:     e.BillingMode,
		ManualAmount:    e.ManualAmount,
	}
}

// ToTimeEntryOutputList converts a slice of TimeEntry entities to TimeEntryOutput DTOs.
func ToTimeEntryOutputList(entities []models.TimeEntry) []dto.TimeEntryOutput {
	if entities == nil {
		return []dto.TimeEntryOutput{}
	}
	result := make([]dto.TimeEntryOutput, len(entities))
	for i, e := range entities {
		result[i] = ToTimeEntryOutput(e)
	}
	return result
}

// ToTimeEntryEntity converts CreateTimeEntryInput DTO to TimeEntry entity.
// Note: InvoiceID defaults to 0 (unassigned) for new entries.
func ToTimeEntryEntity(input dto.CreateTimeEntryInput) models.TimeEntry {
	billingMode := input.BillingMode
	if billingMode == "" {
		billingMode = "hourly"
	}
	return models.TimeEntry{
		ProjectID:       input.ProjectID,
		InvoiceID:       0, // New entries are not assigned to an invoice
		Date:            input.Date,
		StartTime:       input.StartTime,
		EndTime:         input.EndTime,
		DurationSeconds: input.DurationSeconds,
		Description:     input.Description,
		Billable:        input.Billable,
		Invoiced:        input.Invoiced,
		BillingMode:     billingMode,
		ManualAmount:    input.ManualAmount,
	}
}

// ApplyTimeEntryUpdate applies UpdateTimeEntryInput to an existing TimeEntry entity.
func ApplyTimeEntryUpdate(e *models.TimeEntry, input dto.UpdateTimeEntryInput) {
	e.ProjectID = input.ProjectID
	e.InvoiceID = input.InvoiceID
	e.Date = input.Date
	e.StartTime = input.StartTime
	e.EndTime = input.EndTime
	e.DurationSeconds = input.DurationSeconds
	e.Description = input.Description
	e.Billable = input.Billable
	e.Invoiced = input.Invoiced
	e.BillingMode = input.BillingMode
	e.ManualAmount = input.ManualAmount
}
