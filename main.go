package main

import (
	"embed"
	"tally/internal/db"
	"tally/internal/services"
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
	userPreferencesService := services.NewUserPreferencesService(dbConn)
	userTaxSettingsService := services.NewUserTaxSettingsService(dbConn)
	userInvoiceSettingsService := services.NewUserInvoiceSettingsService(dbConn)
	reportService := services.NewReportService(dbConn)
	statusBarService := services.NewStatusBarService(dbConn)
	financeService := services.NewFinanceService(dbConn)
	servicesDuration := time.Since(servicesStart)

	app.SetBootTimings(BootTimings{
		ProcessStart:   start,
		DbInitMs:       dbDuration.Milliseconds(),
		ServicesInitMs: servicesDuration.Milliseconds(),
		TotalBeforeUI:  time.Since(start).Milliseconds(),
	})

	// Create application with options
	err := wails.Run(&options.App{
		Title:             "Tally",
		Width:             1200,
		Height:            900,
		MinWidth:          1024,
		MinHeight:         800,
		HideWindowOnClose: true, // 关闭窗口时隐藏而非退出，用户可通过 Cmd+Q 真正退出
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
			userPreferencesService,
			userTaxSettingsService,
			userInvoiceSettingsService,
			reportService,
			statusBarService,
			financeService,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
