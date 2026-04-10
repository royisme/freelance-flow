package services

import (
	"tally/internal/dto"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimesheetService_CRUD(t *testing.T) {
	db := setupFullTestDB(t)
	defer func() { _ = db.Close() }()

	auth := NewAuthService(db)
	clientSvc := NewClientService(db)
	projectSvc := NewProjectService(db)
	tsSvc := NewTimesheetService(db)

	user := createTestUser(t, auth, "time_user")
	client := clientSvc.Create(user.ID, dto.CreateClientInput{Name: "Client"})
	project := projectSvc.Create(user.ID, dto.CreateProjectInput{
		ClientID: client.ID, Name: "Proj", HourlyRate: 50, Currency: "USD",
	})

	// Create
	date := time.Now().Format("2006-01-02")
	created := tsSvc.Create(user.ID, dto.CreateTimeEntryInput{
		ProjectID:       project.ID,
		Date:            date,
		StartTime:       "09:00",
		EndTime:         "10:00",
		DurationSeconds: 3600,
		Description:     "work",
		Billable:        true,
		Invoiced:        false,
		BillingMode:     "hourly",
	})
	assert.NotZero(t, created.ID)

	// List (all + by project)
	all := tsSvc.List(user.ID, 0)
	assert.Len(t, all, 1)
	byProject := tsSvc.List(user.ID, project.ID)
	assert.Len(t, byProject, 1)

	// Get
	got, err := tsSvc.Get(user.ID, created.ID)
	assert.NoError(t, err)
	assert.Equal(t, "work", got.Description)
	assert.Equal(t, "hourly", got.BillingMode)

	// Update
	fixedAmount := 500.0
	updated := tsSvc.Update(user.ID, dto.UpdateTimeEntryInput{
		ID:              created.ID,
		ProjectID:       project.ID,
		InvoiceID:       0,
		Date:            date,
		StartTime:       "",
		EndTime:         "",
		DurationSeconds: 1800,
		Description:     "work2",
		Billable:        false,
		Invoiced:        false,
		BillingMode:     "fixed",
		ManualAmount:    &fixedAmount,
	})
	assert.Equal(t, created.ID, updated.ID)
	assert.Equal(t, 1800, updated.DurationSeconds)
	assert.False(t, updated.Billable)
	assert.Equal(t, "fixed", updated.BillingMode)
	assert.NotNil(t, updated.ManualAmount)

	// Delete
	tsSvc.Delete(user.ID, created.ID)
	assert.Len(t, tsSvc.List(user.ID, 0), 0)
}
