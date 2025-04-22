package model

import (
	"testing"
)

func TestCalculateTotalMiles(t *testing.T) {
	trips := []Trip{
		{Origin: "A", Destination: "B", Miles: 10.0},
		{Origin: "C", Destination: "D", Miles: 15.0},
		{Origin: "E", Destination: "F", Miles: 5.0},
	}

	totalMiles := CalculateTotalMiles(trips)
	if totalMiles != 30.0 {
		t.Errorf("Expected total miles to be 30.0, got %.2f", totalMiles)
	}
}

func TestCalculateReimbursement(t *testing.T) {
	trips := []Trip{
		{Origin: "A", Destination: "B", Miles: 10.0},
		{Origin: "C", Destination: "D", Miles: 15.0},
		{Origin: "E", Destination: "F", Miles: 5.0},
	}

	ratePerMile := 0.70
	totalReimbursement := CalculateReimbursement(trips, ratePerMile)
	expectedReimbursement := 30.0 * ratePerMile
	if totalReimbursement != expectedReimbursement {
		t.Errorf("Expected total reimbursement to be $%.2f, got $%.2f", expectedReimbursement, totalReimbursement)
	}
}

func TestTripValidation(t *testing.T) {
	tests := []struct {
		name    string
		trip    Trip
		wantErr bool
	}{
		{
			name:    "valid trip",
			trip:    Trip{Origin: "Home", Destination: "Work", Miles: 5.0},
			wantErr: false,
		},
		{
			name:    "empty origin",
			trip:    Trip{Origin: "", Destination: "Work", Miles: 5.0},
			wantErr: true,
		},
		{
			name:    "empty destination",
			trip:    Trip{Origin: "Home", Destination: "", Miles: 5.0},
			wantErr: true,
		},
		{
			name:    "negative miles",
			trip:    Trip{Origin: "Home", Destination: "Work", Miles: -5.0},
			wantErr: true,
		},
		{
			name:    "zero miles",
			trip:    Trip{Origin: "Home", Destination: "Work", Miles: 0.0},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.trip.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Trip.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
