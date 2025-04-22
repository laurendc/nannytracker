package model

import "errors"

// Trip represents a single trip with origin, destination, and mileage
type Trip struct {
	Origin      string  `json:"origin"`
	Destination string  `json:"destination"`
	Miles       float64 `json:"miles"`
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
