package repositories

import (
	"fmt"
	"time"

	"github.com/ElDelak/free-genia-bootcamp-2025/backend_go/internal/models"
)

// GetWordReviewItems returns all word review items for a study session
func (r *SQLiteRepository) GetWordReviewItems(sessionID int64) ([]models.WordReviewItem, error) {
	query := `
		SELECT id, word_id, study_session_id, is_correct, created_at
		FROM word_review_items
		WHERE study_session_id = ?
	`

	rows, err := r.db.Query(query, sessionID)
	if err != nil {
		return nil, fmt.Errorf("error querying word review items: %v", err)
	}
	defer rows.Close()

	var reviews []models.WordReviewItem
	for rows.Next() {
		var review models.WordReviewItem
		var createdAt string
		err := rows.Scan(
			&review.ID,
			&review.WordID,
			&review.StudySessionID,
			&review.IsCorrect,
			&createdAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning word review item: %v", err)
		}

		review.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
		if err != nil {
			return nil, fmt.Errorf("error parsing created_at: %v", err)
		}

		reviews = append(reviews, review)
	}

	return reviews, nil
}

// CreateWordReviewItem creates a new word review item
func (r *SQLiteRepository) CreateWordReviewItem(review *models.WordReviewItem) error {
	query := `
		INSERT INTO word_review_items (word_id, study_session_id, is_correct)
		VALUES (?, ?, ?)
		RETURNING id, created_at
	`

	var createdAt string
	err := r.db.QueryRow(query, review.WordID, review.StudySessionID, review.IsCorrect).Scan(&review.ID, &createdAt)
	if err != nil {
		return fmt.Errorf("error creating word review item: %v", err)
	}

	review.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
	if err != nil {
		return fmt.Errorf("error parsing created_at: %v", err)
	}

	return nil
}
