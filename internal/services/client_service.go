package services

import (
	"database/sql"
	"freelance-flow/internal/dto"
	"freelance-flow/internal/mapper"
	"freelance-flow/internal/models"
	"log"
)

// ClientService handles all client-related operations.
type ClientService struct {
	db *sql.DB
}

// NewClientService creates a new ClientService instance.
func NewClientService(db *sql.DB) *ClientService {
	return &ClientService{db: db}
}

// List returns all clients for a specific user as DTOs.
func (s *ClientService) List(userID int) []dto.ClientOutput {
	rows, err := s.db.Query("SELECT id, name, email, website, avatar, contact_person, address, currency, status, notes FROM clients WHERE user_id = ?", userID)
	if err != nil {
		log.Println("Error querying clients:", err)
		return []dto.ClientOutput{}
	}
	defer closeWithLog(rows, "closing client rows")

	var clients []models.Client
	for rows.Next() {
		var c models.Client
		err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.Website, &c.Avatar, &c.ContactPerson, &c.Address, &c.Currency, &c.Status, &c.Notes)
		if err != nil {
			log.Println("Error scanning client:", err)
			continue
		}
		clients = append(clients, c)
	}
	return mapper.ToClientOutputList(clients)
}

// Get returns a single client by ID for a specific user.
func (s *ClientService) Get(userID int, id int) (dto.ClientOutput, error) {
	row := s.db.QueryRow("SELECT id, name, email, website, avatar, contact_person, address, currency, status, notes FROM clients WHERE id = ? AND user_id = ?", id, userID)
	var c models.Client
	err := row.Scan(&c.ID, &c.Name, &c.Email, &c.Website, &c.Avatar, &c.ContactPerson, &c.Address, &c.Currency, &c.Status, &c.Notes)
	if err != nil {
		return dto.ClientOutput{}, err
	}
	return mapper.ToClientOutput(c), nil
}

// Create adds a new client for a specific user and returns the created client as DTO.
func (s *ClientService) Create(userID int, input dto.CreateClientInput) dto.ClientOutput {
	entity := mapper.ToClientEntity(input)

	stmt, err := s.db.Prepare("INSERT INTO clients(user_id, name, email, website, avatar, contact_person, address, currency, status, notes) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("Error preparing insert:", err)
		return dto.ClientOutput{}
	}
	defer closeWithLog(stmt, "closing client insert statement")

	res, err := stmt.Exec(userID, entity.Name, entity.Email, entity.Website, entity.Avatar, entity.ContactPerson, entity.Address, entity.Currency, entity.Status, entity.Notes)
	if err != nil {
		log.Println("Error inserting client:", err)
		return dto.ClientOutput{}
	}

	id, _ := res.LastInsertId()
	entity.ID = int(id)
	return mapper.ToClientOutput(entity)
}

// Update modifies an existing client for a specific user and returns the updated client as DTO.
func (s *ClientService) Update(userID int, input dto.UpdateClientInput) dto.ClientOutput {
	stmt, err := s.db.Prepare("UPDATE clients SET name=?, email=?, website=?, avatar=?, contact_person=?, address=?, currency=?, status=?, notes=? WHERE id=? AND user_id=?")
	if err != nil {
		log.Println("Error preparing update:", err)
		return dto.ClientOutput{}
	}
	defer closeWithLog(stmt, "closing client update statement")

	_, err = stmt.Exec(input.Name, input.Email, input.Website, input.Avatar, input.ContactPerson, input.Address, input.Currency, input.Status, input.Notes, input.ID, userID)
	if err != nil {
		log.Println("Error updating client:", err)
		return dto.ClientOutput{}
	}

	// Return the updated client
	output, _ := s.Get(userID, input.ID)
	return output
}

// Delete removes a client by ID for a specific user.
func (s *ClientService) Delete(userID int, id int) {
	_, err := s.db.Exec("DELETE FROM clients WHERE id=? AND user_id=?", id, userID)
	if err != nil {
		log.Println("Error deleting client:", err)
	}
}
