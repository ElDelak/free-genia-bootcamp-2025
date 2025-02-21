package loader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/ElDelak/free-genia-bootcamp-2025/backend_go/internal/models"
	"github.com/ElDelak/free-genia-bootcamp-2025/backend_go/internal/repositories"
)

// JSONLoader handles loading initial data from JSON files
type JSONLoader struct {
	repo repositories.Repository
}

// NewJSONLoader creates a new JSON loader
func NewJSONLoader(repo repositories.Repository) *JSONLoader {
	return &JSONLoader{repo: repo}
}

// LoadInitialData loads initial data from JSON files
func (l *JSONLoader) LoadInitialData() error {
	// Load groups
	groups, err := l.loadGroups()
	if err != nil {
		return fmt.Errorf("error loading groups: %v", err)
	}

	// Load words
	err = l.loadWords(groups)
	if err != nil {
		return fmt.Errorf("error loading words: %v", err)
	}

	return nil
}

// loadGroups loads groups from JSON file
func (l *JSONLoader) loadGroups() (map[string]int64, error) {
	groupMap := make(map[string]int64)

	// Read groups.json
	data, err := ioutil.ReadFile("data/groups.json")
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("groups.json not found, skipping")
			return groupMap, nil
		}
		return nil, fmt.Errorf("error reading groups.json: %v", err)
	}

	var groups []models.Group
	if err := json.Unmarshal(data, &groups); err != nil {
		return nil, fmt.Errorf("error unmarshaling groups: %v", err)
	}

	// Create groups
	for _, group := range groups {
		if err := l.repo.CreateGroup(&group); err != nil {
			return nil, fmt.Errorf("error creating group: %v", err)
		}
		groupMap[group.Name] = group.ID
	}

	return groupMap, nil
}

// loadWords loads words from JSON file
func (l *JSONLoader) loadWords(groupMap map[string]int64) error {
	// Read words.json
	data, err := ioutil.ReadFile("data/words.json")
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("words.json not found, skipping")
			return nil
		}
		return fmt.Errorf("error reading words.json: %v", err)
	}

	var words []models.Word
	if err := json.Unmarshal(data, &words); err != nil {
		return fmt.Errorf("error unmarshaling words: %v", err)
	}

	// Create words
	for _, word := range words {
		if err := l.repo.CreateWord(&word); err != nil {
			return fmt.Errorf("error creating word: %v", err)
		}
	}

	return nil
}
