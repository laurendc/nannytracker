package model

import (
	"sort"
	"testing"
	"time"
)

func TestCalculateTotalMiles(t *testing.T) {
	trips := []Trip{
		{Origin: "A", Destination: "B", Miles: 10.0, Date: "2024-03-20"},
		{Origin: "C", Destination: "D", Miles: 15.0, Date: "2024-03-21"},
		{Origin: "E", Destination: "F", Miles: 5.0, Date: "2024-03-22"},
	}

	totalMiles := CalculateTotalMiles(trips)
	if totalMiles != 30.0 {
		t.Errorf("Expected total miles to be 30.0, got %.2f", totalMiles)
	}
}

func TestCalculateReimbursement(t *testing.T) {
	trips := []Trip{
		{Origin: "A", Destination: "B", Miles: 10.0, Date: "2024-03-20"},
		{Origin: "C", Destination: "D", Miles: 15.0, Date: "2024-03-21"},
		{Origin: "E", Destination: "F", Miles: 5.0, Date: "2024-03-22"},
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
			trip:    Trip{Origin: "Home", Destination: "Work", Miles: 5.0, Date: "2024-03-20"},
			wantErr: false,
		},
		{
			name:    "empty origin",
			trip:    Trip{Origin: "", Destination: "Work", Miles: 5.0, Date: "2024-03-20"},
			wantErr: true,
		},
		{
			name:    "empty destination",
			trip:    Trip{Origin: "Home", Destination: "", Miles: 5.0, Date: "2024-03-20"},
			wantErr: true,
		},
		{
			name:    "negative miles",
			trip:    Trip{Origin: "Home", Destination: "Work", Miles: -5.0, Date: "2024-03-20"},
			wantErr: true,
		},
		{
			name:    "zero miles",
			trip:    Trip{Origin: "Home", Destination: "Work", Miles: 0.0, Date: "2024-03-20"},
			wantErr: true,
		},
		{
			name:    "empty date",
			trip:    Trip{Origin: "Home", Destination: "Work", Miles: 5.0, Date: ""},
			wantErr: true,
		},
		{
			name:    "invalid date format",
			trip:    Trip{Origin: "Home", Destination: "Work", Miles: 5.0, Date: "03-20-2024"},
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

func TestDateValidation(t *testing.T) {
	tests := []struct {
		name    string
		date    string
		wantErr bool
	}{
		{
			name:    "valid current date",
			date:    time.Now().Format("2006-01-02"),
			wantErr: false,
		},
		{
			name:    "valid future date",
			date:    time.Now().AddDate(0, 0, 1).Format("2006-01-02"),
			wantErr: false,
		},
		{
			name:    "valid past date",
			date:    time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
			wantErr: false,
		},
		{
			name:    "invalid month",
			date:    "2024-13-01",
			wantErr: true,
		},
		{
			name:    "invalid day",
			date:    "2024-02-30",
			wantErr: true,
		},
		{
			name:    "invalid year",
			date:    "0000-01-01",
			wantErr: true,
		},
		{
			name:    "invalid format",
			date:    "01/01/2024",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trip := Trip{
				Origin:      "Home",
				Destination: "Work",
				Miles:       5.0,
				Date:        tt.date,
			}
			err := trip.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Date validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTripDateOrdering(t *testing.T) {
	trips := []Trip{
		{Origin: "A", Destination: "B", Miles: 10.0, Date: "2024-03-22"},
		{Origin: "C", Destination: "D", Miles: 15.0, Date: "2024-03-20"},
		{Origin: "E", Destination: "F", Miles: 5.0, Date: "2024-03-21"},
	}

	// Sort trips by date
	sort.Slice(trips, func(i, j int) bool {
		dateI, _ := time.Parse("2006-01-02", trips[i].Date)
		dateJ, _ := time.Parse("2006-01-02", trips[j].Date)
		return dateI.Before(dateJ)
	})

	// Verify dates are in ascending order
	for i := 1; i < len(trips); i++ {
		dateI, _ := time.Parse("2006-01-02", trips[i].Date)
		datePrev, _ := time.Parse("2006-01-02", trips[i-1].Date)
		if dateI.Before(datePrev) {
			t.Errorf("Trips not properly sorted by date. Trip %d (%s) is before Trip %d (%s)",
				i, trips[i].Date, i-1, trips[i-1].Date)
		}
	}
}

func TestTripDateFiltering(t *testing.T) {
	trips := []Trip{
		{Origin: "A", Destination: "B", Miles: 10.0, Date: "2024-03-20"},
		{Origin: "C", Destination: "D", Miles: 15.0, Date: "2024-03-21"},
		{Origin: "E", Destination: "F", Miles: 5.0, Date: "2024-03-22"},
		{Origin: "G", Destination: "H", Miles: 8.0, Date: "2024-03-23"},
	}

	// Filter trips for a specific date
	targetDate := "2024-03-21"
	var filteredTrips []Trip
	for _, trip := range trips {
		if trip.Date == targetDate {
			filteredTrips = append(filteredTrips, trip)
		}
	}

	if len(filteredTrips) != 1 {
		t.Errorf("Expected 1 trip on %s, got %d", targetDate, len(filteredTrips))
	}

	if filteredTrips[0].Date != targetDate {
		t.Errorf("Expected trip date to be %s, got %s", targetDate, filteredTrips[0].Date)
	}
}

func TestCalculateWeeklySummaries(t *testing.T) {
	tests := []struct {
		name        string
		trips       []Trip
		ratePerMile float64
		want        []WeeklySummary
	}{
		{
			name:        "empty trips",
			trips:       []Trip{},
			ratePerMile: 0.655,
			want:        nil,
		},
		{
			name: "single week",
			trips: []Trip{
				{Date: "2024-03-17", Miles: 10.0}, // Sunday
				{Date: "2024-03-18", Miles: 15.0}, // Monday
				{Date: "2024-03-19", Miles: 20.0}, // Tuesday
			},
			ratePerMile: 0.655,
			want: []WeeklySummary{
				{
					WeekStart:   "2024-03-17",
					WeekEnd:     "2024-03-23",
					TotalMiles:  45.0,
					TotalAmount: 29.475,
				},
			},
		},
		{
			name: "multiple weeks",
			trips: []Trip{
				{Date: "2024-03-17", Miles: 10.0}, // Week 1
				{Date: "2024-03-18", Miles: 15.0}, // Week 1
				{Date: "2024-03-24", Miles: 20.0}, // Week 2
				{Date: "2024-03-25", Miles: 25.0}, // Week 2
			},
			ratePerMile: 0.655,
			want: []WeeklySummary{
				{
					WeekStart:   "2024-03-17",
					WeekEnd:     "2024-03-23",
					TotalMiles:  25.0,
					TotalAmount: 16.375,
				},
				{
					WeekStart:   "2024-03-24",
					WeekEnd:     "2024-03-30",
					TotalMiles:  45.0,
					TotalAmount: 29.475,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateWeeklySummaries(tt.trips, tt.ratePerMile)
			if len(got) != len(tt.want) {
				t.Errorf("CalculateWeeklySummaries() got %d summaries, want %d", len(got), len(tt.want))
				return
			}

			for i, summary := range got {
				if summary.WeekStart != tt.want[i].WeekStart {
					t.Errorf("WeekStart = %v, want %v", summary.WeekStart, tt.want[i].WeekStart)
				}
				if summary.WeekEnd != tt.want[i].WeekEnd {
					t.Errorf("WeekEnd = %v, want %v", summary.WeekEnd, tt.want[i].WeekEnd)
				}
				if summary.TotalMiles != tt.want[i].TotalMiles {
					t.Errorf("TotalMiles = %v, want %v", summary.TotalMiles, tt.want[i].TotalMiles)
				}
				if summary.TotalAmount != tt.want[i].TotalAmount {
					t.Errorf("TotalAmount = %v, want %v", summary.TotalAmount, tt.want[i].TotalAmount)
				}
			}
		})
	}
}
