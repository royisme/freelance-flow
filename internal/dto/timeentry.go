package dto

// CreateTimeEntryInput represents the input for creating a new time entry.
type CreateTimeEntryInput struct {
	ProjectID       int    `json:"projectId"`
	InvoiceID       int    `json:"invoiceId"`
	Date            string `json:"date"`
	StartTime       string `json:"startTime"`
	EndTime         string `json:"endTime"`
	DurationSeconds int    `json:"durationSeconds"`
	Description     string `json:"description"`
	Billable        bool   `json:"billable"`
	Invoiced        bool   `json:"invoiced"`
}

// UpdateTimeEntryInput represents the input for updating an existing time entry.
type UpdateTimeEntryInput struct {
	ID              int    `json:"id"`
	ProjectID       int    `json:"projectId"`
	InvoiceID       int    `json:"invoiceId"`
	Date            string `json:"date"`
	StartTime       string `json:"startTime"`
	EndTime         string `json:"endTime"`
	DurationSeconds int    `json:"durationSeconds"`
	Description     string `json:"description"`
	Billable        bool   `json:"billable"`
	Invoiced        bool   `json:"invoiced"`
}

// TimeEntryOutput represents the time entry data returned from API.
type TimeEntryOutput struct {
	ID              int    `json:"id"`
	ProjectID       int    `json:"projectId"`
	InvoiceID       int    `json:"invoiceId"`
	Date            string `json:"date"`
	StartTime       string `json:"startTime"`
	EndTime         string `json:"endTime"`
	DurationSeconds int    `json:"durationSeconds"`
	Description     string `json:"description"`
	Billable        bool   `json:"billable"`
	Invoiced        bool   `json:"invoiced"`
}
