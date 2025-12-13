package services

import (
	"freelance-flow/internal/dto"
	"testing"
	"time"
)

func TestStatusBarService_Get_CurrentMonthAndUnpaid(t *testing.T) {
	db := setupFullTestDB(t)
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("failed to close db: %v", err)
		}
	}()

	authService := NewAuthService(db)
	user := createTestUser(t, authService, "statusbar-user")

	// Force timezone/currency to stable values for this test.
	settingsJSON := `{"timezone":"UTC","currency":"USD"}`
	if _, err := db.Exec("UPDATE users SET settings_json = ? WHERE id = ?", settingsJSON, user.ID); err != nil {
		t.Fatalf("failed to set settings_json: %v", err)
	}

	// Create a client/project for FK constraints.
	clientSvc := NewClientService(db)
	projectSvc := NewProjectService(db)
	client := clientSvc.Create(user.ID, dto.CreateClientInput{Name: "C"})
	project := projectSvc.Create(user.ID, dto.CreateProjectInput{ClientID: client.ID, Name: "P", HourlyRate: 100})

	now := time.Now().UTC()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	prevMonth := monthStart.AddDate(0, -1, 0)

	// Insert two time entries: one current month, one previous month.
	if _, err := db.Exec(
		`INSERT INTO time_entries(user_id, project_id, invoice_id, date, start_time, end_time, duration_seconds, description, billable, invoiced)
		 VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		user.ID, project.ID, nil, monthStart.Format("2006-01-02"), "", "", 3600, "this month", true, false,
	); err != nil {
		t.Fatalf("failed to insert current-month time entry: %v", err)
	}
	if _, err := db.Exec(
		`INSERT INTO time_entries(user_id, project_id, invoice_id, date, start_time, end_time, duration_seconds, description, billable, invoiced)
		 VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		user.ID, project.ID, nil, monthStart.Format("2006-01-02"), "", "", 3600, "invoiced", true, true,
	); err != nil {
		t.Fatalf("failed to insert invoiced time entry: %v", err)
	}
	if _, err := db.Exec(
		`INSERT INTO time_entries(user_id, project_id, invoice_id, date, start_time, end_time, duration_seconds, description, billable, invoiced)
		 VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		user.ID, project.ID, nil, prevMonth.Format("2006-01-02"), "", "", 7200, "prev month", true, false,
	); err != nil {
		t.Fatalf("failed to insert prev-month time entry: %v", err)
	}

	// Insert two invoices: one unpaid, one paid.
	if _, err := db.Exec(
		`INSERT INTO invoices(user_id, client_id, number, issue_date, due_date, subtotal, tax_rate, tax_amount, total, status, items_json)
		 VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		user.ID, client.ID, "INV-test-unpaid", monthStart.Format("2006-01-02"), "", 100, 0, 0, 100, "sent", "[]",
	); err != nil {
		t.Fatalf("failed to insert unpaid invoice: %v", err)
	}
	if _, err := db.Exec(
		`INSERT INTO invoices(user_id, client_id, number, issue_date, due_date, subtotal, tax_rate, tax_amount, total, status, items_json)
		 VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		user.ID, client.ID, "INV-test-draft", monthStart.Format("2006-01-02"), "", 50, 0, 0, 50, "draft", "[]",
	); err != nil {
		t.Fatalf("failed to insert draft invoice: %v", err)
	}
	if _, err := db.Exec(
		`INSERT INTO invoices(user_id, client_id, number, issue_date, due_date, subtotal, tax_rate, tax_amount, total, status, items_json)
		 VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		user.ID, client.ID, "INV-test-paid", monthStart.Format("2006-01-02"), "", 200, 0, 0, 200, "paid", "[]",
	); err != nil {
		t.Fatalf("failed to insert paid invoice: %v", err)
	}

	statusSvc := NewStatusBarService(db)
	out, err := statusSvc.Get(user.ID)
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}

	if out.MonthSeconds != 7200 {
		t.Fatalf("expected MonthSeconds=7200, got %d", out.MonthSeconds)
	}
	if out.UninvoicedTotal != 300 {
		t.Fatalf("expected UninvoicedTotal=300, got %v", out.UninvoicedTotal)
	}
	if out.UnpaidTotal != 100 {
		t.Fatalf("expected UnpaidTotal=100, got %v", out.UnpaidTotal)
	}
	if out.Currency != "USD" {
		t.Fatalf("expected Currency=USD, got %q", out.Currency)
	}
}
