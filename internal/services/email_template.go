package services

import (
	"strings"
	"text/template"

	"freelance-flow/internal/dto"
)

// applyTemplate renders a simple text template using invoice fields.
func applyTemplate(tpl string, invoice dto.InvoiceOutput) string {
	if strings.TrimSpace(tpl) == "" {
		return ""
	}
	data := map[string]interface{}{
		"number":   invoice.Number,
		"total":    invoice.Total,
		"due_date": invoice.DueDate,
		"issue":    invoice.IssueDate,
	}
	t, err := template.New("tpl").Parse(tpl)
	if err != nil {
		return ""
	}
	var b strings.Builder
	if err := t.Execute(&b, data); err != nil {
		return ""
	}
	return b.String()
}
