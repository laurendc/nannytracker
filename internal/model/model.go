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
	Type        string  `json:"type"` // "single" or "round"
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
	if t.Type == "" {
		return errors.New("trip type cannot be empty")
	}
	if t.Type != "single" && t.Type != "round" {
		return errors.New("trip type must be either 'single' or 'round'")
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
		if t.Type == "round" {
			total += t.Miles * 2
		} else {
			total += t.Miles
		}
	}
	return total
}

// CalculateReimbursement calculates the total reimbursement amount
func CalculateReimbursement(trips []Trip, ratePerMile float64) float64 {
	return CalculateTotalMiles(trips) * ratePerMile
}

// Expense represents a reimbursable expense
type Expense struct {
	Date        string  `json:"date"`        // Format: YYYY-MM-DD
	Amount      float64 `json:"amount"`      // Amount in dollars
	Description string  `json:"description"` // Brief description of the expense
}

// Validate checks if an expense is valid
func (e Expense) Validate() error {
	if e.Date == "" {
		return errors.New("date cannot be empty")
	}
	if e.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	if e.Description == "" {
		return errors.New("description cannot be empty")
	}
	// Validate date format (YYYY-MM-DD)
	date, err := time.Parse("2006-01-02", e.Date)
	if err != nil {
		return errors.New("date must be in YYYY-MM-DD format")
	}
	// Check for invalid year (less than 1000)
	if date.Year() < 1000 {
		return errors.New("year must be at least 1000")
	}
	return nil
}

// CalculateTotalExpenses returns the sum of all expenses
func CalculateTotalExpenses(expenses []Expense) float64 {
	var total float64
	for _, e := range expenses {
		total += e.Amount
	}
	return total
}

// WeeklySummary represents the total miles and reimbursement for a week
type WeeklySummary struct {
	WeekStart     string // YYYY-MM-DD format
	WeekEnd       string // YYYY-MM-DD format
	TotalMiles    float64
	TotalAmount   float64
	TotalExpenses float64
}

// CalculateWeeklySummaries groups trips and expenses by week and calculates totals
func CalculateWeeklySummaries(trips []Trip, expenses []Expense, ratePerMile float64) []WeeklySummary {
	if len(trips) == 0 && len(expenses) == 0 {
		return nil
	}

	// Sort trips by date
	sort.Slice(trips, func(i, j int) bool {
		return trips[i].Date < trips[j].Date
	})

	// Sort expenses by date
	sort.Slice(expenses, func(i, j int) bool {
		return expenses[i].Date < expenses[j].Date
	})

	// Group trips and expenses by week
	weeklyTrips := make(map[string][]Trip)
	weeklyExpenses := make(map[string][]Expense)

	// Group trips
	for _, trip := range trips {
		t, err := time.Parse("2006-01-02", trip.Date)
		if err != nil {
			continue
		}
		weekStart := t.AddDate(0, 0, -int(t.Weekday()))
		weekKey := weekStart.Format("2006-01-02")
		weeklyTrips[weekKey] = append(weeklyTrips[weekKey], trip)
	}

	// Group expenses
	for _, expense := range expenses {
		t, err := time.Parse("2006-01-02", expense.Date)
		if err != nil {
			continue
		}
		weekStart := t.AddDate(0, 0, -int(t.Weekday()))
		weekKey := weekStart.Format("2006-01-02")
		weeklyExpenses[weekKey] = append(weeklyExpenses[weekKey], expense)
	}

	// Calculate summaries
	var summaries []WeeklySummary
	allWeeks := make(map[string]bool)

	// Collect all weeks
	for week := range weeklyTrips {
		allWeeks[week] = true
	}
	for week := range weeklyExpenses {
		allWeeks[week] = true
	}

	for weekStart := range allWeeks {
		weekTrips := weeklyTrips[weekStart]
		weekExpenses := weeklyExpenses[weekStart]

		totalMiles := CalculateTotalMiles(weekTrips)
		totalAmount := CalculateReimbursement(weekTrips, ratePerMile)
		totalExpenses := CalculateTotalExpenses(weekExpenses)

		// Calculate week end date
		start, _ := time.Parse("2006-01-02", weekStart)
		weekEnd := start.AddDate(0, 0, 6).Format("2006-01-02")

		summaries = append(summaries, WeeklySummary{
			WeekStart:     weekStart,
			WeekEnd:       weekEnd,
			TotalMiles:    totalMiles,
			TotalAmount:   totalAmount,
			TotalExpenses: totalExpenses,
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
	Expenses        []Expense       `json:"expenses"`
	WeeklySummaries []WeeklySummary `json:"weekly_summaries"`
}

// CalculateAndUpdateWeeklySummaries calculates weekly summaries and updates the storage data
func CalculateAndUpdateWeeklySummaries(data *StorageData, ratePerMile float64) {
	data.WeeklySummaries = CalculateWeeklySummaries(data.Trips, data.Expenses, ratePerMile)
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

// AddExpense adds a new expense to the storage data
func (d *StorageData) AddExpense(expense Expense) error {
	if err := expense.Validate(); err != nil {
		return err
	}
	d.Expenses = append(d.Expenses, expense)
	return nil
}

// EditExpense updates an expense at the specified index
func (d *StorageData) EditExpense(index int, newExpense Expense) error {
	if index < 0 || index >= len(d.Expenses) {
		return errors.New("invalid expense index")
	}
	if err := newExpense.Validate(); err != nil {
		return err
	}
	d.Expenses[index] = newExpense
	return nil
}

// DeleteExpense removes an expense at the specified index
func (d *StorageData) DeleteExpense(index int) error {
	if index < 0 || index >= len(d.Expenses) {
		return errors.New("invalid expense index")
	}
	d.Expenses = append(d.Expenses[:index], d.Expenses[index+1:]...)
	return nil
}
