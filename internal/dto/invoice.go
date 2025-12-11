package dto

// InvoiceItemInput represents an invoice line item in input.
type InvoiceItemInput struct {
	Description string  `json:"description"`
	Quantity    float64 `json:"quantity"`
	UnitPrice   float64 `json:"unitPrice"`
	Amount      float64 `json:"amount"`
}

// InvoiceItemOutput represents an invoice line item in output.
type InvoiceItemOutput struct {
	ID          int     `json:"id"`
	Description string  `json:"description"`
	Quantity    float64 `json:"quantity"`
	UnitPrice   float64 `json:"unitPrice"`
	Amount      float64 `json:"amount"`
}

// CreateInvoiceInput represents the input for creating a new invoice.
type CreateInvoiceInput struct {
	ClientID  int                `json:"clientId"`
	Number    string             `json:"number"`
	IssueDate string             `json:"issueDate"`
	DueDate   string             `json:"dueDate"`
	Subtotal  float64            `json:"subtotal"`
	TaxRate   float64            `json:"taxRate"`
	TaxAmount float64            `json:"taxAmount"`
	Total     float64            `json:"total"`
	Status    string             `json:"status"`
	Items     []InvoiceItemInput `json:"items"`
}

// UpdateInvoiceInput represents the input for updating an existing invoice.
type UpdateInvoiceInput struct {
	ID        int                `json:"id"`
	ClientID  int                `json:"clientId"`
	Number    string             `json:"number"`
	IssueDate string             `json:"issueDate"`
	DueDate   string             `json:"dueDate"`
	Subtotal  float64            `json:"subtotal"`
	TaxRate   float64            `json:"taxRate"`
	TaxAmount float64            `json:"taxAmount"`
	Total     float64            `json:"total"`
	Status    string             `json:"status"`
	Items     []InvoiceItemInput `json:"items"`
}

// InvoiceOutput represents the invoice data returned from API.
type InvoiceOutput struct {
	ID        int                 `json:"id"`
	ClientID  int                 `json:"clientId"`
	Number    string              `json:"number"`
	IssueDate string              `json:"issueDate"`
	DueDate   string              `json:"dueDate"`
	Subtotal  float64             `json:"subtotal"`
	TaxRate   float64             `json:"taxRate"`
	TaxAmount float64             `json:"taxAmount"`
	Total     float64             `json:"total"`
	Status    string              `json:"status"`
	Items     []InvoiceItemOutput `json:"items"`
}

// SetInvoiceTimeEntriesInput links time entries to an invoice.
type SetInvoiceTimeEntriesInput struct {
	InvoiceID    int   `json:"invoiceId"`
	TimeEntryIDs []int `json:"timeEntryIds"`
}
