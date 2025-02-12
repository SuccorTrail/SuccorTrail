package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type Donation struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Phone         string    `json:"phone"`
	MealType      string    `json:"mealType"`
	Quantity      int       `json:"quantity"`
	ExpiryDate    time.Time `json:"expiryDate"`
	Location      string    `json:"location"`
	Notes         string    `json:"notes"`
	CreatedAt     time.Time `json:"createdAt"`
}

type Receiver struct {
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	Phone              string   `json:"phone"`
	Location           string   `json:"location"`
	FamilySize         int      `json:"familySize"`
	DietaryRestrictions []string `json:"dietaryRestrictions"`
	CreatedAt          time.Time `json:"createdAt"`
}

type Feedback struct {
	ID        string    `json:"id"`
	DonationID string    `json:"donationId"`
	Quality   int       `json:"quality"`
	Comments  string    `json:"comments"`
	CreatedAt time.Time `json:"createdAt"`
}

var db *sql.DB

func main() {
	// Initialize database
	var err error
	db, err = initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize router
	r := mux.NewRouter()

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// API routes
	r.HandleFunc("/api/donations", createDonation).Methods("POST")
	r.HandleFunc("/api/receivers", createReceiver).Methods("POST")
	r.HandleFunc("/api/meals", getAvailableMeals).Methods("GET")
	r.HandleFunc("/api/verify-meal", verifyMeal).Methods("POST")
	r.HandleFunc("/api/feedback", submitFeedback).Methods("POST")

	// HTML routes
	r.HandleFunc("/", serveTemplate("index.html"))
	r.HandleFunc("/donor", serveTemplate("donor.html"))
	r.HandleFunc("/receiver", serveTemplate("receiver.html"))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "succortrail.db")
	if err != nil {
		return nil, err
	}

	// Create tables
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS donations (
		id TEXT PRIMARY KEY,
		name TEXT,
		email TEXT,
		phone TEXT,
		meal_type TEXT,
		quantity INTEGER,
		expiry_date DATETIME,
		location TEXT,
		notes TEXT,
		created_at DATETIME
	);
	CREATE TABLE IF NOT EXISTS receivers (
		id TEXT PRIMARY KEY,
		name TEXT,
		phone TEXT,
		location TEXT,
		family_size INTEGER,
		dietary_restrictions TEXT,
		created_at DATETIME
	);
	CREATE TABLE IF NOT EXISTS feedback (
		id TEXT PRIMARY KEY,
		donation_id TEXT,
		quality INTEGER,
		comments TEXT,
		created_at DATETIME,
		FOREIGN KEY(donation_id) REFERENCES donations(id)
	);`

	_, err = db.Exec(sqlStmt)
	return db, err
}

func serveTemplate(tmpl string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join("templates", tmpl))
	}
}

func createDonation(w http.ResponseWriter, r *http.Request) {
	var donation Donation
	if err := json.NewDecoder(r.Body).Decode(&donation); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	donation.ID = uuid.New().String()
	donation.CreatedAt = time.Now()

	_, err := db.Exec(`
		INSERT INTO donations (id, name, email, phone, meal_type, quantity, expiry_date, location, notes, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		donation.ID, donation.Name, donation.Email, donation.Phone, donation.MealType,
		donation.Quantity, donation.ExpiryDate, donation.Location, donation.Notes, donation.CreatedAt)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"donationId": donation.ID})
}

func createReceiver(w http.ResponseWriter, r *http.Request) {
	var receiver Receiver
	if err := json.NewDecoder(r.Body).Decode(&receiver); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	receiver.ID = uuid.New().String()
	receiver.CreatedAt = time.Now()

	restrictionsJSON, _ := json.Marshal(receiver.DietaryRestrictions)

	_, err := db.Exec(`
		INSERT INTO receivers (id, name, phone, location, family_size, dietary_restrictions, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		receiver.ID, receiver.Name, receiver.Phone, receiver.Location,
		receiver.FamilySize, string(restrictionsJSON), receiver.CreatedAt)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"receiverId": receiver.ID})
}

func getAvailableMeals(w http.ResponseWriter, r *http.Request) {
	location := r.URL.Query().Get("location")

	rows, err := db.Query(`
		SELECT id, meal_type, quantity, expiry_date, location
		FROM donations
		WHERE location LIKE ? AND quantity > 0 AND expiry_date > datetime('now')
		ORDER BY expiry_date ASC`,
		"%"+location+"%")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var meals []map[string]interface{}
	for rows.Next() {
		var meal struct {
			ID         string
			MealType   string
			Quantity   int
			ExpiryDate time.Time
			Location   string
		}
		if err := rows.Scan(&meal.ID, &meal.MealType, &meal.Quantity, &meal.ExpiryDate, &meal.Location); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		meals = append(meals, map[string]interface{}{
			"id":         meal.ID,
			"type":       meal.MealType,
			"quantity":   meal.Quantity,
			"expiryDate": meal.ExpiryDate,
			"location":   meal.Location,
		})
	}

	json.NewEncoder(w).Encode(meals)
}

func verifyMeal(w http.ResponseWriter, r *http.Request) {
	var data struct {
		DonationID string `json:"donationId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := db.Exec(`
		UPDATE donations
		SET quantity = quantity - 1
		WHERE id = ? AND quantity > 0`,
		data.DonationID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Invalid donation ID or no meals available", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func submitFeedback(w http.ResponseWriter, r *http.Request) {
	var feedback Feedback
	if err := json.NewDecoder(r.Body).Decode(&feedback); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	feedback.ID = uuid.New().String()
	feedback.CreatedAt = time.Now()

	_, err := db.Exec(`
		INSERT INTO feedback (id, donation_id, quality, comments, created_at)
		VALUES (?, ?, ?, ?, ?)`,
		feedback.ID, feedback.DonationID, feedback.Quality, feedback.Comments, feedback.CreatedAt)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
