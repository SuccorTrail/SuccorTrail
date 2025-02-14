package router

import (
	"net/http"

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

	// Serve static files
	fs := http.FileServer(http.Dir("../../web/static"))
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
