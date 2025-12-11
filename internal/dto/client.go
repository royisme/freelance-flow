// Package dto defines data transfer objects shared between layers.
package dto

// CreateClientInput represents the input for creating a new client.
// Note: ID is not included as it will be auto-generated.
type CreateClientInput struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	Website       string `json:"website"`
	Avatar        string `json:"avatar"`
	ContactPerson string `json:"contactPerson"`
	Address       string `json:"address"`
	Currency      string `json:"currency"`
	Status        string `json:"status"`
	Notes         string `json:"notes"`
}

// UpdateClientInput represents the input for updating an existing client.
// ID is required to identify the client to update.
type UpdateClientInput struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Website       string `json:"website"`
	Avatar        string `json:"avatar"`
	ContactPerson string `json:"contactPerson"`
	Address       string `json:"address"`
	Currency      string `json:"currency"`
	Status        string `json:"status"`
	Notes         string `json:"notes"`
}

// ClientOutput represents the client data returned from API.
type ClientOutput struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Website       string `json:"website"`
	Avatar        string `json:"avatar"`
	ContactPerson string `json:"contactPerson"`
	Address       string `json:"address"`
	Currency      string `json:"currency"`
	Status        string `json:"status"`
	Notes         string `json:"notes"`
}
