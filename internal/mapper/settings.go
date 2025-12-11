package mapper

import (
	"freelance-flow/internal/dto"
	"freelance-flow/internal/models"
)

// ToUserSettingsDTO converts model settings to dto.
func ToUserSettingsDTO(settings models.UserSettings) dto.UserSettings {
	return dto.UserSettings{
		Currency:               settings.Currency,
		DefaultTaxRate:         settings.DefaultTaxRate,
		Language:               settings.Language,
		Theme:                  settings.Theme,
		DateFormat:             settings.DateFormat,
		Timezone:               settings.Timezone,
		SenderName:             settings.SenderName,
		SenderCompany:          settings.SenderCompany,
		SenderAddress:          settings.SenderAddress,
		SenderPhone:            settings.SenderPhone,
		SenderEmail:            settings.SenderEmail,
		SenderPostalCode:       settings.SenderPostalCode,
		InvoiceTerms:           settings.InvoiceTerms,
		DefaultMessageTemplate: settings.DefaultMessageTemplate,
	}
}

// ToUserSettingsModel converts dto to model settings.
func ToUserSettingsModel(settings dto.UserSettings) models.UserSettings {
	return models.UserSettings{
		Currency:               settings.Currency,
		DefaultTaxRate:         settings.DefaultTaxRate,
		Language:               settings.Language,
		Theme:                  settings.Theme,
		DateFormat:             settings.DateFormat,
		Timezone:               settings.Timezone,
		SenderName:             settings.SenderName,
		SenderCompany:          settings.SenderCompany,
		SenderAddress:          settings.SenderAddress,
		SenderPhone:            settings.SenderPhone,
		SenderEmail:            settings.SenderEmail,
		SenderPostalCode:       settings.SenderPostalCode,
		InvoiceTerms:           settings.InvoiceTerms,
		DefaultMessageTemplate: settings.DefaultMessageTemplate,
	}
}
