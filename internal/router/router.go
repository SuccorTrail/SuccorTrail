package router

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"github.com/SuccorTrail/SuccorTrail/internal/handler"
	"github.com/SuccorTrail/SuccorTrail/internal/middleware"
	"github.com/SuccorTrail/SuccorTrail/internal/repository"
	"github.com/SuccorTrail/SuccorTrail/internal/util"
	"github.com/gorilla/mux"
)

func InitRouter() *mux.Router {
	r := mux.NewRouter()

	// Middleware
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.RecoveryMiddleware)

	// Get absolute path to the static directory
	cwd, _ := os.Getwd()
	staticPath := filepath.Join(cwd, "web/static")

	// Serve static files
	fs := http.FileServer(http.Dir(staticPath))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// Initialize repositories
	donationRepo := repository.NewDonationRepository()
	receiverRepo := repository.NewReceiverRepository()

	// Initialize handlers
	donationHandler := &handler.DonationHandler{Repo: donationRepo}
	receiverHandler := &handler.ReceiverHandler{Repo: receiverRepo}

	// API routes
	r.HandleFunc("/api/donations", donationHandler.CreateDonation).Methods("POST")
	r.HandleFunc("/api/receivers", receiverHandler.CreateReceiver).Methods("POST")
	
	// Add meals API route
	r.HandleFunc("/api/meals", func(w http.ResponseWriter, r *http.Request) {
		location := r.URL.Query().Get("location")
		if location == "" {
			http.Error(w, "Location is required", http.StatusBadRequest)
			return
		}

		// Simulate different meal availability scenarios
		var meals []map[string]interface{}
		
		// Uncomment the scenario you want to test
		
		// Scenario 1: No meals available
		meals = []map[string]interface{}{}
		
		// Scenario 2: Some meals available
		// meals = []map[string]interface{}{
		// 	{
		// 		"id": "meal1",
		// 		"type": "Vegetarian Pasta",
		// 		"quantity": 5,
		// 		"location": location,
		// 		"expiryDate": "2025-02-15T22:30:00Z",
		// 	},
		// }

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(meals)
	}).Methods("GET")

	// HTML routes
	r.HandleFunc("/", donationHandler.RenderDonationForm)
	r.HandleFunc("/donor", donationHandler.RenderDonationForm)
	r.HandleFunc("/receiver", receiverHandler.RenderReceiverForm)
	
	// Add meal-finder route
	r.HandleFunc("/meal-finder", func(w http.ResponseWriter, r *http.Request) {
		util.RenderTemplate(w, "meal-finder.html", nil)
	})

	return r
}
