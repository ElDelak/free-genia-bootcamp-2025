package models

import (
	"encoding/json"
	"time"
)

// Word represents a vocabulary word with translations
type Word struct {
	ID        int64           `json:"id"`
	Arabic    string         `json:"arabic"`
	Romaji    string         `json:"romaji"`
	English   string         `json:"english"`
	Parts     json.RawMessage `json:"parts"` // JSON data for additional word metadata
	CreatedAt time.Time      `json:"created_at"`
	
	// Relations
	Groups          []Group          `json:"groups,omitempty"`
	WordReviewItems []WordReviewItem `json:"word_review_items,omitempty"`
}

// WordGroup represents the many-to-many relationship between words and groups
type WordGroup struct {
	ID        int64     `json:"id"`
	WordID    int64     `json:"word_id"`
	GroupID   int64     `json:"group_id"`
	CreatedAt time.Time `json:"created_at"`
	
	// Relations
	Word  *Word  `json:"word,omitempty"`
	Group *Group `json:"group,omitempty"`
}
