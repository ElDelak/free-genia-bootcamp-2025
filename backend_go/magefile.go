// +build mage

package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	_ "github.com/mattn/go-sqlite3"
)

const (
	binName     = "genia-api"
	dbPath      = "./database.db"
	migrateCmd  = "migrate"
	gooseCmd    = "goose"
)

// Default target to run when none is specified
var Default = Build

// Build builds the application
func Build() error {
	fmt.Println("Building...")
	return sh.Run("go", "build", "-o", binName)
}

// Clean removes compiled files and database
func Clean() error {
	fmt.Println("Cleaning...")
	os.Remove(binName)
	os.Remove(dbPath)
	return nil
}

// Test runs the test suite
func Test() error {
	fmt.Println("Running tests...")
	return sh.Run("go", "test", "./...")
}

// Lint runs golangci-lint
func Lint() error {
	fmt.Println("Running linter...")
	return sh.Run("golangci-lint", "run")
}

// Dev runs the application in development mode
func Dev() error {
	mg.Deps(InstallDeps)
	fmt.Println("Running in development mode...")
	return sh.RunV("go", "run", "main.go")
}

// InstallDeps installs dependencies
func InstallDeps() error {
	fmt.Println("Installing dependencies...")
	deps := []string{
		"github.com/golang-migrate/migrate/v4/cmd/migrate",
		"github.com/golangci/golangci-lint/cmd/golangci-lint@latest",
	}
	for _, dep := range deps {
		if err := sh.Run("go", "install", dep); err != nil {
			return err
		}
	}
	return nil
}

// MigrateCreate creates a new migration file
func MigrateCreate(name string) error {
	if name == "" {
		return fmt.Errorf("migration name is required")
	}
	
	migrationsDir := filepath.Join("internal", "db", "migrations")
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102150405")
	upFile := filepath.Join(migrationsDir, fmt.Sprintf("%s_%s.up.sql", timestamp, name))
	downFile := filepath.Join(migrationsDir, fmt.Sprintf("%s_%s.down.sql", timestamp, name))

	if err := os.WriteFile(upFile, []byte("-- Migration Up"), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(downFile, []byte("-- Migration Down"), 0644); err != nil {
		return err
	}

	fmt.Printf("Created migration files:\n%s\n%s\n", upFile, downFile)
	return nil
}

// MigrateUp runs all pending migrations
func MigrateUp() error {
	fmt.Println("Running migrations...")
	
	// Create database if it doesn't exist
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		file, err := os.Create(dbPath)
		if err != nil {
			return fmt.Errorf("error creating database file: %v", err)
		}
		file.Close()
	}

	// Open database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	// Run migration
	migrationsPath := filepath.Join("internal", "db", "migrations")
	if err := os.MkdirAll(migrationsPath, 0755); err != nil {
		return fmt.Errorf("error creating migrations directory: %v", err)
	}
	
	return sh.RunV(migrateCmd, 
		"-path", migrationsPath,
		"-database", fmt.Sprintf("sqlite3://%s", dbPath),
		"up")
}

// MigrateDown rolls back all migrations
func MigrateDown() error {
	fmt.Println("Rolling back migrations...")
	migrationsPath := filepath.Join("internal", "db", "migrations")
	return sh.RunV(migrateCmd,
		"-path", migrationsPath,
		"-database", fmt.Sprintf("sqlite3://%s", dbPath),
		"down")
}

// MigrateReset resets the database by running down and then up
func MigrateReset() error {
	mg.SerialDeps(MigrateDown, MigrateUp)
	return nil
}

// LoadData loads initial data from JSON files
func LoadData() error {
	fmt.Println("Loading initial data...")
	dataDir := "data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return err
	}

	// Example word data
	wordData := `[
		{
			"arabic": "مرحبا",
			"romaji": "marhaba",
			"english": "hello",
			"parts": {
				"type": "greeting",
				"formality": "neutral"
			}
		}
	]`

	// Example group data
	groupData := `[
		{
			"name": "Basic Vocabulary",
			"description": "Essential words for beginners"
		}
	]`

	if err := os.WriteFile(filepath.Join(dataDir, "words.json"), []byte(wordData), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dataDir, "groups.json"), []byte(groupData), 0644); err != nil {
		return err
	}

	fmt.Println("Created example JSON data files in", dataDir)
	return nil
}

// Setup sets up the development environment
func Setup() error {
	mg.SerialDeps(InstallDeps, Clean, MigrateUp, LoadData)
	return nil
}

// Run runs the application
func Run() error {
	mg.Deps(Build)
	return sh.RunV("./"+binName)
}
