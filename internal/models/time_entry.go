// Package models defines database-backed domain models.
package models

// TimeEntry records tracked work for a project.
type TimeEntry struct {
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
