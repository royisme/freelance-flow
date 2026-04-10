package services

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"tally/internal/dto"
)

// ReportService provides aggregated income/hours reports.
type ReportService struct {
	db *sql.DB
}

// NewReportService creates a ReportService instance.
func NewReportService(db *sql.DB) *ReportService {
	return &ReportService{db: db}
}

// Get aggregates income and hours for the given user and filters.
func (s *ReportService) Get(userID int, filter dto.ReportFilter) (dto.ReportOutput, error) {
	rows, totals, err := s.queryTableRows(userID, filter)
	if err != nil {
		return dto.ReportOutput{}, err
	}
	chart, err := s.queryChartSeries(userID, filter)
	if err != nil {
		return dto.ReportOutput{}, err
	}

	return dto.ReportOutput{
		TotalHours:  totals.totalHours,
		TotalIncome: totals.totalIncome,
		Rows:        rows,
		Chart:       chart,
	}, nil
}

type reportTotals struct {
	totalHours  float64
	totalIncome float64
}

func (s *ReportService) buildWhere(userID int, filter dto.ReportFilter) (string, []any) {
	clauses := []string{"te.user_id = ?", "te.billable = 1"}
	args := []any{userID}

	if filter.StartDate != "" {
		clauses = append(clauses, "te.date >= ?")
		args = append(args, filter.StartDate)
	}
	if filter.EndDate != "" {
		clauses = append(clauses, "te.date <= ?")
		args = append(args, filter.EndDate)
	}
	if filter.ClientID > 0 {
		clauses = append(clauses, "c.id = ?")
		args = append(args, filter.ClientID)
	}
	if filter.ProjectID > 0 {
		clauses = append(clauses, "p.id = ?")
		args = append(args, filter.ProjectID)
	}

	return "WHERE " + strings.Join(clauses, " AND "), args
}

func (s *ReportService) queryTableRows(userID int, filter dto.ReportFilter) ([]dto.ReportRow, reportTotals, error) {
	where, args := s.buildWhere(userID, filter)

	// #nosec G202 -- where clause is composed from fixed predicates with parameter binding.
	query := `
SELECT te.date,
       c.id, c.name,
       p.id, p.name,
       SUM(te.duration_seconds) / 3600.0 AS hours,
       SUM(CASE
           WHEN te.billing_mode = 'fixed' THEN COALESCE(te.manual_amount, 0)
           ELSE (te.duration_seconds / 3600.0) * COALESCE(p.hourly_rate, 0)
       END) AS income
FROM time_entries te
JOIN projects p ON te.project_id = p.id
JOIN clients c ON p.client_id = c.id
` + where + `
GROUP BY te.date, c.id, p.id
ORDER BY te.date ASC, c.name ASC, p.name ASC`

	dbRows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, reportTotals{}, fmt.Errorf("failed to query report rows: %w", err)
	}
	defer closeWithLog(dbRows, "closing report rows")

	var out []dto.ReportRow
	var totals reportTotals
	for dbRows.Next() {
		var r dto.ReportRow
		if err := dbRows.Scan(
			&r.Date,
			&r.ClientID, &r.ClientName,
			&r.ProjectID, &r.ProjectName,
			&r.Hours, &r.Income,
		); err != nil {
			log.Println("Error scanning report row:", err)
			continue
		}
		out = append(out, r)
		totals.totalHours += r.Hours
		totals.totalIncome += r.Income
	}
	return out, totals, nil
}

func (s *ReportService) queryChartSeries(userID int, filter dto.ReportFilter) (dto.ReportChartSeries, error) {
	where, args := s.buildWhere(userID, filter)

	// #nosec G202 -- where clause is composed from fixed predicates with parameter binding.
	query := `
SELECT te.date,
       SUM(te.duration_seconds) / 3600.0 AS hours,
       SUM(CASE
           WHEN te.billing_mode = 'fixed' THEN COALESCE(te.manual_amount, 0)
           ELSE (te.duration_seconds / 3600.0) * COALESCE(p.hourly_rate, 0)
       END) AS income
FROM time_entries te
JOIN projects p ON te.project_id = p.id
JOIN clients c ON p.client_id = c.id
` + where + `
GROUP BY te.date
ORDER BY te.date ASC`

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return dto.ReportChartSeries{}, fmt.Errorf("failed to query report chart: %w", err)
	}
	defer closeWithLog(rows, "closing report chart rows")

	var dates []string
	var hours []float64
	var revenue []float64
	for rows.Next() {
		var d string
		var h, inc float64
		if err := rows.Scan(&d, &h, &inc); err != nil {
			log.Println("Error scanning report chart row:", err)
			continue
		}
		dates = append(dates, d)
		hours = append(hours, h)
		revenue = append(revenue, inc)
	}

	return dto.ReportChartSeries{
		Dates:   dates,
		Hours:   hours,
		Revenue: revenue,
	}, nil
}
