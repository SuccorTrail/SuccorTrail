package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/SuccorTrail/SuccorTrail/internal/model"
	"github.com/SuccorTrail/SuccorTrail/internal/repository"
	"github.com/SuccorTrail/SuccorTrail/internal/util"
	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	UserRepo repository.UserRepository
}

// RenderLoginForm renders the login page
func (h *AuthHandler) RenderLoginForm(w http.ResponseWriter, r *http.Request) {
	util.RenderTemplate(w, "login.html", nil)
}

// RenderSignupForm renders the signup page
func (h *AuthHandler) RenderSignupForm(w http.ResponseWriter, r *http.Request) {
	util.RenderTemplate(w, "signup.html", nil)
}

// RenderLandingPage renders the landing page
func (h *AuthHandler) RenderLandingPage(w http.ResponseWriter, r *http.Request) {
	util.RenderTemplate(w, "landing.html", nil)
}

// SignUp handles user registration
func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	// Parse JSON request
	var request struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		UserType string `json:"userType"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		logrus.WithError(err).Error("Error decoding JSON request")
		sendJSONResponse(w, http.StatusBadRequest, false, "Invalid request format")
		return
	}

	// Validate required fields
	if request.Name == "" || request.Email == "" || request.Password == "" || request.UserType == "" {
		sendJSONResponse(w, http.StatusBadRequest, false, "Missing required fields")
		return
	}

	// Check if user already exists
	exists, err := h.UserRepo.UserExists(request.Email)
	if err != nil {
		logrus.WithError(err).Error("Error checking if user exists")
		sendJSONResponse(w, http.StatusInternalServerError, false, "Error creating account")
		return
	}

	if exists {
		sendJSONResponse(w, http.StatusConflict, false, "Email already registered")
		return
	}

	// Hash password (in a real app)
	// hashedPassword := util.HashPassword(request.Password)

	// Create user
	user := &model.User{
		ID:        util.GenerateUUID(),
		Name:      request.Name,
		Email:     request.Email,
		Password:  request.Password,
		UserType:  request.UserType,
		CreatedAt: time.Now(),
	}

	// Save user to database
	err = h.UserRepo.Create(user)
	if err != nil {
		logrus.WithError(err).Error("Error creating user")
		sendJSONResponse(w, http.StatusInternalServerError, false, "Error creating account")
		return
	}

	// Return success response
	sendJSONResponse(w, http.StatusOK, true, "Account created successfully!")
}

// Login handles user authentication
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Parse JSON request
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		logrus.WithError(err).Error("Error decoding JSON request")
		sendJSONResponse(w, http.StatusBadRequest, false, "Invalid request format")
		return
	}

	// Validate required fields
	if request.Email == "" || request.Password == "" {
		sendJSONResponse(w, http.StatusBadRequest, false, "Email and password are required")
		return
	}

	// Get user by email
	user, err := h.UserRepo.GetByEmail(request.Email)
	if err != nil {
		logrus.WithError(err).Error("Error retrieving user")
		sendJSONResponse(w, http.StatusUnauthorized, false, "Invalid credentials")
		return
	}

	// In a real app, verify password
	// if !util.VerifyPassword(password, user.Password) {
	if request.Password != user.Password {
		sendJSONResponse(w, http.StatusUnauthorized, false, "Invalid credentials")
		return
	}

	// Determine redirect URL based on user type
	var redirectURL string
	switch user.UserType {
	case "donor":
		redirectURL = "/donor"
	case "organization":
		redirectURL = "/organization"
	case "recipient":
		redirectURL = "/receiver"
	default:
		redirectURL = "/dashboard"
	}

	// Return success response with redirect URL
	response := map[string]interface{}{
		"success":     true,
		"message":     "Login successful",
		"redirectURL": redirectURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Helper function to send JSON responses
func sendJSONResponse(w http.ResponseWriter, statusCode int, success bool, message string) {
	response := map[string]interface{}{
		"success": success,
		"message": message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
