package maps

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

// DistanceCalculator is an interface for calculating distances between two points
type DistanceCalculator interface {
	CalculateDistance(ctx context.Context, origin, destination string) (float64, error)
}

// Client represents a Google Maps Distance Matrix API client
type Client struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

// DistanceMatrixResponse represents the response from the Distance Matrix API
type DistanceMatrixResponse struct {
	Rows []struct {
		Elements []struct {
			Distance struct {
				Value int64  `json:"value"` // Distance in meters
				Text  string `json:"text"`  // Human-readable distance
			} `json:"distance"`
			Duration struct {
				Value int64  `json:"value"` // Duration in seconds
				Text  string `json:"text"`  // Human-readable duration
			} `json:"duration"`
			Status string `json:"status"`
		} `json:"elements"`
	} `json:"rows"`
	Status string `json:"status"`
}

// NewClient creates a new Distance Matrix API client
func NewClient() (*Client, error) {
	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GOOGLE_MAPS_API_KEY environment variable not set. Please set it in your .env file")
	}

	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{},
		baseURL:    "https://maps.googleapis.com/maps/api/distancematrix/json",
	}, nil
}

// CalculateDistance calculates the distance between two addresses
func (c *Client) CalculateDistance(ctx context.Context, origin, destination string) (float64, error) {
	if origin == "" || destination == "" {
		return 0, fmt.Errorf("origin and destination addresses cannot be empty")
	}

	// Build the URL with query parameters
	params := url.Values{}
	params.Add("origins", origin)
	params.Add("destinations", destination)
	params.Add("key", c.apiKey)
	params.Add("units", "imperial") // Use miles instead of kilometers

	reqURL := fmt.Sprintf("%s?%s", c.baseURL, params.Encode())

	// Create the request
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	// Make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return 0, fmt.Errorf("invalid or unauthorized API key. Please check your GOOGLE_MAPS_API_KEY in .env file")
	}
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API request failed with status: %s", resp.Status)
	}

	// Parse the response
	var result DistanceMatrixResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	// Check API status
	if result.Status == "REQUEST_DENIED" {
		return 0, fmt.Errorf("API request denied. Please check your API key and billing status")
	}
	if result.Status != "OK" {
		return 0, fmt.Errorf("API returned error status: %s", result.Status)
	}

	// Check if we have valid results
	if len(result.Rows) == 0 || len(result.Rows[0].Elements) == 0 {
		return 0, fmt.Errorf("no distance information found")
	}

	element := result.Rows[0].Elements[0]
	if element.Status != "OK" {
		return 0, fmt.Errorf("distance calculation failed: %s", element.Status)
	}

	// Convert meters to miles
	return metersToMiles(element.Distance.Value), nil
}

// metersToMiles converts meters to miles
func metersToMiles(meters int64) float64 {
	return float64(meters) * 0.000621371
}
