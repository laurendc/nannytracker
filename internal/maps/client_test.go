package maps

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestNewClient(t *testing.T) {
	// Save original env var and restore after test
	originalKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	defer os.Setenv("GOOGLE_MAPS_API_KEY", originalKey)

	// Test case: No API key
	os.Unsetenv("GOOGLE_MAPS_API_KEY")
	_, err := NewClient()
	if err == nil {
		t.Error("Expected error when API key is not set")
	}

	// Test case: Valid API key
	os.Setenv("GOOGLE_MAPS_API_KEY", "test-key")
	client, err := NewClient()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if client == nil {
		t.Error("Expected client to be created")
	}
}

func TestCalculateDistance(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the request is valid
		if r.URL.Query().Get("origins") == "" || r.URL.Query().Get("destinations") == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Return a mock response (16093 meters is approximately 10 miles)
		response := `{
			"rows": [{
				"elements": [{
					"distance": {
						"value": 16093,
						"text": "10.0 mi"
					},
					"duration": {
						"value": 1200,
						"text": "20 mins"
					},
					"status": "OK"
				}]
			}],
			"status": "OK"
		}`
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer server.Close()

	// Create client with test server URL
	client := &Client{
		apiKey:     "test-key",
		httpClient: server.Client(),
		baseURL:    server.URL,
	}

	// Test successful case
	distance, err := client.CalculateDistance(context.Background(), "123 Main St", "456 Oak Ave")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expectedDistance := metersToMiles(16093)
	if distance != expectedDistance {
		t.Errorf("Expected distance %f, got %f", expectedDistance, distance)
	}

	// Test error cases
	_, err = client.CalculateDistance(context.Background(), "", "456 Oak Ave")
	if err == nil {
		t.Error("Expected error for empty origin")
	}

	_, err = client.CalculateDistance(context.Background(), "123 Main St", "")
	if err == nil {
		t.Error("Expected error for empty destination")
	}
}

func TestMetersToMiles(t *testing.T) {
	const conversionFactor = 0.000621371
	tests := []struct {
		meters int64
		want   float64
	}{
		{1609, float64(1609) * conversionFactor},
		{8046, float64(8046) * conversionFactor},
		{16093, float64(16093) * conversionFactor},
		{0, 0.0},
	}

	for _, tt := range tests {
		got := metersToMiles(tt.meters)
		if got != tt.want {
			t.Errorf("metersToMiles(%d) = %f, want %f", tt.meters, got, tt.want)
		}
	}
}
