package model

import (
	"encoding/json"
	"testing"
)

func TestTripTemplate(t *testing.T) {
	// Test creating a valid trip template
	template := TripTemplate{
		Name:        "Work Commute",
		Origin:      "123 Home St",
		Destination: "456 Work Ave",
		TripType:    "single",
		Notes:       "Regular work commute",
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(template)
	if err != nil {
		t.Errorf("Failed to marshal template: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaledTemplate TripTemplate
	err = json.Unmarshal(jsonData, &unmarshaledTemplate)
	if err != nil {
		t.Errorf("Failed to unmarshal template: %v", err)
	}

	// Verify all fields are preserved
	if unmarshaledTemplate.Name != template.Name {
		t.Errorf("Expected Name to be %s, got %s", template.Name, unmarshaledTemplate.Name)
	}
	if unmarshaledTemplate.Origin != template.Origin {
		t.Errorf("Expected Origin to be %s, got %s", template.Origin, unmarshaledTemplate.Origin)
	}
	if unmarshaledTemplate.Destination != template.Destination {
		t.Errorf("Expected Destination to be %s, got %s", template.Destination, unmarshaledTemplate.Destination)
	}
	if unmarshaledTemplate.TripType != template.TripType {
		t.Errorf("Expected TripType to be %s, got %s", template.TripType, unmarshaledTemplate.TripType)
	}
	if unmarshaledTemplate.Notes != template.Notes {
		t.Errorf("Expected Notes to be %s, got %s", template.Notes, unmarshaledTemplate.Notes)
	}
}

func TestTripTemplateValidation(t *testing.T) {
	tests := []struct {
		name        string
		template    TripTemplate
		expectError bool
	}{
		{
			name: "valid template",
			template: TripTemplate{
				Name:        "Work Commute",
				Origin:      "123 Home St",
				Destination: "456 Work Ave",
				TripType:    "single",
				Notes:       "Regular work commute",
			},
			expectError: false,
		},
		{
			name: "empty name",
			template: TripTemplate{
				Name:        "",
				Origin:      "123 Home St",
				Destination: "456 Work Ave",
				TripType:    "single",
			},
			expectError: true,
		},
		{
			name: "empty origin",
			template: TripTemplate{
				Name:        "Work Commute",
				Origin:      "",
				Destination: "456 Work Ave",
				TripType:    "single",
			},
			expectError: true,
		},
		{
			name: "empty destination",
			template: TripTemplate{
				Name:        "Work Commute",
				Origin:      "123 Home St",
				Destination: "",
				TripType:    "single",
			},
			expectError: true,
		},
		{
			name: "invalid trip type",
			template: TripTemplate{
				Name:        "Work Commute",
				Origin:      "123 Home St",
				Destination: "456 Work Ave",
				TripType:    "invalid",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.template.Validate()
			if (err != nil) != tt.expectError {
				t.Errorf("TripTemplate.Validate() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}
