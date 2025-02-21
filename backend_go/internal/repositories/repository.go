package repositories

import (
	"database/sql"
	"github.com/ElDelak/free-genia-bootcamp-2025/backend_go/internal/models"
)

// Repository defines all database operations
type Repository interface {
	// Word operations
	GetWords(groupID int64, search string, page, pageSize int) ([]models.Word, int, error)
	GetWordByID(id int64) (*models.Word, error)
	CreateWord(word *models.Word) error
	UpdateWord(word *models.Word) error
	DeleteWord(id int64) error

	// Group operations
	GetGroups() ([]models.Group, error)
	GetGroupByID(id int64) (*models.Group, error)
	CreateGroup(group *models.Group) error
	UpdateGroup(group *models.Group) error
	DeleteGroup(id int64) error

	// Study session operations
	GetLastStudySession() (*models.StudySession, error)
	GetStudySessionsByActivityID(activityID int64) ([]models.StudySession, error)
	CreateStudySession(session *models.StudySession) error

	// Study activity operations
	GetStudyActivities() ([]models.StudyActivity, error)
	GetStudyActivity(id int64) (*models.StudyActivity, error)
	CreateStudyActivity(activity *models.StudyActivity) error
	GetStudyProgress(days int) ([]models.StudyActivity, error)

	// Word review operations
	GetWordReviewItems(sessionID int64) ([]models.WordReviewItem, error)
	CreateWordReviewItem(review *models.WordReviewItem) error
	GetQuickStats() (*models.DashboardStats, error)
}

// SQLiteRepository implements Repository interface
type SQLiteRepository struct {
	db *sql.DB
}

// NewSQLiteRepository creates a new SQLite repository
func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{db: db}
}
