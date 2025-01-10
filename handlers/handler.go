package handlers

import (
	"anchor-backend/models"
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)

type LockRequest struct {
	RackID    uint   `json:"rack_id"`
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
	UserPhone string `json:"user_phone"`
}

type LockResponse struct {
	RackID      uint `json:"rack_id"`
	UserID      uint `json:"user_id"`
	LockSuccess bool `json:"lock_success"`
}

type UnlockRequest struct {
	RackID    uint   `json:"rack_id"`
	UserName  string `json:"user_name"`
	UserEmail string `json:"user_email"`
	UserPhone string `json:"user_phone"`
}

type UnlockResponse struct {
	RackID        uint `json:"rack_id"`
	UserID        uint `json:"user_id"`
	UnlockSuccess bool `json:"unlock_success"`
}

func LockHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var req LockRequest

	// Handle request body.
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "POST Lock: Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.RackID == 0 || req.UserName == "" || req.UserEmail == "" || req.UserPhone == "" {
		http.Error(w, "POST Lock: Missing field in request body", http.StatusBadRequest)
		return
	}

	// Fetch rack and check occupancy.
	rack, err := fetchRackByID(db, req.RackID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "POST Lock: Rack doesn't exist in DB", http.StatusBadRequest)
			return
		}
		http.Error(w, "POST Lock: Failed to retrieve rack from DB", http.StatusInternalServerError)
		return
	}
	if rack.CurrUserID != 0 {
		http.Error(w, "POST Lock: Rack occupied by another user", http.StatusBadRequest)
		return
	}

	// Create a new user if necessary.
	user, err := fetchUserByMetadata(db, req.UserName, req.UserEmail, req.UserPhone)
	if err != nil {
		http.Error(w, "POST Lock: Failed to verify user", http.StatusInternalServerError)
		return
	}

	// Update rack occupancy.
	res := db.Model(&models.Rack{}).Where("id = ?", rack.ID).Update("curr_user_id", user.ID)
	if res.Error != nil || res.RowsAffected == 0 {
		http.Error(w, "POST Lock: Failed updating rack in DB", http.StatusInternalServerError)
		return
	}

	response := LockResponse{
		RackID:      rack.ID,
		UserID:      user.ID,
		LockSuccess: true,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func UnlockHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var req UnlockRequest

	// Handle request body.
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "POST Unlock: Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.RackID == 0 || req.UserName == "" || req.UserEmail == "" || req.UserPhone == "" {
		http.Error(w, "POST Unlock: Missing field in request body", http.StatusBadRequest)
		return
	}

	// Fetch rack.
	rack, err := fetchRackByID(db, req.RackID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "POST Unlock: Rack doesn't exist in DB", http.StatusBadRequest)
			return
		}
		http.Error(w, "POST Unlock: Failed to retrieve rack from DB", http.StatusInternalServerError)
		return
	}

	// Fetch user with lock.
	user, err := fetchUserByID(db, rack.CurrUserID)
	if err != nil {
		http.Error(w, "POST Unlock: Failed to fetch user", http.StatusInternalServerError)
		return
	}

	// Authenticate unlocking user.
	if user.Name != req.UserName || user.Email != req.UserEmail || user.Phone != req.UserPhone {
		http.Error(w, "POST Unlock: Failed to authenticate user", http.StatusBadRequest)
		return
	}

	// TODO: 2FA if time and scope permits.

	// Update rack occupancy.
	res := db.Model(&models.Rack{}).Where("id = ?", rack.ID).Update("curr_user_id", nil)
	if res.Error != nil || res.RowsAffected == 0 {
		http.Error(w, "POST Unlock: Failed updating rack in DB", http.StatusInternalServerError)
		return
	}

	response := UnlockResponse{
		RackID:        rack.ID,
		UserID:        user.ID,
		UnlockSuccess: true,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func fetchUserByID(db *gorm.DB, userID uint) (*models.User, error) {
	// Check if user exists and/or is occupied.
	var user models.User
	res := db.Where("id = ?", userID).First(&user)
	if res.Error != nil {
		return nil, res.Error
	}

	return &user, nil
}

func fetchUserByMetadata(db *gorm.DB, userName string, userEmail string, userPhone string) (*models.User, error) {
	// Check if user exists and/or is occupied.
	var user models.User
	res := db.Where("name = ?", userName).Where("email = ?", userEmail).Where("phone = ?", userPhone).First(&user)
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		return nil, res.Error
	}
	if res.Error == nil {
		return &user, nil
	}
	newUser := models.User{
		Name:  userName,
		Email: userEmail,
		Phone: userPhone,
	}

	// Create new user.
	if err := db.Create(&newUser).Error; err != nil {
		return nil, err
	}

	return &newUser, nil
}

func fetchRackByID(db *gorm.DB, rackID uint) (*models.Rack, error) {
	// Check if user exists and/or is occupied.
	var rack models.Rack
	res := db.Where("id = ?", rackID).First(&rack)
	if res.Error != nil {
		return nil, res.Error
	}

	return &rack, nil
}
