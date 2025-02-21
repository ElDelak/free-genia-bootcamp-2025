package main

import (
	"log"

	"github.com/ElDelak/free-genia-bootcamp-2025/backend_go/internal/db"
	"github.com/ElDelak/free-genia-bootcamp-2025/backend_go/internal/handlers"
	"github.com/ElDelak/free-genia-bootcamp-2025/backend_go/internal/loader"
	"github.com/ElDelak/free-genia-bootcamp-2025/backend_go/internal/repositories"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	database, err := db.InitDB("./database.db")
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	// Create repository
	repo := repositories.NewSQLiteRepository(database)

	// Initialize handlers
	handler := handlers.NewHandler(repo)

	// Initialize JSON loader
	jsonLoader := loader.NewJSONLoader(repo)

	// Load initial data
	if err := jsonLoader.LoadInitialData(); err != nil {
		log.Printf("Error loading initial data: %v", err)
	}

	// Initialize router
	router := gin.Default()

	// Enable CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// API routes
	api := router.Group("/api")
	{
		// Dashboard endpoints
		api.GET("/dashboard/last-session", handler.GetLastStudySession)
		api.GET("/dashboard/study-progress", handler.GetStudyProgress)
		api.GET("/dashboard/quick-stats", handler.GetQuickStats)

		// Word endpoints
		api.GET("/words", handler.GetWords)
		api.POST("/words", handler.CreateWord)
		api.GET("/words/:id", handler.GetWordByID)
		api.PUT("/words/:id", handler.UpdateWord)
		api.DELETE("/words/:id", handler.DeleteWord)

		// Group endpoints
		api.GET("/groups", handler.GetGroups)
		api.POST("/groups", handler.CreateGroup)
		api.GET("/groups/:id", handler.GetGroupByID)
		api.PUT("/groups/:id", handler.UpdateGroup)
		api.DELETE("/groups/:id", handler.DeleteGroup)

		// Study activity endpoints
		api.GET("/study-activities", handler.GetStudyActivities)
		api.POST("/study-activities", handler.CreateStudyActivity)
		api.GET("/study-activities/:id", handler.GetStudyActivity)
		api.GET("/study-activities/:id/study-sessions", handler.GetStudyActivitySessions)
		api.POST("/study-activities/:id/study-sessions", handler.CreateStudySession)

		// Word review endpoints
		api.GET("/reviews/session/:session_id", handler.GetWordReviewItems)
		api.POST("/reviews", handler.CreateWordReviewItem)
	}

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
