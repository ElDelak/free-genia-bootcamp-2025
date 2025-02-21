package repositories

import (
	"fmt"
	"time"

	"github.com/ElDelak/free-genia-bootcamp-2025/backend_go/internal/models"
)

// GetGroups returns all groups
func (r *SQLiteRepository) GetGroups() ([]models.Group, error) {
	query := `
		SELECT id, name, description, created_at
		FROM groups
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying groups: %v", err)
	}
	defer rows.Close()

	var groups []models.Group
	for rows.Next() {
		var group models.Group
		var createdAt string
		err := rows.Scan(
			&group.ID,
			&group.Name,
			&group.Description,
			&createdAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning group: %v", err)
		}

		group.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
		if err != nil {
			return nil, fmt.Errorf("error parsing created_at: %v", err)
		}

		groups = append(groups, group)
	}

	return groups, nil
}

// GetGroupByID returns a specific group
func (r *SQLiteRepository) GetGroupByID(id int64) (*models.Group, error) {
	query := `
		SELECT id, name, description, created_at
		FROM groups
		WHERE id = ?
	`

	var group models.Group
	var createdAt string
	err := r.db.QueryRow(query, id).Scan(
		&group.ID,
		&group.Name,
		&group.Description,
		&createdAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error querying group: %v", err)
	}

	group.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
	if err != nil {
		return nil, fmt.Errorf("error parsing created_at: %v", err)
	}

	return &group, nil
}

// CreateGroup creates a new group
func (r *SQLiteRepository) CreateGroup(group *models.Group) error {
	query := `
		INSERT INTO groups (name, description)
		VALUES (?, ?)
		RETURNING id, created_at
	`

	var createdAt string
	err := r.db.QueryRow(query, group.Name, group.Description).Scan(&group.ID, &createdAt)
	if err != nil {
		return fmt.Errorf("error creating group: %v", err)
	}

	group.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
	if err != nil {
		return fmt.Errorf("error parsing created_at: %v", err)
	}

	return nil
}

// UpdateGroup updates an existing group
func (r *SQLiteRepository) UpdateGroup(group *models.Group) error {
	query := `
		UPDATE groups
		SET name = ?, description = ?
		WHERE id = ?
	`

	result, err := r.db.Exec(query, group.Name, group.Description, group.ID)
	if err != nil {
		return fmt.Errorf("error updating group: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("group not found")
	}

	return nil
}

// DeleteGroup deletes a group
func (r *SQLiteRepository) DeleteGroup(id int64) error {
	query := `
		DELETE FROM groups
		WHERE id = ?
	`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting group: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("group not found")
	}

	return nil
}
