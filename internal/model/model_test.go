package model

import (
	"encoding/json"
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
			name: "valid trip",
			trip: Trip{
				Origin:      "Home",
				Destination: "Work",
				Miles:       5.0,
				Date:        "2024-03-20",
				Type:        "single",
			},
			wantErr: false,
		},
		{
			name: "empty origin",
			trip: Trip{
				Origin:      "",
				Destination: "Work",
				Miles:       5.0,
				Date:        "2024-03-20",
				Type:        "single",
			},
			wantErr: true,
		},
		{
			name: "empty destination",
			trip: Trip{
				Origin:      "Home",
				Destination: "",
				Miles:       5.0,
				Date:        "2024-03-20",
				Type:        "single",
			},
			wantErr: true,
		},
		{
			name: "negative miles",
			trip: Trip{
				Origin:      "Home",
				Destination: "Work",
				Miles:       -5.0,
				Date:        "2024-03-20",
				Type:        "single",
			},
			wantErr: true,
		},
		{
			name: "zero miles",
			trip: Trip{
				Origin:      "Home",
				Destination: "Work",
				Miles:       0.0,
				Date:        "2024-03-20",
				Type:        "single",
			},
			wantErr: true,
		},
		{
			name: "empty date",
			trip: Trip{
				Origin:      "Home",
				Destination: "Work",
				Miles:       5.0,
				Date:        "",
				Type:        "single",
			},
			wantErr: true,
		},
		{
			name: "invalid date format",
			trip: Trip{
				Origin:      "Home",
				Destination: "Work",
				Miles:       5.0,
				Date:        "03-20-2024",
				Type:        "single",
			},
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
				Type:        "single",
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
		{Origin: "A", Destination: "B", Miles: 10.0, Date: "2024-03-22", Type: "single"},
		{Origin: "C", Destination: "D", Miles: 15.0, Date: "2024-03-20", Type: "single"},
		{Origin: "E", Destination: "F", Miles: 5.0, Date: "2024-03-21", Type: "single"},
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
		{Origin: "A", Destination: "B", Miles: 10.0, Date: "2024-03-20", Type: "single"},
		{Origin: "C", Destination: "D", Miles: 15.0, Date: "2024-03-21", Type: "single"},
		{Origin: "E", Destination: "F", Miles: 5.0, Date: "2024-03-22", Type: "single"},
		{Origin: "G", Destination: "H", Miles: 8.0, Date: "2024-03-23", Type: "single"},
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
		expenses    []Expense
		ratePerMile float64
		want        []WeeklySummary
	}{
		{
			name:        "no trips or expenses",
			trips:       []Trip{},
			expenses:    []Expense{},
			ratePerMile: 0.70,
			want:        nil,
		},
		{
			name: "single week with trips and expenses",
			trips: []Trip{
				{Date: "2024-03-20", Origin: "Home", Destination: "Work", Miles: 10, Type: "single"},
				{Date: "2024-03-21", Origin: "Work", Destination: "Home", Miles: 10, Type: "round"},
			},
			expenses: []Expense{
				{Date: "2024-03-20", Amount: 25.50, Description: "Lunch"},
				{Date: "2024-03-21", Amount: 15.75, Description: "Snacks"},
			},
			ratePerMile: 0.70,
			want: []WeeklySummary{
				{
					WeekStart:     "2024-03-17", // Sunday of the week
					WeekEnd:       "2024-03-23", // Saturday of the week
					TotalMiles:    30.0,         // 10 + (10 * 2)
					TotalAmount:   21.0,         // 30 miles * 0.70
					TotalExpenses: 41.25,        // 25.50 + 15.75
					Trips: []Trip{
						{Date: "2024-03-20", Origin: "Home", Destination: "Work", Miles: 10, Type: "single"},
						{Date: "2024-03-21", Origin: "Work", Destination: "Home", Miles: 10, Type: "round"},
					},
					Expenses: []Expense{
						{Date: "2024-03-20", Amount: 25.50, Description: "Lunch"},
						{Date: "2024-03-21", Amount: 15.75, Description: "Snacks"},
					},
				},
			},
		},
		{
			name: "multiple weeks",
			trips: []Trip{
				{Date: "2024-03-20", Origin: "Home", Destination: "Work", Miles: 10, Type: "single"},
				{Date: "2024-03-27", Origin: "Work", Destination: "Home", Miles: 10, Type: "round"},
			},
			expenses: []Expense{
				{Date: "2024-03-20", Amount: 25.50, Description: "Lunch"},
				{Date: "2024-03-27", Amount: 15.75, Description: "Snacks"},
			},
			ratePerMile: 0.70,
			want: []WeeklySummary{
				{
					WeekStart:     "2024-03-17", // First week
					WeekEnd:       "2024-03-23",
					TotalMiles:    10.0,
					TotalAmount:   7.0,   // 10 miles * 0.70
					TotalExpenses: 25.50, // First week expense
					Trips: []Trip{
						{Date: "2024-03-20", Origin: "Home", Destination: "Work", Miles: 10, Type: "single"},
					},
					Expenses: []Expense{
						{Date: "2024-03-20", Amount: 25.50, Description: "Lunch"},
					},
				},
				{
					WeekStart:     "2024-03-24", // Second week
					WeekEnd:       "2024-03-30",
					TotalMiles:    20.0,  // 10 * 2 for round trip
					TotalAmount:   14.0,  // 20 miles * 0.70
					TotalExpenses: 15.75, // Second week expense
					Trips: []Trip{
						{Date: "2024-03-27", Origin: "Work", Destination: "Home", Miles: 10, Type: "round"},
					},
					Expenses: []Expense{
						{Date: "2024-03-27", Amount: 15.75, Description: "Snacks"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateWeeklySummaries(tt.trips, tt.expenses, tt.ratePerMile)
			if len(got) != len(tt.want) {
				t.Errorf("CalculateWeeklySummaries() got %d summaries, want %d", len(got), len(tt.want))
				return
			}

			for i := range got {
				if got[i].WeekStart != tt.want[i].WeekStart {
					t.Errorf("Week %d: WeekStart got %v, want %v", i, got[i].WeekStart, tt.want[i].WeekStart)
				}
				if got[i].WeekEnd != tt.want[i].WeekEnd {
					t.Errorf("Week %d: WeekEnd got %v, want %v", i, got[i].WeekEnd, tt.want[i].WeekEnd)
				}
				if got[i].TotalMiles != tt.want[i].TotalMiles {
					t.Errorf("Week %d: TotalMiles got %v, want %v", i, got[i].TotalMiles, tt.want[i].TotalMiles)
				}
				if got[i].TotalAmount != tt.want[i].TotalAmount {
					t.Errorf("Week %d: TotalAmount got %v, want %v", i, got[i].TotalAmount, tt.want[i].TotalAmount)
				}
				if got[i].TotalExpenses != tt.want[i].TotalExpenses {
					t.Errorf("Week %d: TotalExpenses got %v, want %v", i, got[i].TotalExpenses, tt.want[i].TotalExpenses)
				}

				// Verify trips
				if len(got[i].Trips) != len(tt.want[i].Trips) {
					t.Errorf("Week %d: got %d trips, want %d", i, len(got[i].Trips), len(tt.want[i].Trips))
					continue
				}
				for j := range got[i].Trips {
					if got[i].Trips[j].Date != tt.want[i].Trips[j].Date {
						t.Errorf("Week %d, Trip %d: Date got %v, want %v", i, j, got[i].Trips[j].Date, tt.want[i].Trips[j].Date)
					}
					if got[i].Trips[j].Origin != tt.want[i].Trips[j].Origin {
						t.Errorf("Week %d, Trip %d: Origin got %v, want %v", i, j, got[i].Trips[j].Origin, tt.want[i].Trips[j].Origin)
					}
					if got[i].Trips[j].Destination != tt.want[i].Trips[j].Destination {
						t.Errorf("Week %d, Trip %d: Destination got %v, want %v", i, j, got[i].Trips[j].Destination, tt.want[i].Trips[j].Destination)
					}
					if got[i].Trips[j].Miles != tt.want[i].Trips[j].Miles {
						t.Errorf("Week %d, Trip %d: Miles got %v, want %v", i, j, got[i].Trips[j].Miles, tt.want[i].Trips[j].Miles)
					}
					if got[i].Trips[j].Type != tt.want[i].Trips[j].Type {
						t.Errorf("Week %d, Trip %d: Type got %v, want %v", i, j, got[i].Trips[j].Type, tt.want[i].Trips[j].Type)
					}
				}

				// Verify expenses
				if len(got[i].Expenses) != len(tt.want[i].Expenses) {
					t.Errorf("Week %d: got %d expenses, want %d", i, len(got[i].Expenses), len(tt.want[i].Expenses))
					continue
				}
				for j := range got[i].Expenses {
					if got[i].Expenses[j].Date != tt.want[i].Expenses[j].Date {
						t.Errorf("Week %d, Expense %d: Date got %v, want %v", i, j, got[i].Expenses[j].Date, tt.want[i].Expenses[j].Date)
					}
					if got[i].Expenses[j].Amount != tt.want[i].Expenses[j].Amount {
						t.Errorf("Week %d, Expense %d: Amount got %v, want %v", i, j, got[i].Expenses[j].Amount, tt.want[i].Expenses[j].Amount)
					}
					if got[i].Expenses[j].Description != tt.want[i].Expenses[j].Description {
						t.Errorf("Week %d, Expense %d: Description got %v, want %v", i, j, got[i].Expenses[j].Description, tt.want[i].Expenses[j].Description)
					}
				}
			}
		})
	}
}

func TestEditTrip(t *testing.T) {
	data := &StorageData{
		Trips: []Trip{
			{Date: "2024-03-20", Origin: "Home", Destination: "Work", Miles: 5.0, Type: "single"},
			{Date: "2024-03-21", Origin: "Work", Destination: "Store", Miles: 2.5, Type: "round"},
		},
	}

	// Test valid edit
	newTrip := Trip{Date: "2024-03-22", Origin: "Home", Destination: "Gym", Miles: 3.0, Type: "single"}
	if err := data.EditTrip(0, newTrip); err != nil {
		t.Errorf("EditTrip failed: %v", err)
	}
	if data.Trips[0] != newTrip {
		t.Errorf("Expected trip to be updated, got %+v", data.Trips[0])
	}

	// Test invalid index
	if err := data.EditTrip(2, newTrip); err == nil {
		t.Error("Expected error for invalid index")
	}

	// Test invalid trip
	invalidTrip := Trip{Date: "invalid", Origin: "Home", Destination: "Work", Miles: 5.0, Type: "single"}
	if err := data.EditTrip(0, invalidTrip); err == nil {
		t.Error("Expected error for invalid trip")
	}
}

func TestDeleteTrip(t *testing.T) {
	data := &StorageData{
		Trips: []Trip{
			{Date: "2024-03-20", Origin: "Home", Destination: "Work", Miles: 5.0},
			{Date: "2024-03-21", Origin: "Work", Destination: "Store", Miles: 2.5},
		},
	}

	// Test valid delete
	if err := data.DeleteTrip(0); err != nil {
		t.Errorf("DeleteTrip failed: %v", err)
	}
	if len(data.Trips) != 1 {
		t.Errorf("Expected 1 trip after deletion, got %d", len(data.Trips))
	}
	if data.Trips[0].Date != "2024-03-21" {
		t.Errorf("Expected remaining trip to be the second one, got %+v", data.Trips[0])
	}

	// Test invalid index
	if err := data.DeleteTrip(1); err == nil {
		t.Error("Expected error for invalid index")
	}
}

func TestTripTypeValidation(t *testing.T) {
	tests := []struct {
		name    string
		trip    Trip
		wantErr bool
	}{
		{
			name: "valid single trip",
			trip: Trip{
				Origin:      "Home",
				Destination: "Work",
				Miles:       5.0,
				Date:        "2024-03-20",
				Type:        "single",
			},
			wantErr: false,
		},
		{
			name: "valid round trip",
			trip: Trip{
				Origin:      "Home",
				Destination: "Work",
				Miles:       5.0,
				Date:        "2024-03-20",
				Type:        "round",
			},
			wantErr: false,
		},
		{
			name: "empty trip type",
			trip: Trip{
				Origin:      "Home",
				Destination: "Work",
				Miles:       5.0,
				Date:        "2024-03-20",
				Type:        "",
			},
			wantErr: true,
		},
		{
			name: "invalid trip type",
			trip: Trip{
				Origin:      "Home",
				Destination: "Work",
				Miles:       5.0,
				Date:        "2024-03-20",
				Type:        "invalid",
			},
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

func TestRoundTripMileageCalculation(t *testing.T) {
	trips := []Trip{
		{
			Origin:      "Home",
			Destination: "Work",
			Miles:       10.0,
			Date:        "2024-03-20",
			Type:        "single",
		},
		{
			Origin:      "Home",
			Destination: "Store",
			Miles:       5.0,
			Date:        "2024-03-20",
			Type:        "round",
		},
	}

	totalMiles := CalculateTotalMiles(trips)
	// First trip: 10 miles (single)
	// Second trip: 5 miles * 2 (round trip)
	expectedMiles := 20.0

	if totalMiles != expectedMiles {
		t.Errorf("CalculateTotalMiles() = %v, want %v", totalMiles, expectedMiles)
	}
}

func TestTripTypeSerialization(t *testing.T) {
	originalTrip := Trip{
		Origin:      "Home",
		Destination: "Work",
		Miles:       5.0,
		Date:        "2024-03-20",
		Type:        "round",
	}

	// Serialize to JSON
	data, err := json.Marshal(originalTrip)
	if err != nil {
		t.Fatalf("Failed to marshal trip: %v", err)
	}

	// Deserialize back
	var deserializedTrip Trip
	if err := json.Unmarshal(data, &deserializedTrip); err != nil {
		t.Fatalf("Failed to unmarshal trip: %v", err)
	}

	// Verify type was preserved
	if deserializedTrip.Type != originalTrip.Type {
		t.Errorf("Trip type not preserved during serialization. Got %s, want %s",
			deserializedTrip.Type, originalTrip.Type)
	}
}

func TestExpenseValidation(t *testing.T) {
	tests := []struct {
		name    string
		expense Expense
		wantErr bool
	}{
		{
			name: "valid expense",
			expense: Expense{
				Date:        "2024-03-20",
				Amount:      25.50,
				Description: "Lunch for kids",
			},
			wantErr: false,
		},
		{
			name: "empty date",
			expense: Expense{
				Date:        "",
				Amount:      25.50,
				Description: "Lunch for kids",
			},
			wantErr: true,
		},
		{
			name: "invalid date format",
			expense: Expense{
				Date:        "03/20/2024",
				Amount:      25.50,
				Description: "Lunch for kids",
			},
			wantErr: true,
		},
		{
			name: "zero amount",
			expense: Expense{
				Date:        "2024-03-20",
				Amount:      0,
				Description: "Lunch for kids",
			},
			wantErr: true,
		},
		{
			name: "negative amount",
			expense: Expense{
				Date:        "2024-03-20",
				Amount:      -25.50,
				Description: "Lunch for kids",
			},
			wantErr: true,
		},
		{
			name: "empty description",
			expense: Expense{
				Date:        "2024-03-20",
				Amount:      25.50,
				Description: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.expense.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Expense.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCalculateTotalExpenses(t *testing.T) {
	expenses := []Expense{
		{Date: "2024-03-20", Amount: 25.50, Description: "Lunch"},
		{Date: "2024-03-21", Amount: 15.75, Description: "Snacks"},
		{Date: "2024-03-22", Amount: 30.00, Description: "Activities"},
	}

	expected := 71.25
	got := CalculateTotalExpenses(expenses)
	if got != expected {
		t.Errorf("CalculateTotalExpenses() = %v, want %v", got, expected)
	}
}

func TestWeeklySummaries(t *testing.T) {
	trips := []Trip{
		{Date: "2024-03-20", Origin: "Home", Destination: "Work", Miles: 10, Type: "single"},
		{Date: "2024-03-21", Origin: "Work", Destination: "Home", Miles: 10, Type: "single"},
	}

	expenses := []Expense{} // Empty expenses list for this test
	ratePerMile := 0.70
	summaries := CalculateWeeklySummaries(trips, expenses, ratePerMile)

	if len(summaries) != 1 {
		t.Errorf("Expected 1 weekly summary, got %d", len(summaries))
	}

	summary := summaries[0]
	expectedMiles := 20.0
	expectedAmount := 14.0 // 20 miles * 0.70

	if summary.TotalMiles != expectedMiles {
		t.Errorf("Expected %v miles, got %v", expectedMiles, summary.TotalMiles)
	}
	if summary.TotalAmount != expectedAmount {
		t.Errorf("Expected $%v amount, got $%v", expectedAmount, summary.TotalAmount)
	}
}

func TestStorageDataExpenseOperations(t *testing.T) {
	data := &StorageData{
		Trips:    []Trip{},
		Expenses: []Expense{},
	}

	// Test AddExpense
	expense := Expense{
		Date:        "2024-03-20",
		Amount:      25.50,
		Description: "Lunch",
	}
	err := data.AddExpense(expense)
	if err != nil {
		t.Errorf("AddExpense() error = %v", err)
	}
	if len(data.Expenses) != 1 {
		t.Errorf("Expected 1 expense, got %d", len(data.Expenses))
	}

	// Test EditExpense
	editedExpense := Expense{
		Date:        "2024-03-20",
		Amount:      30.00,
		Description: "Lunch and snacks",
	}
	err = data.EditExpense(0, editedExpense)
	if err != nil {
		t.Errorf("EditExpense() error = %v", err)
	}
	if data.Expenses[0].Amount != 30.00 {
		t.Errorf("Expected amount $30.00, got $%v", data.Expenses[0].Amount)
	}

	// Test DeleteExpense
	err = data.DeleteExpense(0)
	if err != nil {
		t.Errorf("DeleteExpense() error = %v", err)
	}
	if len(data.Expenses) != 0 {
		t.Errorf("Expected 0 expenses, got %d", len(data.Expenses))
	}

	// Test invalid operations
	err = data.EditExpense(0, editedExpense)
	if err == nil {
		t.Error("Expected error for invalid index in EditExpense")
	}

	err = data.DeleteExpense(0)
	if err == nil {
		t.Error("Expected error for invalid index in DeleteExpense")
	}
}

func TestStorage(t *testing.T) {
	data := &StorageData{
		Trips:    []Trip{},
		Expenses: []Expense{},
	}

	// Add a trip
	trip := Trip{
		Date:        "2024-03-20",
		Origin:      "Home",
		Destination: "Work",
		Miles:       10,
		Type:        "single",
	}
	err := data.AddTrip(trip)
	if err != nil {
		t.Errorf("AddTrip() error = %v", err)
	}

	// Calculate weekly summaries
	ratePerMile := 0.70
	CalculateAndUpdateWeeklySummaries(data, ratePerMile)

	if len(data.WeeklySummaries) != 1 {
		t.Errorf("Expected 1 weekly summary, got %d", len(data.WeeklySummaries))
	}

	summary := data.WeeklySummaries[0]
	expectedMiles := 10.0
	expectedAmount := 7.0 // 10 miles * 0.70

	if summary.TotalMiles != expectedMiles {
		t.Errorf("Expected %v miles, got %v", expectedMiles, summary.TotalMiles)
	}
	if summary.TotalAmount != expectedAmount {
		t.Errorf("Expected $%v amount, got $%v", expectedAmount, summary.TotalAmount)
	}
}

func TestAddingTrip(t *testing.T) {
	data := &StorageData{
		Trips:    []Trip{},
		Expenses: []Expense{},
	}

	trip := Trip{
		Date:        "2024-03-20",
		Origin:      "Home",
		Destination: "Work",
		Miles:       10,
		Type:        "single",
	}

	err := data.AddTrip(trip)
	if err != nil {
		t.Errorf("AddTrip() error = %v", err)
	}

	if len(data.Trips) != 1 {
		t.Errorf("Expected 1 trip, got %d", len(data.Trips))
	}

	if data.Trips[0].Date != trip.Date {
		t.Errorf("Expected date %v, got %v", trip.Date, data.Trips[0].Date)
	}
	if data.Trips[0].Origin != trip.Origin {
		t.Errorf("Expected origin %v, got %v", trip.Origin, data.Trips[0].Origin)
	}
	if data.Trips[0].Destination != trip.Destination {
		t.Errorf("Expected destination %v, got %v", trip.Destination, data.Trips[0].Destination)
	}
	if data.Trips[0].Miles != trip.Miles {
		t.Errorf("Expected miles %v, got %v", trip.Miles, data.Trips[0].Miles)
	}
	if data.Trips[0].Type != trip.Type {
		t.Errorf("Expected type %v, got %v", trip.Type, data.Trips[0].Type)
	}
}

func TestStorageDataTripTemplateOperations(t *testing.T) {
	data := &StorageData{
		Trips:         []Trip{},
		Expenses:      []Expense{},
		TripTemplates: []TripTemplate{},
	}

	// Test adding a valid template
	template := TripTemplate{
		Name:        "Work Commute",
		Origin:      "123 Home St",
		Destination: "456 Work Ave",
		TripType:    "single",
		Notes:       "Regular work commute",
	}

	// Test adding template to storage
	if err := data.AddTripTemplate(template); err != nil {
		t.Errorf("AddTripTemplate() error = %v", err)
	}
	if len(data.TripTemplates) != 1 {
		t.Errorf("Expected 1 template, got %d", len(data.TripTemplates))
	}

	// Test editing template
	editedTemplate := TripTemplate{
		Name:        "Work Commute (Updated)",
		Origin:      "123 Home St",
		Destination: "456 Work Ave",
		TripType:    "round",
		Notes:       "Updated work commute",
	}
	if err := data.EditTripTemplate(0, editedTemplate); err != nil {
		t.Errorf("EditTripTemplate() error = %v", err)
	}
	if data.TripTemplates[0].Name != editedTemplate.Name {
		t.Errorf("Expected name %v, got %v", editedTemplate.Name, data.TripTemplates[0].Name)
	}
	if data.TripTemplates[0].TripType != editedTemplate.TripType {
		t.Errorf("Expected type %v, got %v", editedTemplate.TripType, data.TripTemplates[0].TripType)
	}

	// Test deleting template
	if err := data.DeleteTripTemplate(0); err != nil {
		t.Errorf("DeleteTripTemplate() error = %v", err)
	}
	if len(data.TripTemplates) != 0 {
		t.Errorf("Expected 0 templates after deletion, got %d", len(data.TripTemplates))
	}

	// Test invalid operations
	if err := data.EditTripTemplate(0, editedTemplate); err == nil {
		t.Error("Expected error for invalid index in EditTripTemplate")
	}

	if err := data.DeleteTripTemplate(0); err == nil {
		t.Error("Expected error for invalid index in DeleteTripTemplate")
	}
}

func TestTripTemplateSerialization(t *testing.T) {
	originalTemplate := TripTemplate{
		Name:        "Work Commute",
		Origin:      "123 Home St",
		Destination: "456 Work Ave",
		TripType:    "round",
		Notes:       "Regular work commute",
	}

	// Serialize to JSON
	data, err := json.Marshal(originalTemplate)
	if err != nil {
		t.Fatalf("Failed to marshal template: %v", err)
	}

	// Deserialize back
	var deserializedTemplate TripTemplate
	if err := json.Unmarshal(data, &deserializedTemplate); err != nil {
		t.Fatalf("Failed to unmarshal template: %v", err)
	}

	// Verify all fields were preserved
	if deserializedTemplate.Name != originalTemplate.Name {
		t.Errorf("Name not preserved during serialization. Got %s, want %s",
			deserializedTemplate.Name, originalTemplate.Name)
	}
	if deserializedTemplate.Origin != originalTemplate.Origin {
		t.Errorf("Origin not preserved during serialization. Got %s, want %s",
			deserializedTemplate.Origin, originalTemplate.Origin)
	}
	if deserializedTemplate.Destination != originalTemplate.Destination {
		t.Errorf("Destination not preserved during serialization. Got %s, want %s",
			deserializedTemplate.Destination, originalTemplate.Destination)
	}
	if deserializedTemplate.TripType != originalTemplate.TripType {
		t.Errorf("TripType not preserved during serialization. Got %s, want %s",
			deserializedTemplate.TripType, originalTemplate.TripType)
	}
	if deserializedTemplate.Notes != originalTemplate.Notes {
		t.Errorf("Notes not preserved during serialization. Got %s, want %s",
			deserializedTemplate.Notes, originalTemplate.Notes)
	}
}
