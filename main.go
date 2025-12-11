package main

import (
	"embed"
	"freelance-flow/internal/db"
	"freelance-flow/internal/services"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Initialize Database
	dbConn := db.Init()

	// Initialize Services
	authService := services.NewAuthService(dbConn)
	clientService := services.NewClientService(dbConn)
	projectService := services.NewProjectService(dbConn)
	timesheetService := services.NewTimesheetService(dbConn)
	invoiceService := services.NewInvoiceService(dbConn)
	updateService := services.NewUpdateService()
	settingsService := services.NewSettingsService(dbConn)
	invoiceEmailSettingsService := services.NewInvoiceEmailSettingsService(dbConn)

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "FreelanceFlow",
		Width:  1024,
		Height: 768,
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
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
