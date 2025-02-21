package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ElDelak/free-genia-bootcamp-2025/backend_go/internal/models"
)

// GetWordReviewItems returns all review items for a study session
func (h *Handler) GetWordReviewItems(c *gin.Context) {
	sessionID, err := strconv.ParseInt(c.Param("session_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}

	reviews, err := h.repo.GetWordReviewItems(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

// CreateWordReviewItem creates a new word review item
func (h *Handler) CreateWordReviewItem(c *gin.Context) {
	var review models.WordReviewItem
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.CreateWordReviewItem(&review); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, review)
}
