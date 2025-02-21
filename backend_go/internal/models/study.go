package models

import "time"

// StudyActivity represents a study activity
type StudyActivity struct {
	ID            int64     `json:"id"`
	GroupID       int64     `json:"group_id"`
	ActivityCount int       `json:"activity_count"`
	ReviewCount   int       `json:"review_count"`
	CorrectCount  int       `json:"correct_count"`
	CreatedAt     time.Time `json:"created_at"`
}

// StudySession represents a study session
type StudySession struct {
	ID              int64     `json:"id"`
	StudyActivityID int64     `json:"study_activity_id"`
	GroupID         int64     `json:"group_id"`
	Group           *Group    `json:"group,omitempty"`
	WordsReviewed   int       `json:"words_reviewed"`
	CorrectCount    int       `json:"correct_count"`
	CreatedAt       time.Time `json:"created_at"`
}

// WordReviewItem represents a word review in a study session
type WordReviewItem struct {
	ID             int64     `json:"id"`
	WordID         int64     `json:"word_id"`
	StudySessionID int64     `json:"study_session_id"`
	IsCorrect      bool      `json:"is_correct"`
	CreatedAt      time.Time `json:"created_at"`
}
