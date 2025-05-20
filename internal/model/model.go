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

// RecurringTrip represents a trip that occurs weekly
type RecurringTrip struct {
	Origin      string  `json:"origin"`
	Destination string  `json:"destination"`
	Miles       float64 `json:"miles"`
	StartDate   string  `json:"start_date"` // Format: YYYY-MM-DD
	EndDate     string  `json:"end_date"`   // Format: YYYY-MM-DD, optional
	Type        string  `json:"type"`       // "single" or "round"
	Weekday     int     `json:"weekday"`    // 0-6, where 0 is Sunday
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

// Validate checks if a recurring trip is valid
func (rt RecurringTrip) Validate() error {
	if rt.Origin == "" {
		return errors.New("origin cannot be empty")
	}
	if rt.Destination == "" {
		return errors.New("destination cannot be empty")
	}
	if rt.Miles <= 0 {
		return errors.New("miles must be greater than 0")
	}
	if rt.StartDate == "" {
		return errors.New("start date cannot be empty")
	}
	if rt.Type == "" {
		return errors.New("trip type cannot be empty")
	}
	if rt.Type != "single" && rt.Type != "round" {
		return errors.New("trip type must be either 'single' or 'round'")
	}
	if rt.Weekday < 0 || rt.Weekday > 6 {
		return errors.New("weekday must be between 0 (Sunday) and 6 (Saturday)")
	}

	// Validate start date format
	startDate, err := time.Parse("2006-01-02", rt.StartDate)
	if err != nil {
		return errors.New("start date must be in YYYY-MM-DD format")
	}
	if startDate.Year() < 1000 {
		return errors.New("start year must be at least 1000")
	}

	// Validate end date if provided
	if rt.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", rt.EndDate)
		if err != nil {
			return errors.New("end date must be in YYYY-MM-DD format")
		}
		if endDate.Before(startDate) {
			return errors.New("end date must be after start date")
		}
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
	Trips         []Trip    // Itemized list of trips for this week
	Expenses      []Expense // Itemized list of expenses for this week
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
			Trips:         weekTrips,
			Expenses:      weekExpenses,
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
	RecurringTrips  []RecurringTrip `json:"recurring_trips"`
	Expenses        []Expense       `json:"expenses"`
	WeeklySummaries []WeeklySummary `json:"weekly_summaries"`
	ReferenceDate   string          `json:"reference_date,omitempty"` // For testing purposes
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

// GenerateTrips generates individual trips from a recurring trip for a given date range
func (rt RecurringTrip) GenerateTrips(startDate, endDate time.Time) []Trip {
	var trips []Trip
	current := startDate

	// If the start date is not the target weekday, find the next occurrence
	if current.Weekday() != time.Weekday(rt.Weekday) {
		daysUntilNext := (rt.Weekday - int(current.Weekday()) + 7) % 7
		current = current.AddDate(0, 0, daysUntilNext)
	}

	// Generate trips for each occurrence until end date
	for !current.After(endDate) {
		trip := Trip{
			Origin:      rt.Origin,
			Destination: rt.Destination,
			Miles:       rt.Miles,
			Date:        current.Format("2006-01-02"),
			Type:        rt.Type,
		}
		trips = append(trips, trip)
		current = current.AddDate(0, 0, 7) // Add one week
	}

	return trips
}

// GenerateTripsFromRecurring generates individual trips from all recurring trips
func (d *StorageData) GenerateTripsFromRecurring() error {
	// Get current date and end of current month
	var now time.Time
	var err error
	if d.ReferenceDate != "" {
		now, err = time.Parse("2006-01-02", d.ReferenceDate)
		if err != nil {
			return err
		}
	} else {
		now = time.Now()
	}
	endOfMonth := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, now.Location())

	// Create a map to track existing trip dates
	existingDates := make(map[string]bool)
	for _, trip := range d.Trips {
		existingDates[trip.Date] = true
	}

	// Generate trips for each recurring trip
	for _, rt := range d.RecurringTrips {
		startDate, err := time.Parse("2006-01-02", rt.StartDate)
		if err != nil {
			return err
		}

		// Use end date from recurring trip if provided, otherwise use end of month
		endDate := endOfMonth
		if rt.EndDate != "" {
			parsedEndDate, err := time.Parse("2006-01-02", rt.EndDate)
			if err != nil {
				return err
			}
			if parsedEndDate.Before(endDate) {
				endDate = parsedEndDate
			}
		}

		// Generate trips and add them to the storage
		trips := rt.GenerateTrips(startDate, endDate)
		for _, trip := range trips {
			// Only add the trip if it doesn't already exist for that date
			if !existingDates[trip.Date] {
				if err := d.AddTrip(trip); err != nil {
					return err
				}
				existingDates[trip.Date] = true
			}
		}
	}

	return nil
}

// AddRecurringTrip adds a new recurring trip to the storage data
func (d *StorageData) AddRecurringTrip(trip RecurringTrip) error {
	if err := trip.Validate(); err != nil {
		return err
	}
	d.RecurringTrips = append(d.RecurringTrips, trip)
	return nil
}

// EditRecurringTrip updates a recurring trip at the specified index
func (d *StorageData) EditRecurringTrip(index int, newTrip RecurringTrip) error {
	if index < 0 || index >= len(d.RecurringTrips) {
		return errors.New("invalid recurring trip index")
	}
	if err := newTrip.Validate(); err != nil {
		return err
	}
	d.RecurringTrips[index] = newTrip
	return nil
}

// DeleteRecurringTrip removes a recurring trip at the specified index
func (d *StorageData) DeleteRecurringTrip(index int) error {
	if index < 0 || index >= len(d.RecurringTrips) {
		return errors.New("invalid recurring trip index")
	}
	d.RecurringTrips = append(d.RecurringTrips[:index], d.RecurringTrips[index+1:]...)
	return nil
}

// AddTrip adds a new trip to the storage data
func (d *StorageData) AddTrip(trip Trip) error {
	if err := trip.Validate(); err != nil {
		return err
	}
	d.Trips = append(d.Trips, trip)
	return nil
}
