package services

import (
	"database/sql"
	"freelance-flow/internal/dto"
	"time"
)

// StatusBarService provides aggregate metrics for the footer/status bar.
type StatusBarService struct {
	db *sql.DB
}

// NewStatusBarService creates a StatusBarService instance.
func NewStatusBarService(db *sql.DB) *StatusBarService {
	return &StatusBarService{db: db}
}

// Get returns aggregate status bar metrics for the given user.
func (s *StatusBarService) Get(userID int) (dto.StatusBarOutput, error) {
	settingsSvc := NewSettingsService(s.db)
	settings := settingsSvc.Get(userID)

	loc := time.UTC
	if settings.Timezone != "" {
		if parsed, err := time.LoadLocation(settings.Timezone); err == nil {
			loc = parsed
		}
	}

	now := time.Now().In(loc)
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
	nextMonthStart := monthStart.AddDate(0, 1, 0)

	startDate := monthStart.Format("2006-01-02")
	endDateExclusive := nextMonthStart.Format("2006-01-02")

	var monthSeconds int
	if err := s.db.QueryRow(
		`SELECT COALESCE(SUM(duration_seconds), 0) FROM time_entries WHERE user_id = ? AND date >= ? AND date < ?`,
		userID,
		startDate,
		endDateExclusive,
	).Scan(&monthSeconds); err != nil {
		return dto.StatusBarOutput{}, err
	}

	var unpaidTotal float64
	if err := s.db.QueryRow(
		`SELECT COALESCE(SUM(total), 0) FROM invoices WHERE user_id = ? AND status IN ('sent', 'overdue')`,
		userID,
	).Scan(&unpaidTotal); err != nil {
		return dto.StatusBarOutput{}, err
	}

	var uninvoicedTotal float64
	if err := s.db.QueryRow(
		`SELECT COALESCE(SUM((te.duration_seconds / 3600.0) * COALESCE(p.hourly_rate, 0)), 0)
		 FROM time_entries te
		 JOIN projects p ON p.id = te.project_id AND p.user_id = te.user_id
		 WHERE te.user_id = ?
		   AND te.billable = 1
		   AND te.invoiced = 0`,
		userID,
	).Scan(&uninvoicedTotal); err != nil {
		return dto.StatusBarOutput{}, err
	}

	currency := settings.Currency
	if currency == "" {
		currency = "USD"
	}

	return dto.StatusBarOutput{
		MonthSeconds:    monthSeconds,
		UninvoicedTotal: uninvoicedTotal,
		UnpaidTotal:     unpaidTotal,
		Currency:        currency,
	}, nil
}
