package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/laurendc/nannytracker/pkg/config"
	core "github.com/laurendc/nannytracker/pkg/core"
	"github.com/laurendc/nannytracker/pkg/core/storage"
	"github.com/laurendc/nannytracker/pkg/version"
)

type Server struct {
	store *storage.FileStorage
	cfg   *config.Config
}

func NewServer(cfg *config.Config) *Server {
	store := storage.New(cfg.DataPath())
	return &Server{
		store: store,
		cfg:   cfg,
	}
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
	var trip core.Trip
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
	var expense core.Expense
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
	summaries := core.CalculateWeeklySummaries(data.Trips, data.Expenses, s.cfg.RatePerMile)

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
	server := NewServer(cfg)

	// Set up routes
	http.HandleFunc("/health", server.handleHealth)
	http.HandleFunc("/version", server.handleVersion)
	http.HandleFunc("/api/trips", server.handleTrips)
	http.HandleFunc("/api/expenses", server.handleExpenses)
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
	log.Printf("  GET  /api/expenses")
	log.Printf("  POST /api/expenses")
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
