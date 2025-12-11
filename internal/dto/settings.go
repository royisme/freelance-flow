package dto

// UserSettings represents user preferences stored in settings_json.
// This is exposed to the frontend via Wails.
type UserSettings struct {
	Currency               string  `json:"currency"`
	DefaultTaxRate         float64 `json:"defaultTaxRate"`
	Language               string  `json:"language"`
	Theme                  string  `json:"theme"`
	DateFormat             string  `json:"dateFormat"`
	Timezone               string  `json:"timezone"`
	SenderName             string  `json:"senderName"`
	SenderCompany          string  `json:"senderCompany"`
	SenderAddress          string  `json:"senderAddress"`
	SenderPhone            string  `json:"senderPhone"`
	SenderEmail            string  `json:"senderEmail"`
	SenderPostalCode       string  `json:"senderPostalCode"`
	InvoiceTerms           string  `json:"invoiceTerms"`
	DefaultMessageTemplate string  `json:"defaultMessageTemplate"`
}
