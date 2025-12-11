// Package mapper converts between models and DTOs.
package mapper

import (
	"freelance-flow/internal/dto"
	"freelance-flow/internal/models"
)

// ToClientOutput converts a Client entity to ClientOutput DTO.
func ToClientOutput(e models.Client) dto.ClientOutput {
	return dto.ClientOutput{
		ID:            e.ID,
		Name:          e.Name,
		Email:         e.Email,
		Website:       e.Website,
		Avatar:        e.Avatar,
		ContactPerson: e.ContactPerson,
		Address:       e.Address,
		Currency:      e.Currency,
		Status:        e.Status,
		Notes:         e.Notes,
	}
}

// ToClientOutputList converts a slice of Client entities to ClientOutput DTOs.
func ToClientOutputList(entities []models.Client) []dto.ClientOutput {
	if entities == nil {
		return []dto.ClientOutput{}
	}
	result := make([]dto.ClientOutput, len(entities))
	for i, e := range entities {
		result[i] = ToClientOutput(e)
	}
	return result
}

// ToClientEntity converts CreateClientInput DTO to Client entity.
func ToClientEntity(input dto.CreateClientInput) models.Client {
	return models.Client{
		Name:          input.Name,
		Email:         input.Email,
		Website:       input.Website,
		Avatar:        input.Avatar,
		ContactPerson: input.ContactPerson,
		Address:       input.Address,
		Currency:      input.Currency,
		Status:        input.Status,
		Notes:         input.Notes,
	}
}

// ApplyClientUpdate applies UpdateClientInput to an existing Client entity.
func ApplyClientUpdate(e *models.Client, input dto.UpdateClientInput) {
	e.Name = input.Name
	e.Email = input.Email
	e.Website = input.Website
	e.Avatar = input.Avatar
	e.ContactPerson = input.ContactPerson
	e.Address = input.Address
	e.Currency = input.Currency
	e.Status = input.Status
	e.Notes = input.Notes
}
