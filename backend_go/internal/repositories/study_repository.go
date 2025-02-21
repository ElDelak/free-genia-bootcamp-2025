package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ElDelak/free-genia-bootcamp-2025/backend_go/internal/models"
)

// GetLastStudySession retrieves the most recent study session
func (r *SQLiteRepository) GetLastStudySession() (*models.StudySession, error) {
	var session models.StudySession
	var groupName string

	err := r.db.QueryRow(`
		SELECT 
			s.id, s.group_id, s.created_at,
			g.name,
			COUNT(w.id) as words_reviewed,
			SUM(CASE WHEN w.correct THEN 1 ELSE 0 END) as correct_count
		FROM study_sessions s
		JOIN groups g ON s.group_id = g.id
		LEFT JOIN word_review_items w ON s.id = w.study_session_id
		GROUP BY s.id
		ORDER BY s.created_at DESC
		LIMIT 1
	`).Scan(
		&session.ID,
		&session.GroupID,
		&session.CreatedAt,
		&groupName,
		&session.WordsReviewed,
		&session.CorrectCount,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error querying last study session: %v", err)
	}

	session.Group = &models.Group{
		ID:   session.GroupID,
		Name: groupName,
	}

	return &session, nil
}

// GetStudyActivities returns all study activities
func (r *SQLiteRepository) GetStudyActivities() ([]models.StudyActivity, error) {
	query := `
		SELECT 
			sa.id, 
			sa.group_id,
			COUNT(DISTINCT ss.id) as activity_count,
			COUNT(wri.id) as review_count,
			COUNT(CASE WHEN wri.is_correct THEN 1 END) as correct_count,
			sa.created_at
		FROM study_activities sa
		LEFT JOIN study_sessions ss ON ss.study_activity_id = sa.id
		LEFT JOIN word_review_items wri ON wri.study_session_id = ss.id
		GROUP BY sa.id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying study activities: %v", err)
	}
	defer rows.Close()

	var activities []models.StudyActivity
	for rows.Next() {
		var activity models.StudyActivity
		var createdAt string
		err := rows.Scan(
			&activity.ID,
			&activity.GroupID,
			&activity.ActivityCount,
			&activity.ReviewCount,
			&activity.CorrectCount,
			&createdAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning study activity: %v", err)
		}

		activity.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
		if err != nil {
			return nil, fmt.Errorf("error parsing created_at: %v", err)
		}

		activities = append(activities, activity)
	}

	return activities, nil
}

// GetStudyActivity returns a specific study activity
func (r *SQLiteRepository) GetStudyActivity(id int64) (*models.StudyActivity, error) {
	query := `
		SELECT 
			sa.id, 
			sa.group_id,
			COUNT(DISTINCT ss.id) as activity_count,
			COUNT(wri.id) as review_count,
			COUNT(CASE WHEN wri.is_correct THEN 1 END) as correct_count,
			sa.created_at
		FROM study_activities sa
		LEFT JOIN study_sessions ss ON ss.study_activity_id = sa.id
		LEFT JOIN word_review_items wri ON wri.study_session_id = ss.id
		WHERE sa.id = ?
		GROUP BY sa.id
	`

	var activity models.StudyActivity
	var createdAt string
	err := r.db.QueryRow(query, id).Scan(
		&activity.ID,
		&activity.GroupID,
		&activity.ActivityCount,
		&activity.ReviewCount,
		&activity.CorrectCount,
		&createdAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error querying study activity: %v", err)
	}

	activity.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
	if err != nil {
		return nil, fmt.Errorf("error parsing created_at: %v", err)
	}

	return &activity, nil
}

// CreateStudyActivity creates a new study activity
func (r *SQLiteRepository) CreateStudyActivity(activity *models.StudyActivity) error {
	query := `
		INSERT INTO study_activities (group_id)
		VALUES (?)
		RETURNING id, created_at
	`

	var createdAt string
	err := r.db.QueryRow(query, activity.GroupID).Scan(&activity.ID, &createdAt)
	if err != nil {
		return fmt.Errorf("error creating study activity: %v", err)
	}

	activity.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
	if err != nil {
		return fmt.Errorf("error parsing created_at: %v", err)
	}

	return nil
}

// GetStudySessionsByActivityID returns all study sessions for an activity
func (r *SQLiteRepository) GetStudySessionsByActivityID(activityID int64) ([]models.StudySession, error) {
	query := `
		SELECT id, study_activity_id, group_id, created_at
		FROM study_sessions
		WHERE study_activity_id = ?
	`

	rows, err := r.db.Query(query, activityID)
	if err != nil {
		return nil, fmt.Errorf("error querying study sessions: %v", err)
	}
	defer rows.Close()

	var sessions []models.StudySession
	for rows.Next() {
		var session models.StudySession
		var createdAt string
		err := rows.Scan(
			&session.ID,
			&session.StudyActivityID,
			&session.GroupID,
			&createdAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning study session: %v", err)
		}

		session.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
		if err != nil {
			return nil, fmt.Errorf("error parsing created_at: %v", err)
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}

// CreateStudySession creates a new study session
func (r *SQLiteRepository) CreateStudySession(session *models.StudySession) error {
	query := `
		INSERT INTO study_sessions (study_activity_id, group_id)
		VALUES (?, ?)
		RETURNING id, created_at
	`

	var createdAt string
	err := r.db.QueryRow(query, session.StudyActivityID, session.GroupID).Scan(&session.ID, &createdAt)
	if err != nil {
		return fmt.Errorf("error creating study session: %v", err)
	}

	session.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
	if err != nil {
		return fmt.Errorf("error parsing created_at: %v", err)
	}

	return nil
}

// GetStudyProgress returns study progress for the last n days
func (r *SQLiteRepository) GetStudyProgress(days int) ([]models.StudyActivity, error) {
	query := `
		SELECT 
			sa.id, 
			sa.group_id,
			COUNT(DISTINCT ss.id) as activity_count,
			COUNT(wri.id) as review_count,
			COUNT(CASE WHEN wri.is_correct THEN 1 END) as correct_count,
			sa.created_at
		FROM study_activities sa
		LEFT JOIN study_sessions ss ON ss.study_activity_id = sa.id
		LEFT JOIN word_review_items wri ON wri.study_session_id = ss.id
		WHERE sa.created_at >= datetime('now', ?)
		GROUP BY sa.id
		ORDER BY sa.created_at DESC
	`

	rows, err := r.db.Query(query, fmt.Sprintf("-%d days", days))
	if err != nil {
		return nil, fmt.Errorf("error querying study progress: %v", err)
	}
	defer rows.Close()

	var activities []models.StudyActivity
	for rows.Next() {
		var activity models.StudyActivity
		var createdAt string
		err := rows.Scan(
			&activity.ID,
			&activity.GroupID,
			&activity.ActivityCount,
			&activity.ReviewCount,
			&activity.CorrectCount,
			&createdAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning study activity: %v", err)
		}

		activity.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAt)
		if err != nil {
			return nil, fmt.Errorf("error parsing created_at: %v", err)
		}

		activities = append(activities, activity)
	}

	return activities, nil
}

// GetQuickStats returns quick statistics for the dashboard
func (r *SQLiteRepository) GetQuickStats() (*models.DashboardStats, error) {
	query := `
		WITH stats AS (
			SELECT 
				COUNT(DISTINCT w.id) as total_words,
				COUNT(DISTINCT g.id) as total_groups,
				COUNT(DISTINCT ss.id) as total_sessions,
				COUNT(wri.id) as review_count,
				COUNT(CASE WHEN wri.is_correct THEN 1 END) as correct_count,
				MAX(ss.created_at) as last_session_date
			FROM words w
			CROSS JOIN groups g
			LEFT JOIN study_sessions ss ON ss.group_id = g.id
			LEFT JOIN word_review_items wri ON wri.study_session_id = ss.id
		)
		SELECT 
			total_words,
			total_groups,
			total_sessions,
			review_count,
			correct_count,
			CAST(CASE 
				WHEN review_count > 0 
				THEN ROUND(CAST(correct_count AS FLOAT) / review_count * 100, 2)
				ELSE 0 
			END AS FLOAT) as accuracy_rate,
			last_session_date
		FROM stats
	`

	stats := &models.DashboardStats{}
	var lastSessionDate sql.NullString

	err := r.db.QueryRow(query).Scan(
		&stats.TotalWords,
		&stats.TotalGroups,
		&stats.TotalSessions,
		&stats.ReviewCount,
		&stats.CorrectCount,
		&stats.AccuracyRate,
		&lastSessionDate,
	)
	if err != nil {
		return nil, fmt.Errorf("error querying quick stats: %v", err)
	}

	if lastSessionDate.Valid {
		stats.LastSessionDate = lastSessionDate.String
	}

	return stats, nil
}
