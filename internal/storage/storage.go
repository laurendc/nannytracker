package storage

import (
	"encoding/json"
	"os"

	"github.com/lauren/nannytracker/internal/model"
)

// Storage defines the interface for trip data persistence
type Storage interface {
	SaveData(data *model.StorageData) error
	LoadData() (*model.StorageData, error)
}

// FileStorage implements Storage using a JSON file
type FileStorage struct {
	filePath string
}

// New creates a new FileStorage instance
func New(filePath string) *FileStorage {
	return &FileStorage{
		filePath: filePath,
	}
}

// SaveData saves the complete data structure to the file
func (s *FileStorage) SaveData(data *model.StorageData) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filePath, jsonData, 0600)
}

// LoadData loads the complete data structure from the file
func (s *FileStorage) LoadData() (*model.StorageData, error) {
	data := &model.StorageData{
		Trips:           make([]model.Trip, 0),
		WeeklySummaries: make([]model.WeeklySummary, 0),
		TripTemplates:   make([]model.TripTemplate, 0),
	}

	fileData, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return data, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(fileData, data); err != nil {
		return nil, err
	}

	return data, nil
}
