package main

import (
	"embed"
	"freelance-flow/internal/db"
	"freelance-flow/internal/services"
	"time"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

var assets embed.FS

func main() {
	start := time.Now()
	// Create an instance of the app structure
	app := NewApp()

	// Initialize Database
	dbStart := time.Now()
	dbConn := db.Init()
	dbDuration := time.Since(dbStart)

	// Initialize Services
	servicesStart := time.Now()
	authService := services.NewAuthService(dbConn)
	clientService := services.NewClientService(dbConn)
	projectService := services.NewProjectService(dbConn)
	timesheetService := services.NewTimesheetService(dbConn)
	invoiceService := services.NewInvoiceService(dbConn)
	updateService := services.NewUpdateService()
	settingsService := services.NewSettingsService(dbConn)
	invoiceEmailSettingsService := services.NewInvoiceEmailSettingsService(dbConn)
	reportService := services.NewReportService(dbConn)
	statusBarService := services.NewStatusBarService(dbConn)
	servicesDuration := time.Since(servicesStart)

	app.SetBootTimings(BootTimings{
		ProcessStart:    start,
		DbInitMs:        dbDuration.Milliseconds(),
		ServicesInitMs:  servicesDuration.Milliseconds(),
		TotalBeforeUI:   time.Since(start).Milliseconds(),
	})

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "FreelanceFlow",
		Width:     1200,
		Height:    900,
		MinWidth:  1024,
		MinHeight: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
			authService,
			clientService,
			projectService,
			timesheetService,
			invoiceService,
			updateService,
			settingsService,
			invoiceEmailSettingsService,
			reportService,
			statusBarService,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
