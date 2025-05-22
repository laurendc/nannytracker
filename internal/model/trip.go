package model

import (
	"fmt"
)

// TripTemplate represents a saved trip template.
type TripTemplate struct {
	Name        string `json:"name"`
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	TripType    string `json:"tripType"`
	Notes       string `json:"notes"`
}

// Validate checks if the trip template is valid.
func (t *TripTemplate) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("template name cannot be empty")
	}
	if t.Origin == "" {
		return fmt.Errorf("origin cannot be empty")
	}
	if t.Destination == "" {
		return fmt.Errorf("destination cannot be empty")
	}
	if t.TripType != "single" && t.TripType != "round" {
		return fmt.Errorf("invalid trip type: must be 'single' or 'round'")
	}
	return nil
}
