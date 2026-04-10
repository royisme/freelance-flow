package dto

// CreateTimeEntryInput represents the input for creating a new time entry.
// Note: InvoiceID is intentionally omitted - time entries are assigned to invoices later.
type CreateTimeEntryInput struct {
	ProjectID       int      `json:"projectId"`
	Date            string   `json:"date"`
	StartTime       string   `json:"startTime"`
	EndTime         string   `json:"endTime"`
	DurationSeconds int      `json:"durationSeconds"`
	Description     string   `json:"description"`
	Billable        bool     `json:"billable"`
	Invoiced        bool     `json:"invoiced"`
	BillingMode     string   `json:"billingMode"`
	ManualAmount    *float64 `json:"manualAmount"`
}

// UpdateTimeEntryInput represents the input for updating an existing time entry.
type UpdateTimeEntryInput struct {
	ID              int      `json:"id"`
	ProjectID       int      `json:"projectId"`
	InvoiceID       int      `json:"invoiceId"`
	Date            string   `json:"date"`
	StartTime       string   `json:"startTime"`
	EndTime         string   `json:"endTime"`
	DurationSeconds int      `json:"durationSeconds"`
	Description     string   `json:"description"`
	Billable        bool     `json:"billable"`
	Invoiced        bool     `json:"invoiced"`
	BillingMode     string   `json:"billingMode"`
	ManualAmount    *float64 `json:"manualAmount"`
}

// TimeEntryOutput represents the time entry data returned from API.
type TimeEntryOutput struct {
	ID              int      `json:"id"`
	ProjectID       int      `json:"projectId"`
	InvoiceID       int      `json:"invoiceId"`
	Date            string   `json:"date"`
	StartTime       string   `json:"startTime"`
	EndTime         string   `json:"endTime"`
	DurationSeconds int      `json:"durationSeconds"`
	Description     string   `json:"description"`
	Billable        bool     `json:"billable"`
	Invoiced        bool     `json:"invoiced"`
	BillingMode     string   `json:"billingMode"`
	ManualAmount    *float64 `json:"manualAmount"`
}
