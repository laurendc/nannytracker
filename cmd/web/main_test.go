package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/laurendc/nannytracker/pkg/config"
	core "github.com/laurendc/nannytracker/pkg/core"
)

func setupTestServer(t *testing.T) (*Server, string, func()) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "nannytracker-web-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create data directory
	dataDir := filepath.Join(tempDir, ".nannytracker")
	if err := os.MkdirAll(dataDir, 0750); err != nil {
		t.Fatalf("Failed to create data dir: %v", err)
	}

	// Create config with test data path
	cfg := &config.Config{
		DataDir:     dataDir,
		DataFile:    "trips.json",
		RatePerMile: 0.70,
	}

	// Create server
	server := NewServer(cfg)

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return server, tempDir, cleanup
}

func TestHealthEndpoint(t *testing.T) {
	server, _, cleanup := setupTestServer(t)
	defer cleanup()

	// Test GET /health
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	server.handleHealth(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]string
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", response["status"])
	}

	if response["service"] != "nannytracker-api" {
		t.Errorf("Expected service 'nannytracker-api', got '%s'", response["service"])
	}

	// Test wrong method
	req = httptest.NewRequest(http.MethodPost, "/health", nil)
	w = httptest.NewRecorder()
	server.handleHealth(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestTripsEndpoint(t *testing.T) {
	server, _, cleanup := setupTestServer(t)
	defer cleanup()

	// Test GET /api/trips (empty)
	req := httptest.NewRequest(http.MethodGet, "/api/trips", nil)
	w := httptest.NewRecorder()
	server.handleTrips(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["count"] != float64(0) {
		t.Errorf("Expected count 0, got %v", response["count"])
	}

	// Test POST /api/trips
	trip := core.Trip{
		Date:        "2024-12-18",
		Origin:      "Test Home",
		Destination: "Test Work",
		Miles:       5.0,
		Type:        "single",
	}

	tripJSON, _ := json.Marshal(trip)
	req = httptest.NewRequest(http.MethodPost, "/api/trips", bytes.NewBuffer(tripJSON))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	server.handleTrips(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var createdTrip core.Trip
	if err := json.NewDecoder(w.Body).Decode(&createdTrip); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if createdTrip.Date != trip.Date {
		t.Errorf("Expected date %s, got %s", trip.Date, createdTrip.Date)
	}

	// Test GET /api/trips (with data)
	req = httptest.NewRequest(http.MethodGet, "/api/trips", nil)
	w = httptest.NewRecorder()
	server.handleTrips(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["count"] != float64(1) {
		t.Errorf("Expected count 1, got %v", response["count"])
	}

	// Test CORS preflight
	req = httptest.NewRequest(http.MethodOptions, "/api/trips", nil)
	w = httptest.NewRecorder()
	server.handleTrips(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 for OPTIONS, got %d", w.Code)
	}

	// Test invalid method
	req = httptest.NewRequest(http.MethodPut, "/api/trips", nil)
	w = httptest.NewRecorder()
	server.handleTrips(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestTripsValidation(t *testing.T) {
	server, _, cleanup := setupTestServer(t)
	defer cleanup()

	// Test invalid trip data
	invalidTrip := map[string]interface{}{
		"date":        "invalid-date",
		"origin":      "",
		"destination": "Test Work",
		"miles":       -5.0,
		"type":        "invalid",
	}

	tripJSON, _ := json.Marshal(invalidTrip)
	req := httptest.NewRequest(http.MethodPost, "/api/trips", bytes.NewBuffer(tripJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.handleTrips(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid data, got %d", w.Code)
	}

	// Test invalid JSON
	req = httptest.NewRequest(http.MethodPost, "/api/trips", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	server.handleTrips(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid JSON, got %d", w.Code)
	}
}

func TestExpensesEndpoint(t *testing.T) {
	server, _, cleanup := setupTestServer(t)
	defer cleanup()

	// Test GET /api/expenses (empty)
	req := httptest.NewRequest(http.MethodGet, "/api/expenses", nil)
	w := httptest.NewRecorder()
	server.handleExpenses(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["count"] != float64(0) {
		t.Errorf("Expected count 0, got %v", response["count"])
	}

	// Test POST /api/expenses
	expense := core.Expense{
		Date:        "2024-12-18",
		Amount:      25.50,
		Description: "Test expense",
	}

	expenseJSON, _ := json.Marshal(expense)
	req = httptest.NewRequest(http.MethodPost, "/api/expenses", bytes.NewBuffer(expenseJSON))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	server.handleExpenses(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var createdExpense core.Expense
	if err := json.NewDecoder(w.Body).Decode(&createdExpense); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if createdExpense.Amount != expense.Amount {
		t.Errorf("Expected amount %.2f, got %.2f", expense.Amount, createdExpense.Amount)
	}

	// Test GET /api/expenses (with data)
	req = httptest.NewRequest(http.MethodGet, "/api/expenses", nil)
	w = httptest.NewRecorder()
	server.handleExpenses(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["count"] != float64(1) {
		t.Errorf("Expected count 1, got %v", response["count"])
	}
}

func TestExpensesValidation(t *testing.T) {
	server, _, cleanup := setupTestServer(t)
	defer cleanup()

	// Test invalid expense data
	invalidExpense := map[string]interface{}{
		"date":        "invalid-date",
		"amount":      -10.0,
		"description": "",
	}

	expenseJSON, _ := json.Marshal(invalidExpense)
	req := httptest.NewRequest(http.MethodPost, "/api/expenses", bytes.NewBuffer(expenseJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.handleExpenses(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid data, got %d", w.Code)
	}
}

func TestWeeklySummariesEndpoint(t *testing.T) {
	server, _, cleanup := setupTestServer(t)
	defer cleanup()

	// Test GET /api/summaries (empty)
	req := httptest.NewRequest(http.MethodGet, "/api/summaries", nil)
	w := httptest.NewRecorder()
	server.handleWeeklySummaries(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["count"] != float64(0) {
		t.Errorf("Expected count 0, got %v", response["count"])
	}

	// Test wrong method
	req = httptest.NewRequest(http.MethodPost, "/api/summaries", nil)
	w = httptest.NewRecorder()
	server.handleWeeklySummaries(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestWeeklySummariesWithData(t *testing.T) {
	server, _, cleanup := setupTestServer(t)
	defer cleanup()

	// Add some test data
	data, err := server.store.LoadData()
	if err != nil {
		t.Fatalf("Failed to load data: %v", err)
	}

	// Add a trip
	trip := core.Trip{
		Date:        "2024-12-18",
		Origin:      "Home",
		Destination: "Work",
		Miles:       10.0,
		Type:        "single",
	}
	data.Trips = append(data.Trips, trip)

	// Add an expense
	expense := core.Expense{
		Date:        "2024-12-18",
		Amount:      25.50,
		Description: "Lunch",
	}
	data.Expenses = append(data.Expenses, expense)

	if err := server.store.SaveData(data); err != nil {
		t.Fatalf("Failed to save data: %v", err)
	}

	// Test GET /api/summaries (with data)
	req := httptest.NewRequest(http.MethodGet, "/api/summaries", nil)
	w := httptest.NewRecorder()
	server.handleWeeklySummaries(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	summariesInterface, ok := response["summaries"]
	if !ok {
		t.Fatal("Expected 'summaries' key in response")
	}

	summaries, ok := summariesInterface.([]interface{})
	if !ok {
		t.Fatal("Expected summaries to be an array")
	}

	if len(summaries) == 0 {
		t.Error("Expected at least one weekly summary")
	}

	// Check that the summary contains the expected data
	summaryInterface, ok := summaries[0].(map[string]interface{})
	if !ok {
		t.Fatal("Expected summary to be a map")
	}

	if summaryInterface["TotalMiles"] == nil {
		t.Error("Expected TotalMiles in summary")
	}
}

func TestCORSHeaders(t *testing.T) {
	server, _, cleanup := setupTestServer(t)
	defer cleanup()

	// Test CORS headers on trips endpoint
	req := httptest.NewRequest(http.MethodGet, "/api/trips", nil)
	w := httptest.NewRecorder()
	server.handleTrips(w, req)

	headers := w.Header()
	if headers.Get("Access-Control-Allow-Origin") != "*" {
		t.Error("Expected CORS Access-Control-Allow-Origin header")
	}

	if headers.Get("Access-Control-Allow-Methods") == "" {
		t.Error("Expected CORS Access-Control-Allow-Methods header")
	}

	if headers.Get("Access-Control-Allow-Headers") == "" {
		t.Error("Expected CORS Access-Control-Allow-Headers header")
	}
}

func TestServerCreation(t *testing.T) {
	// Test server creation with valid config
	cfg := &config.Config{
		DataDir:     "/tmp/test",
		DataFile:    "test.json",
		RatePerMile: 0.70,
	}

	server := NewServer(cfg)
	if server == nil {
		t.Fatal("Expected server to be created")
	}

	if server.store == nil {
		t.Fatal("Expected storage to be initialized")
	}

	if server.cfg != cfg {
		t.Error("Expected config to be set")
	}
}

func setupBenchmarkServer(b *testing.B) (*Server, func()) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "nannytracker-web-bench")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create data directory
	dataDir := filepath.Join(tempDir, ".nannytracker")
	if err := os.MkdirAll(dataDir, 0750); err != nil {
		b.Fatalf("Failed to create data dir: %v", err)
	}

	// Create config with test data path
	cfg := &config.Config{
		DataDir:     dataDir,
		DataFile:    "trips.json",
		RatePerMile: 0.70,
	}

	// Create server
	server := NewServer(cfg)

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return server, cleanup
}

func BenchmarkTripsEndpoint(b *testing.B) {
	server, cleanup := setupBenchmarkServer(b)
	defer cleanup()

	// Add some test data
	data, err := server.store.LoadData()
	if err != nil {
		b.Fatalf("Failed to load data: %v", err)
	}

	for i := 0; i < 100; i++ {
		trip := core.Trip{
			Date:        fmt.Sprintf("2024-12-%02d", i%30+1),
			Origin:      fmt.Sprintf("Home %d", i),
			Destination: fmt.Sprintf("Work %d", i),
			Miles:       float64(i) + 1.0,
			Type:        "single",
		}
		data.Trips = append(data.Trips, trip)
	}

	if err := server.store.SaveData(data); err != nil {
		b.Fatalf("Failed to save data: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/api/trips", nil)
		w := httptest.NewRecorder()
		server.handleTrips(w, req)
	}
}

func TestHTTPServerConfig(t *testing.T) {
	port := "12345"
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      nil,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if srv.ReadTimeout != 10*time.Second {
		t.Errorf("Expected ReadTimeout to be 10s, got %v", srv.ReadTimeout)
	}
	if srv.WriteTimeout != 10*time.Second {
		t.Errorf("Expected WriteTimeout to be 10s, got %v", srv.WriteTimeout)
	}
	if srv.IdleTimeout != 60*time.Second {
		t.Errorf("Expected IdleTimeout to be 60s, got %v", srv.IdleTimeout)
	}
}
