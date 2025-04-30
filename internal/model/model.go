package model

import (
	"errors"
	"sort"
	"time"
)

// Trip represents a single trip with origin, destination, and mileage
type Trip struct {
	Origin      string  `json:"origin"`
	Destination string  `json:"destination"`
	Miles       float64 `json:"miles"`
	Date        string  `json:"date"` // Format: YYYY-MM-DD
}

// Validate checks if a trip is valid
func (t Trip) Validate() error {
	if t.Origin == "" {
		return errors.New("origin cannot be empty")
	}
	if t.Destination == "" {
		return errors.New("destination cannot be empty")
	}
	if t.Miles <= 0 {
		return errors.New("miles must be greater than 0")
	}
	if t.Date == "" {
		return errors.New("date cannot be empty")
	}
	// Validate date format (YYYY-MM-DD)
	date, err := time.Parse("2006-01-02", t.Date)
	if err != nil {
		return errors.New("date must be in YYYY-MM-DD format")
	}
	// Check for invalid year (less than 1000)
	if date.Year() < 1000 {
		return errors.New("year must be at least 1000")
	}
	return nil
}

// CalculateTotalMiles returns the sum of miles for all trips
func CalculateTotalMiles(trips []Trip) float64 {
	var total float64
	for _, t := range trips {
		total += t.Miles
	}
	return total
}

// CalculateReimbursement calculates the total reimbursement amount
func CalculateReimbursement(trips []Trip, ratePerMile float64) float64 {
	return CalculateTotalMiles(trips) * ratePerMile
}

// WeeklySummary represents the total miles and reimbursement for a week
type WeeklySummary struct {
	WeekStart   string // YYYY-MM-DD format
	WeekEnd     string // YYYY-MM-DD format
	TotalMiles  float64
	TotalAmount float64
}

// CalculateWeeklySummaries groups trips by week and calculates totals
func CalculateWeeklySummaries(trips []Trip, ratePerMile float64) []WeeklySummary {
	if len(trips) == 0 {
		return nil
	}

	// Sort trips by date
	sort.Slice(trips, func(i, j int) bool {
		return trips[i].Date < trips[j].Date
	})

	// Group trips by week
	weeklyTrips := make(map[string][]Trip)
	for _, trip := range trips {
		// Parse the date
		t, err := time.Parse("2006-01-02", trip.Date)
		if err != nil {
			continue
		}

		// Get the start of the week (Sunday)
		weekStart := t.AddDate(0, 0, -int(t.Weekday()))
		weekKey := weekStart.Format("2006-01-02")
		weeklyTrips[weekKey] = append(weeklyTrips[weekKey], trip)
	}

	// Calculate summaries
	var summaries []WeeklySummary
	for weekStart, weekTrips := range weeklyTrips {
		totalMiles := CalculateTotalMiles(weekTrips)
		totalAmount := CalculateReimbursement(weekTrips, ratePerMile)

		// Calculate week end date
		start, _ := time.Parse("2006-01-02", weekStart)
		weekEnd := start.AddDate(0, 0, 6).Format("2006-01-02")

		summaries = append(summaries, WeeklySummary{
			WeekStart:   weekStart,
			WeekEnd:     weekEnd,
			TotalMiles:  totalMiles,
			TotalAmount: totalAmount,
		})
	}

	// Sort summaries by week start date
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].WeekStart < summaries[j].WeekStart
	})

	return summaries
}

// StorageData represents the complete data structure stored in the JSON file
type StorageData struct {
	Trips           []Trip          `json:"trips"`
	WeeklySummaries []WeeklySummary `json:"weekly_summaries"`
}

// CalculateAndUpdateWeeklySummaries calculates weekly summaries and updates the storage data
func CalculateAndUpdateWeeklySummaries(data *StorageData, ratePerMile float64) {
	data.WeeklySummaries = CalculateWeeklySummaries(data.Trips, ratePerMile)
}

// EditTrip updates a trip at the specified index
func (d *StorageData) EditTrip(index int, newTrip Trip) error {
	if index < 0 || index >= len(d.Trips) {
		return errors.New("invalid trip index")
	}
	if err := newTrip.Validate(); err != nil {
		return err
	}
	d.Trips[index] = newTrip
	return nil
}

// DeleteTrip removes a trip at the specified index
func (d *StorageData) DeleteTrip(index int) error {
	if index < 0 || index >= len(d.Trips) {
		return errors.New("invalid trip index")
	}
	d.Trips = append(d.Trips[:index], d.Trips[index+1:]...)
	return nil
}
