package handlers

import (
	"anchor-backend/models"
	"encoding/json"
	"net/http"
	"regexp"

	"gorm.io/gorm"
)

type LockRequest struct {
	RackID      string `json:"rack_id"`
	UserID      string `json:"user_id"`
	UserName    string `json:"user_name"`
	UserContact string `json:"user_contact"`
}

type LockResponse struct {
	RackID      string `json:"rack_id"`
	UserID      string `json:"user_id"`
	LockSuccess bool   `json:"lock_success"`
}

func LockHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var req LockRequest

	// Handle request body.
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "POST Lock: Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.RackID == "" || req.UserID == "" {
		http.Error(w, "POST Lock: Missing rack_id and/or user_id", http.StatusBadRequest)
		return
	}

	// Check if rack exists and/or is occupied.
	var rack models.Rack
	res := db.Where("id = ?", req.RackID).First(&rack)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
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
	user, err := createUser(db, req.UserName, req.UserContact)
	if err != nil {
		http.Error(w, "POST Lock: Failed to verify user", http.StatusInternalServerError)
		return
	}

	// Update rack occupancy.
	res = db.Model(&models.Rack{}).Where("id = ?", rack.ID).Update("curr_user_id", user.ID)
	if res.Error != nil || res.RowsAffected == 0 {
		http.Error(w, "POST Lock: Failed updating rack in DB", http.StatusInternalServerError)
		return
	}

	response := LockResponse{
		RackID:      req.RackID,
		UserID:      req.UserID,
		LockSuccess: true,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func createUser(db *gorm.DB, userName string, userContact string) (*models.User, error) {
	// Check if user exists and/or is occupied.
	var user models.User
	res := db.Where("name = ?", userName).First(&user)
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		return nil, res.Error
	}
	if res.Error == nil {
		return &user, nil
	}

	contactType := "unknown"
	if isEmail(userContact) {
		contactType = "phone"
	} else if isPhoneNumber(userContact) {
		contactType = "email"
	}
	newUser := models.User{
		Name:        userName,
		Contact:     userContact,
		ContactType: contactType,
	}

	if err := db.Create(&newUser).Error; err != nil {
		return nil, err
	}

	return &newUser, nil
}

func isEmail(s string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(s)
}

func isPhoneNumber(s string) bool {
	phoneRegex := `^\+?[0-9]{1,3}?[-. ]?[0-9]{1,4}[-. ]?[0-9]{1,4}[-. ]?[0-9]{1,9}$`
	re := regexp.MustCompile(phoneRegex)
	return re.MatchString(s)
}

func jsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
