package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/laurendc/nannytracker/pkg/config"
	model "github.com/laurendc/nannytracker/pkg/core"
	"github.com/laurendc/nannytracker/pkg/core/maps"
	"github.com/laurendc/nannytracker/pkg/core/storage"
	"github.com/laurendc/nannytracker/pkg/version"
)

type Server struct {
	store      *storage.FileStorage
	cfg        *config.Config
	mapsClient maps.DistanceCalculator
}

func NewServer(cfg *config.Config) (*Server, error) {
	store := storage.New(cfg.DataPath())

	// Initialize Google Maps client
	var mapsClient maps.DistanceCalculator
	realClient, err := maps.NewClient()
	if err != nil {
		// Fall back to mock client if Google Maps API is not available
		log.Printf("Google Maps API not available, using mock client: %v", err)
		mapsClient = maps.NewMockClient()
	} else {
		mapsClient = realClient
	}

	return &Server{
		store:      store,
		cfg:        cfg,
		mapsClient: mapsClient,
	}, nil
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"service": "nannytracker-api",
	}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleVersion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(version.Get()); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleTrips(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle CORS preflight
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.getTrips(w, r)
	case http.MethodPost:
		s.createTrip(w, r)
	case http.MethodPut:
		s.updateTrip(w, r)
	case http.MethodDelete:
		s.deleteTrip(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) getTrips(w http.ResponseWriter, r *http.Request) {
	data, err := s.store.LoadData()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load data: %v", err), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"trips": data.Trips,
		"count": len(data.Trips),
	}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (s *Server) createTrip(w http.ResponseWriter, r *http.Request) {
	// Create a struct for the incoming trip data without miles
	var tripData struct {
		Date        string `json:"date"`
		Origin      string `json:"origin"`
		Destination string `json:"destination"`
		Type        string `json:"type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&tripData); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Basic validation of required fields
	if tripData.Date == "" {
		http.Error(w, "Date is required", http.StatusBadRequest)
		return
	}
	if tripData.Origin == "" {
		http.Error(w, "Origin is required", http.StatusBadRequest)
		return
	}
	if tripData.Destination == "" {
		http.Error(w, "Destination is required", http.StatusBadRequest)
		return
	}
	if tripData.Type == "" {
		http.Error(w, "Type is required", http.StatusBadRequest)
		return
	}
	if tripData.Type != "single" && tripData.Type != "round" {
		http.Error(w, "Type must be 'single' or 'round'", http.StatusBadRequest)
		return
	}

	// Calculate miles using Google Maps API
	distance, err := s.mapsClient.CalculateDistance(context.Background(), tripData.Origin, tripData.Destination)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to calculate distance: %v", err), http.StatusInternalServerError)
		return
	}

	// Create the complete trip with calculated miles
	trip := model.Trip{
		Date:        tripData.Date,
		Origin:      tripData.Origin,
		Destination: tripData.Destination,
		Type:        tripData.Type,
		Miles:       distance,
	}

	// Validate the complete trip
	if err := trip.Validate(); err != nil {
		http.Error(w, fmt.Sprintf("Invalid trip data: %v", err), http.StatusBadRequest)
		return
	}

	// Load existing data
	data, err := s.store.LoadData()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load data: %v", err), http.StatusInternalServerError)
		return
	}

	// Add the new trip
	data.Trips = append(data.Trips, trip)

	// Save the updated data
	if err := s.store.SaveData(data); err != nil {
		http.Error(w, fmt.Sprintf("Failed to save data: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(trip); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (s *Server) updateTrip(w http.ResponseWriter, r *http.Request) {
	// Extract index from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/trips/")
	if path == "" {
		http.Error(w, "Trip index is required", http.StatusBadRequest)
		return
	}

	index, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid trip index", http.StatusBadRequest)
		return
	}

	var trip model.Trip
	if err := json.NewDecoder(r.Body).Decode(&trip); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Validate the trip
	if err := trip.Validate(); err != nil {
		http.Error(w, fmt.Sprintf("Invalid trip data: %v", err), http.StatusBadRequest)
		return
	}

	// Load existing data
	data, err := s.store.LoadData()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load data: %v", err), http.StatusInternalServerError)
		return
	}

	// Update the trip
	if err := data.EditTrip(index, trip); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update trip: %v", err), http.StatusBadRequest)
		return
	}

	// Save the updated data
	if err := s.store.SaveData(data); err != nil {
		http.Error(w, fmt.Sprintf("Failed to save data: %v", err), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(trip); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (s *Server) deleteTrip(w http.ResponseWriter, r *http.Request) {
	// Extract index from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/trips/")
	if path == "" {
		http.Error(w, "Trip index is required", http.StatusBadRequest)
		return
	}

	index, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid trip index", http.StatusBadRequest)
		return
	}

	// Load existing data
	data, err := s.store.LoadData()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load data: %v", err), http.StatusInternalServerError)
		return
	}

	// Delete the trip
	if err := data.DeleteTrip(index); err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete trip: %v", err), http.StatusBadRequest)
		return
	}

	// Save the updated data
	if err := s.store.SaveData(data); err != nil {
		http.Error(w, fmt.Sprintf("Failed to save data: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleExpenses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle CORS preflight
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	switch r.Method {
	case http.MethodGet:
		s.getExpenses(w, r)
	case http.MethodPost:
		s.createExpense(w, r)
	case http.MethodPut:
		s.updateExpense(w, r)
	case http.MethodDelete:
		s.deleteExpense(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) getExpenses(w http.ResponseWriter, r *http.Request) {
	data, err := s.store.LoadData()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load data: %v", err), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"expenses": data.Expenses,
		"count":    len(data.Expenses),
	}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (s *Server) createExpense(w http.ResponseWriter, r *http.Request) {
	var expense model.Expense
	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Validate the expense
	if err := expense.Validate(); err != nil {
		http.Error(w, fmt.Sprintf("Invalid expense data: %v", err), http.StatusBadRequest)
		return
	}

	// Load existing data
	data, err := s.store.LoadData()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load data: %v", err), http.StatusInternalServerError)
		return
	}

	// Add the new expense
	data.Expenses = append(data.Expenses, expense)

	// Save the updated data
	if err := s.store.SaveData(data); err != nil {
		http.Error(w, fmt.Sprintf("Failed to save data: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(expense); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (s *Server) updateExpense(w http.ResponseWriter, r *http.Request) {
	// Extract index from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/expenses/")
	if path == "" {
		http.Error(w, "Expense index is required", http.StatusBadRequest)
		return
	}

	index, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid expense index", http.StatusBadRequest)
		return
	}

	var expense model.Expense
	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Validate the expense
	if err := expense.Validate(); err != nil {
		http.Error(w, fmt.Sprintf("Invalid expense data: %v", err), http.StatusBadRequest)
		return
	}

	// Load existing data
	data, err := s.store.LoadData()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load data: %v", err), http.StatusInternalServerError)
		return
	}

	// Update the expense
	if err := data.EditExpense(index, expense); err != nil {
		http.Error(w, fmt.Sprintf("Failed to update expense: %v", err), http.StatusBadRequest)
		return
	}

	// Save the updated data
	if err := s.store.SaveData(data); err != nil {
		http.Error(w, fmt.Sprintf("Failed to save data: %v", err), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(expense); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (s *Server) deleteExpense(w http.ResponseWriter, r *http.Request) {
	// Extract index from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/expenses/")
	if path == "" {
		http.Error(w, "Expense index is required", http.StatusBadRequest)
		return
	}

	index, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid expense index", http.StatusBadRequest)
		return
	}

	// Load existing data
	data, err := s.store.LoadData()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load data: %v", err), http.StatusInternalServerError)
		return
	}

	// Delete the expense
	if err := data.DeleteExpense(index); err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete expense: %v", err), http.StatusBadRequest)
		return
	}

	// Save the updated data
	if err := s.store.SaveData(data); err != nil {
		http.Error(w, fmt.Sprintf("Failed to save data: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleWeeklySummaries(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	data, err := s.store.LoadData()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to load data: %v", err), http.StatusInternalServerError)
		return
	}

	// Calculate weekly summaries
	summaries := model.CalculateWeeklySummaries(data.Trips, data.Expenses, s.cfg.RatePerMile)

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"summaries": summaries,
		"count":     len(summaries),
	}); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func main() {
	// Parse command line flags
	var showVersion bool
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.BoolVar(&showVersion, "v", false, "Show version information")
	flag.Parse()

	// Show version if requested
	if showVersion {
		fmt.Println(version.FullString())
		os.Exit(0)
	}

	// Load .env file from project root
	config.LoadEnv()

	// Load configuration
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create server
	server, err := NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Set up routes
	http.HandleFunc("/health", server.handleHealth)
	http.HandleFunc("/version", server.handleVersion)
	http.HandleFunc("/api/trips", server.handleTrips)
	http.HandleFunc("/api/trips/", server.handleTrips) // Handle /api/trips/{index}
	http.HandleFunc("/api/expenses", server.handleExpenses)
	http.HandleFunc("/api/expenses/", server.handleExpenses) // Handle /api/expenses/{index}
	http.HandleFunc("/api/summaries", server.handleWeeklySummaries)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting NannyTracker API server on port %s", port)
	log.Printf("API endpoints:")
	log.Printf("  GET  /health")
	log.Printf("  GET  /version")
	log.Printf("  GET  /api/trips")
	log.Printf("  POST /api/trips")
	log.Printf("  PUT  /api/trips/{index}")
	log.Printf("  DELETE /api/trips/{index}")
	log.Printf("  GET  /api/expenses")
	log.Printf("  POST /api/expenses")
	log.Printf("  PUT  /api/expenses/{index}")
	log.Printf("  DELETE /api/expenses/{index}")
	log.Printf("  GET  /api/summaries")

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      nil,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
