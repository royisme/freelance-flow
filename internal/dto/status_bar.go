package dto

// StatusBarOutput provides lightweight aggregate metrics for the app footer/status bar.
type StatusBarOutput struct {
	MonthSeconds    int     `json:"monthSeconds"`
	UninvoicedTotal float64 `json:"uninvoicedTotal"`
	UnpaidTotal     float64 `json:"unpaidTotal"`
	Currency        string  `json:"currency"`
}
