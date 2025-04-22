package storage

import (
	"encoding/json"
	"os"

	"github.com/lauren/nannytracker/internal/model"
)

// Storage defines the interface for trip data persistence
type Storage interface {
	SaveTrips(trips []model.Trip) error
	LoadTrips() ([]model.Trip, error)
}

// FileStorage implements Storage using JSON files
type FileStorage struct {
	filePath string
}

func New(filePath string) *FileStorage {
	return &FileStorage{
		filePath: filePath,
	}
}

func (s *FileStorage) SaveTrips(trips []model.Trip) error {
	data, err := json.MarshalIndent(trips, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.filePath, data, 0644)
}

func (s *FileStorage) LoadTrips() ([]model.Trip, error) {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []model.Trip{}, nil
		}
		return nil, err
	}

	var trips []model.Trip
	if err := json.Unmarshal(data, &trips); err != nil {
		return nil, err
	}

	return trips, nil
}
