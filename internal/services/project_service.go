package services

import (
	"database/sql"
	"freelance-flow/internal/dto"
	"freelance-flow/internal/mapper"
	"freelance-flow/internal/models"
	"log"
	"strings"
)

// ProjectService handles all project-related operations.
type ProjectService struct {
	db *sql.DB
}

// NewProjectService creates a new ProjectService instance.
func NewProjectService(db *sql.DB) *ProjectService {
	return &ProjectService{db: db}
}

// List returns all projects for a specific user as DTOs.
func (s *ProjectService) List(userID int) []dto.ProjectOutput {
	rows, err := s.db.Query("SELECT id, client_id, name, description, hourly_rate, currency, status, deadline, tags FROM projects WHERE user_id = ?", userID)
	if err != nil {
		log.Println("Error querying projects:", err)
		return []dto.ProjectOutput{}
	}
	defer closeWithLog(rows, "closing project list rows")

	var projects []models.Project
	for rows.Next() {
		var p models.Project
		var tagsStr string
		err := rows.Scan(&p.ID, &p.ClientID, &p.Name, &p.Description, &p.HourlyRate, &p.Currency, &p.Status, &p.Deadline, &tagsStr)
		if err != nil {
			log.Println("Error scanning project:", err)
			continue
		}
		if tagsStr != "" {
			p.Tags = strings.Split(tagsStr, ",")
		} else {
			p.Tags = []string{}
		}
		projects = append(projects, p)
	}
	return mapper.ToProjectOutputList(projects)
}

// ListByClient returns all projects for a specific client of a specific user.
func (s *ProjectService) ListByClient(userID int, clientID int) []dto.ProjectOutput {
	rows, err := s.db.Query("SELECT id, client_id, name, description, hourly_rate, currency, status, deadline, tags FROM projects WHERE client_id = ? AND user_id = ?", clientID, userID)
	if err != nil {
		log.Println("Error querying projects by client:", err)
		return []dto.ProjectOutput{}
	}
	defer closeWithLog(rows, "closing project by client rows")

	var projects []models.Project
	for rows.Next() {
		var p models.Project
		var tagsStr string
		err := rows.Scan(&p.ID, &p.ClientID, &p.Name, &p.Description, &p.HourlyRate, &p.Currency, &p.Status, &p.Deadline, &tagsStr)
		if err != nil {
			log.Println("Error scanning project:", err)
			continue
		}
		if tagsStr != "" {
			p.Tags = strings.Split(tagsStr, ",")
		} else {
			p.Tags = []string{}
		}
		projects = append(projects, p)
	}
	return mapper.ToProjectOutputList(projects)
}

// Get returns a single project by ID for a specific user.
func (s *ProjectService) Get(userID int, id int) (dto.ProjectOutput, error) {
	row := s.db.QueryRow("SELECT id, client_id, name, description, hourly_rate, currency, status, deadline, tags FROM projects WHERE id = ? AND user_id = ?", id, userID)
	var p models.Project
	var tagsStr string
	err := row.Scan(&p.ID, &p.ClientID, &p.Name, &p.Description, &p.HourlyRate, &p.Currency, &p.Status, &p.Deadline, &tagsStr)
	if err != nil {
		return dto.ProjectOutput{}, err
	}
	if tagsStr != "" {
		p.Tags = strings.Split(tagsStr, ",")
	} else {
		p.Tags = []string{}
	}
	return mapper.ToProjectOutput(p), nil
}

// Create adds a new project for a specific user and returns the created project as DTO.
func (s *ProjectService) Create(userID int, input dto.CreateProjectInput) dto.ProjectOutput {
	entity := mapper.ToProjectEntity(input)
	tagsStr := strings.Join(entity.Tags, ",")

	stmt, err := s.db.Prepare("INSERT INTO projects(user_id, client_id, name, description, hourly_rate, currency, status, deadline, tags) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("Error preparing project insert:", err)
		return dto.ProjectOutput{}
	}
	defer closeWithLog(stmt, "closing project insert statement")

	res, err := stmt.Exec(userID, entity.ClientID, entity.Name, entity.Description, entity.HourlyRate, entity.Currency, entity.Status, entity.Deadline, tagsStr)
	if err != nil {
		log.Println("Error inserting project:", err)
		return dto.ProjectOutput{}
	}

	id, _ := res.LastInsertId()
	entity.ID = int(id)
	return mapper.ToProjectOutput(entity)
}

// Update modifies an existing project for a specific user and returns the updated project as DTO.
func (s *ProjectService) Update(userID int, input dto.UpdateProjectInput) dto.ProjectOutput {
	tagsStr := strings.Join(input.Tags, ",")

	stmt, err := s.db.Prepare("UPDATE projects SET client_id=?, name=?, description=?, hourly_rate=?, currency=?, status=?, deadline=?, tags=? WHERE id=? AND user_id=?")
	if err != nil {
		log.Println("Error preparing project update:", err)
		return dto.ProjectOutput{}
	}
	defer closeWithLog(stmt, "closing project update statement")

	_, err = stmt.Exec(input.ClientID, input.Name, input.Description, input.HourlyRate, input.Currency, input.Status, input.Deadline, tagsStr, input.ID, userID)
	if err != nil {
		log.Println("Error updating project:", err)
		return dto.ProjectOutput{}
	}

	output, _ := s.Get(userID, input.ID)
	return output
}

// Delete removes a project by ID for a specific user.
func (s *ProjectService) Delete(userID int, id int) {
	_, err := s.db.Exec("DELETE FROM projects WHERE id=? AND user_id=?", id, userID)
	if err != nil {
		log.Println("Error deleting project:", err)
	}
}
