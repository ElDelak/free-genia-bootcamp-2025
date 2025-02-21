package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ElDelak/free-genia-bootcamp-2025/backend_go/internal/models"
)

// WordsQueryParams represents query parameters for the words endpoint
type WordsQueryParams struct {
	GroupID  int64  `form:"group_id"`
	Search   string `form:"search"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=20"`
}

// GetWords returns a list of words with optional filtering and pagination
func (h *Handler) GetWords(c *gin.Context) {
	var params WordsQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate pagination parameters
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 20
	}

	words, total, err := h.repo.GetWords(params.GroupID, params.Search, params.Page, params.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"words": words,
		"pagination": gin.H{
			"total":     total,
			"page":      params.Page,
			"page_size": params.PageSize,
		},
	})
}

// GetWordByID returns a specific word
func (h *Handler) GetWordByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid word id"})
		return
	}

	word, err := h.repo.GetWordByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if word == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "word not found"})
		return
	}

	c.JSON(http.StatusOK, word)
}

// CreateWord creates a new word
func (h *Handler) CreateWord(c *gin.Context) {
	var word models.Word
	if err := c.ShouldBindJSON(&word); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.CreateWord(&word); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, word)
}

// UpdateWord updates an existing word
func (h *Handler) UpdateWord(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid word id"})
		return
	}

	var word models.Word
	if err := c.ShouldBindJSON(&word); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	word.ID = id

	if err := h.repo.UpdateWord(&word); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, word)
}

// DeleteWord deletes a word
func (h *Handler) DeleteWord(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid word id"})
		return
	}

	if err := h.repo.DeleteWord(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
