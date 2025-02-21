package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ElDelak/free-genia-bootcamp-2025/backend_go/internal/models"
)

// GetWords retrieves words with filtering and pagination
func (r *SQLiteRepository) GetWords(groupID int64, search string, page, pageSize int) ([]models.Word, int, error) {
	offset := (page - 1) * pageSize

	// First get total count
	countQuery := "SELECT COUNT(*) FROM words w"
	args := []interface{}{}

	// Add filters
	wheres := []string{}
	if groupID > 0 {
		wheres = append(wheres, "w.id IN (SELECT word_id FROM words_groups WHERE group_id = ?)")
		args = append(args, groupID)
	}
	if search != "" {
		wheres = append(wheres, "(w.arabic LIKE ? OR w.romaji LIKE ? OR w.english LIKE ?)")
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern, searchPattern)
	}
	if len(wheres) > 0 {
		countQuery += " WHERE " + strings.Join(wheres, " AND ")
	}

	var totalCount int
	err := r.db.QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting total count: %v", err)
	}

	// Then get the actual data
	query := `
		SELECT w.id, w.arabic, w.romaji, w.english, w.parts, w.created_at
		FROM words w
	`
	if len(wheres) > 0 {
		query += " WHERE " + strings.Join(wheres, " AND ")
	}
	query += " ORDER BY w.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, pageSize, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("error querying words: %v", err)
	}
	defer rows.Close()

	var words []models.Word
	for rows.Next() {
		var word models.Word
		var partsStr string
		err := rows.Scan(
			&word.ID,
			&word.Arabic,
			&word.Romaji,
			&word.English,
			&partsStr,
			&word.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning word row: %v", err)
		}

		// Parse the JSON string into RawMessage
		if partsStr != "" {
			word.Parts = json.RawMessage(partsStr)
		}
		words = append(words, word)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating word rows: %v", err)
	}

	return words, totalCount, nil
}

// GetWordByID retrieves a single word by ID
func (r *SQLiteRepository) GetWordByID(id int64) (*models.Word, error) {
	var word models.Word
	var partsStr string
	err := r.db.QueryRow(`
		SELECT id, arabic, romaji, english, parts, created_at
		FROM words
		WHERE id = ?
	`, id).Scan(
		&word.ID,
		&word.Arabic,
		&word.Romaji,
		&word.English,
		&partsStr,
		&word.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error scanning word: %v", err)
	}

	// Parse the JSON string into RawMessage
	if partsStr != "" {
		word.Parts = json.RawMessage(partsStr)
	}

	return &word, nil
}

// CreateWord creates a new word
func (r *SQLiteRepository) CreateWord(word *models.Word) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	result, err := tx.Exec(`
		INSERT INTO words (arabic, romaji, english, parts)
		VALUES (?, ?, ?, ?)
	`, word.Arabic, word.Romaji, word.English, word.Parts)
	if err != nil {
		return fmt.Errorf("error inserting word: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last insert id: %v", err)
	}
	word.ID = id

	return tx.Commit()
}

// UpdateWord updates an existing word
func (r *SQLiteRepository) UpdateWord(word *models.Word) error {
	result, err := r.db.Exec(`
		UPDATE words
		SET arabic = ?, romaji = ?, english = ?, parts = ?
		WHERE id = ?
	`, word.Arabic, word.Romaji, word.English, word.Parts, word.ID)
	if err != nil {
		return fmt.Errorf("error updating word: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("word not found: %d", word.ID)
	}

	return nil
}

// DeleteWord deletes a word by ID
func (r *SQLiteRepository) DeleteWord(id int64) error {
	result, err := r.db.Exec("DELETE FROM words WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("error deleting word: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %v", err)
	}
	if rows == 0 {
		return fmt.Errorf("word not found: %d", id)
	}

	return nil
}
