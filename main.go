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
	UpdatedAt     time.Time `json:"updatedAt"`
	Status        string    `json:"status"`
}

type Receiver struct {
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	Phone              string   `json:"phone"`
	Location           string   `json:"location"`
	FamilySize         int      `json:"familySize"`
	DietaryRestrictions []string `json:"dietaryRestrictions"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
	VerificationStatus string    `json:"verificationStatus"`
}

type Feedback struct {
	ID        string    `json:"id"`
	DistributionID string    `json:"distributionId"`
	Quality   int       `json:"quality"`
	Comments  string    `json:"comments"`
	CreatedAt time.Time `json:"createdAt"`
}

type MealDistribution struct {
	ID             string    `json:"id"`
	DonationID     string    `json:"donationId"`
	ReceiverID     string    `json:"receiverId"`
	Quantity       int       `json:"quantity"`
	DistributionDate time.Time `json:"distributionDate"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"createdAt"`
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
	r.HandleFunc("/api/distributions", createMealDistribution).Methods("POST")

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

	// Enable foreign key constraints
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return nil, err
	}

	// Create tables with proper constraints and indices
	sqlStmt := `
	-- Donations table
	CREATE TABLE IF NOT EXISTS donations (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT NOT NULL,
		phone TEXT NOT NULL,
		meal_type TEXT NOT NULL,
		quantity INTEGER NOT NULL CHECK (quantity >= 0),
		expiry_date DATETIME NOT NULL,
		location TEXT NOT NULL,
		notes TEXT,
		status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'completed', 'expired')),
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);

	-- Create indices for frequently queried columns
	CREATE INDEX IF NOT EXISTS idx_donations_location ON donations(location);
	CREATE INDEX IF NOT EXISTS idx_donations_expiry ON donations(expiry_date);
	CREATE INDEX IF NOT EXISTS idx_donations_status ON donations(status);

	-- Receivers table
	CREATE TABLE IF NOT EXISTS receivers (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		phone TEXT NOT NULL UNIQUE,
		location TEXT NOT NULL,
		family_size INTEGER NOT NULL CHECK (family_size > 0),
		dietary_restrictions TEXT,
		verification_status TEXT NOT NULL DEFAULT 'pending' CHECK (verification_status IN ('pending', 'verified', 'rejected')),
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);

	-- Create indices for receivers
	CREATE INDEX IF NOT EXISTS idx_receivers_location ON receivers(location);
	CREATE INDEX IF NOT EXISTS idx_receivers_phone ON receivers(phone);

	-- Meal distributions table (tracks actual distributions)
	CREATE TABLE IF NOT EXISTS meal_distributions (
		id TEXT PRIMARY KEY,
		donation_id TEXT NOT NULL,
		receiver_id TEXT NOT NULL,
		quantity INTEGER NOT NULL CHECK (quantity > 0),
		distribution_date DATETIME NOT NULL,
		status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'completed', 'cancelled')),
		created_at DATETIME NOT NULL,
		FOREIGN KEY (donation_id) REFERENCES donations(id),
		FOREIGN KEY (receiver_id) REFERENCES receivers(id)
	);

	-- Create indices for meal distributions
	CREATE INDEX IF NOT EXISTS idx_distributions_donation ON meal_distributions(donation_id);
	CREATE INDEX IF NOT EXISTS idx_distributions_receiver ON meal_distributions(receiver_id);
	CREATE INDEX IF NOT EXISTS idx_distributions_date ON meal_distributions(distribution_date);

	-- Feedback table
	CREATE TABLE IF NOT EXISTS feedback (
		id TEXT PRIMARY KEY,
		distribution_id TEXT NOT NULL,
		quality_rating INTEGER NOT NULL CHECK (quality_rating BETWEEN 1 AND 5),
		comments TEXT,
		created_at DATETIME NOT NULL,
		FOREIGN KEY (distribution_id) REFERENCES meal_distributions(id)
	);

	-- Create index for feedback
	CREATE INDEX IF NOT EXISTS idx_feedback_distribution ON feedback(distribution_id);

	-- Notifications table
	CREATE TABLE IF NOT EXISTS notifications (
		id TEXT PRIMARY KEY,
		receiver_id TEXT NOT NULL,
		message TEXT NOT NULL,
		type TEXT NOT NULL CHECK (type IN ('meal_available', 'distribution_reminder', 'feedback_request', 'system')),
		status TEXT NOT NULL DEFAULT 'unread' CHECK (status IN ('unread', 'read')),
		created_at DATETIME NOT NULL,
		read_at DATETIME,
		FOREIGN KEY (receiver_id) REFERENCES receivers(id)
	);

	-- Create indices for notifications
	CREATE INDEX IF NOT EXISTS idx_notifications_receiver ON notifications(receiver_id);
	CREATE INDEX IF NOT EXISTS idx_notifications_status ON notifications(status);

	-- Audit log table
	CREATE TABLE IF NOT EXISTS audit_log (
		id TEXT PRIMARY KEY,
		table_name TEXT NOT NULL,
		record_id TEXT NOT NULL,
		action TEXT NOT NULL CHECK (action IN ('insert', 'update', 'delete')),
		old_value TEXT,
		new_value TEXT,
		created_at DATETIME NOT NULL,
		created_by TEXT NOT NULL
	);

	-- Create indices for audit log
	CREATE INDEX IF NOT EXISTS idx_audit_table_record ON audit_log(table_name, record_id);
	CREATE INDEX IF NOT EXISTS idx_audit_created ON audit_log(created_at);
	`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, err
	}

	// Create triggers for updated_at timestamps
	triggers := `
	-- Trigger for donations updated_at
	CREATE TRIGGER IF NOT EXISTS trg_donations_updated_at 
	AFTER UPDATE ON donations
	BEGIN
		UPDATE donations SET updated_at = DATETIME('now') WHERE id = NEW.id;
	END;

	-- Trigger for receivers updated_at
	CREATE TRIGGER IF NOT EXISTS trg_receivers_updated_at 
	AFTER UPDATE ON receivers
	BEGIN
		UPDATE receivers SET updated_at = DATETIME('now') WHERE id = NEW.id;
	END;
	`

	_, err = db.Exec(triggers)
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
		log.Printf("Error decoding donation: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Log the received data
	log.Printf("Received donation data: %+v", donation)

	donation.ID = uuid.New().String()
	donation.CreatedAt = time.Now()
	donation.UpdatedAt = time.Now()
	donation.Status = "active"

	// Log the data being inserted
	log.Printf("Inserting donation with ID: %s, Status: %s", donation.ID, donation.Status)

	result, err := db.Exec(`
		INSERT INTO donations (id, name, email, phone, meal_type, quantity, expiry_date, location, notes, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		donation.ID, donation.Name, donation.Email, donation.Phone, donation.MealType,
		donation.Quantity, donation.ExpiryDate, donation.Location, donation.Notes, donation.Status, donation.CreatedAt, donation.UpdatedAt)

	if err != nil {
		log.Printf("Error inserting donation: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
	} else {
		log.Printf("Rows affected: %d", rowsAffected)
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
	receiver.UpdatedAt = time.Now()
	receiver.VerificationStatus = "pending"

	restrictionsJSON, _ := json.Marshal(receiver.DietaryRestrictions)

	_, err := db.Exec(`
		INSERT INTO receivers (id, name, phone, location, family_size, dietary_restrictions, verification_status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		receiver.ID, receiver.Name, receiver.Phone, receiver.Location,
		receiver.FamilySize, string(restrictionsJSON), receiver.VerificationStatus, receiver.CreatedAt, receiver.UpdatedAt)

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
		WHERE location LIKE ? AND quantity > 0 AND expiry_date > datetime('now') AND status = 'active'
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
		SET quantity = quantity - 1, status = 'completed'
		WHERE id = ? AND quantity > 0 AND status = 'active'`,
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
		INSERT INTO feedback (id, distribution_id, quality_rating, comments, created_at)
		VALUES (?, ?, ?, ?, ?)`,
		feedback.ID, feedback.DistributionID, feedback.Quality, feedback.Comments, feedback.CreatedAt)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func createMealDistribution(w http.ResponseWriter, r *http.Request) {
	var distribution MealDistribution
	if err := json.NewDecoder(r.Body).Decode(&distribution); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	distribution.ID = uuid.New().String()
	distribution.CreatedAt = time.Now()

	_, err := db.Exec(`
		INSERT INTO meal_distributions (id, donation_id, receiver_id, quantity, distribution_date, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		distribution.ID, distribution.DonationID, distribution.ReceiverID, distribution.Quantity, distribution.DistributionDate, distribution.Status, distribution.CreatedAt)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
