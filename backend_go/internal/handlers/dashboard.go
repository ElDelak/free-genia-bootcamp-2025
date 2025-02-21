package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ElDelak/free-genia-bootcamp-2025/backend_go/internal/repositories"
)

type Handler struct {
	repo repositories.Repository
}

func NewHandler(repo repositories.Repository) *Handler {
	return &Handler{repo: repo}
}

// DashboardStats represents the quick stats response
type DashboardStats struct {
	TotalWords      int64 `json:"total_words"`
	WordsLearned    int64 `json:"words_learned"`
	StudyStreak     int   `json:"study_streak"`
	TodayReviews    int64 `json:"today_reviews"`
}

// GetLastStudySession returns the most recent study session
func (h *Handler) GetLastStudySession(c *gin.Context) {
	session, err := h.repo.GetLastStudySession()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if session == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no study sessions found"})
		return
	}
	c.JSON(http.StatusOK, session)
}

// GetStudyProgress returns study progress for the specified number of days
func (h *Handler) GetStudyProgress(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid days parameter"})
		return
	}

	progress, err := h.repo.GetStudyProgress(days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, progress)
}

// GetQuickStats returns dashboard statistics
func (h *Handler) GetQuickStats(c *gin.Context) {
	stats, err := h.repo.GetQuickStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}
