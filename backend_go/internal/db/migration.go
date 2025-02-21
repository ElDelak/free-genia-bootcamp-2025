package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/ElDelak/free-genia-bootcamp-2025/backend_go/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

// InitDB initializes the database connection and runs migrations
func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	// Read and execute migration SQL
	migrationSQL, err := os.ReadFile(filepath.Join("internal", "db", "migrations", "001_initial_schema.up.sql"))
	if err != nil {
		return nil, fmt.Errorf("error reading migration file: %v", err)
	}

	if _, err := db.Exec(string(migrationSQL)); err != nil {
		return nil, fmt.Errorf("error executing migration: %v", err)
	}

	// Load seed data
	if err := loadSeedData(db); err != nil {
		return nil, fmt.Errorf("error loading seed data: %v", err)
	}

	log.Println("Database initialized successfully")
	return db, nil
}

// loadSeedData loads initial data from JSON files
func loadSeedData(db *sql.DB) error {
	// Load groups
	groupData, err := os.ReadFile("data/groups.json")
	if err != nil {
		return fmt.Errorf("error reading groups.json: %v", err)
	}

	var groups []models.Group
	if err := json.Unmarshal(groupData, &groups); err != nil {
		return fmt.Errorf("error parsing groups.json: %v", err)
	}

	// Insert groups
	for _, group := range groups {
		_, err := db.Exec(`
			INSERT INTO groups (name, description)
			VALUES (?, ?)
		`, group.Name, group.Description)
		if err != nil {
			return fmt.Errorf("error inserting group: %v", err)
		}
	}

	// Load words
	wordData, err := os.ReadFile("data/words.json")
	if err != nil {
		return fmt.Errorf("error reading words.json: %v", err)
	}

	var words []models.Word
	if err := json.Unmarshal(wordData, &words); err != nil {
		return fmt.Errorf("error parsing words.json: %v", err)
	}

	// Insert words
	for _, word := range words {
		partsJSON, err := json.Marshal(word.Parts)
		if err != nil {
			return fmt.Errorf("error marshaling word parts: %v", err)
		}

		_, err = db.Exec(`
			INSERT INTO words (arabic, romaji, english, parts)
			VALUES (?, ?, ?, ?)
		`, word.Arabic, word.Romaji, word.English, string(partsJSON))
		if err != nil {
			return fmt.Errorf("error inserting word: %v", err)
		}
	}

	log.Println("Seed data loaded successfully")
	return nil
}

// ResetDB drops all tables and reruns migrations
func ResetDB(db *sql.DB) error {
	// Drop all tables
	_, err := db.Exec(`
		DROP TABLE IF EXISTS word_review_items;
		DROP TABLE IF EXISTS study_sessions;
		DROP TABLE IF EXISTS study_activities;
		DROP TABLE IF EXISTS words_groups;
		DROP TABLE IF EXISTS words;
		DROP TABLE IF EXISTS groups;
	`)
	if err != nil {
		return fmt.Errorf("error dropping tables: %v", err)
	}

	// Read and execute migration SQL
	migrationSQL, err := os.ReadFile(filepath.Join("internal", "db", "migrations", "001_initial_schema.up.sql"))
	if err != nil {
		return fmt.Errorf("error reading migration file: %v", err)
	}

	if _, err := db.Exec(string(migrationSQL)); err != nil {
		return fmt.Errorf("error executing migration: %v", err)
	}

	// Load seed data
	if err := loadSeedData(db); err != nil {
		return fmt.Errorf("error loading seed data: %v", err)
	}

	log.Println("Database reset successfully")
	return nil
}
