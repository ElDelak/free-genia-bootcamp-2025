package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/ElDelak/free-genia-bootcamp-2025/backend_go/internal/models"
	"strconv"
)

// CreateStudyActivityRequest represents the request body for creating a study activity
type CreateStudyActivityRequest struct {
	GroupID int64 `json:"group_id" binding:"required"`
}

// GetStudyActivities returns all study activities
func (h *Handler) GetStudyActivities(c *gin.Context) {
	activities, err := h.repo.GetStudyActivities()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, activities)
}

// GetStudyActivity returns a specific study activity by ID
func (h *Handler) GetStudyActivity(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid activity id"})
		return
	}

	activity, err := h.repo.GetStudyActivity(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if activity == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "activity not found"})
		return
	}

	c.JSON(http.StatusOK, activity)
}

// GetStudyActivitySessions returns all study sessions for a specific activity
func (h *Handler) GetStudyActivitySessions(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid activity id"})
		return
	}

	sessions, err := h.repo.GetStudySessionsByActivityID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sessions)
}

// CreateStudyActivity creates a new study activity
func (h *Handler) CreateStudyActivity(c *gin.Context) {
	var activity models.StudyActivity
	if err := c.ShouldBindJSON(&activity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.CreateStudyActivity(&activity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, activity)
}

// CreateStudySession creates a new study session for an activity
func (h *Handler) CreateStudySession(c *gin.Context) {
	activityID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid activity id"})
		return
	}

	var session models.StudySession
	if err := c.ShouldBindJSON(&session); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	session.StudyActivityID = activityID

	if err := h.repo.CreateStudySession(&session); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, session)
}
