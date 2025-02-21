package models

// DashboardStats represents statistics for the dashboard
type DashboardStats struct {
	TotalWords      int     `json:"total_words"`
	TotalGroups     int     `json:"total_groups"`
	TotalSessions   int     `json:"total_sessions"`
	ReviewCount     int     `json:"review_count"`
	CorrectCount    int     `json:"correct_count"`
	AccuracyRate    float64 `json:"accuracy_rate"`
	LastSessionDate string  `json:"last_session_date"`
}
