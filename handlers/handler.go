package handlers

import (
	"anchor-backend/models"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

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

type GetRacksResponse struct {
	Racks []RackDetails `json:"rack_details"`
}

type RackDetails struct {
	RackID      uint   `json:"rack_id"`
	UserID      uint   `json:"user_id"`
	UserName    string `json:"user_name"`
	UserEmail   string `json:"user_email"`
	UserPhone   string `json:"user_phone"`
	LastUpdated string `json:"last_updated"`
}

func decryptData(encryptedData string, ivBase64 string, aesKey []byte) (string, error) {
	iv, err := base64.StdEncoding.DecodeString(ivBase64)
	if err != nil {
		return "", fmt.Errorf("failed to decode IV: %v", err)
	}

	encrypted, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return "", fmt.Errorf("failed to decode encrypted data: %v", err)
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %v", err)
	}

	decrypted := make([]byte, len(encrypted))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(decrypted, encrypted)

	decrypted, err = unpad(decrypted)
	if err != nil {
		return "", fmt.Errorf("failed to unpad decrypted data: %v", err)
	}

	return string(decrypted), nil
}

func unpad(data []byte) ([]byte, error) {
	padding := data[len(data)-1]
	if int(padding) > len(data) {
		return nil, fmt.Errorf("invalid padding size")
	}
	return data[:len(data)-int(padding)], nil
}

func LockHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var req = LockRequest{RackID: 1, UserName: "Erick", UserEmail: "erick@gmail.com", UserPhone: "1231231234"}

	// Decrypt data.
	aesKey := []byte("12345678901234567890123456789012")
	var requestBody struct {
		Data string `json:"data"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "POST Lock: Invalid JSON in request body", http.StatusBadRequest)
		return
	}

	requestSections := strings.Split(string(requestBody.Data), ":")
	if len(requestSections) != 2 {
		http.Error(w, "POST Lock: Incorrect section count in request payload", http.StatusBadRequest)
		return
	}

	ivBase64 := requestSections[0]
	encryptedData := requestSections[1]
	decrypted, err := decryptData(encryptedData, ivBase64, aesKey)
	if err != nil {
		http.Error(w, "POST Lock: Invalid encryption of request payload", http.StatusBadRequest)
		return
	}
	log.Printf("Decrypted data: %s", decrypted)

	// Decode decrypted JSON.
	err = json.Unmarshal([]byte(decrypted), &req)
	if err != nil {
		http.Error(w, "POST Lock: Invalid JSON in decrypted payload", http.StatusBadRequest)
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

	NotifyWebSocketClients(fmt.Sprintf("{\"rack_id\": %d, \"user_id\": %d, \"user_name\": \"%s\", \"user_email\": \"%s\", \"user_phone\": \"%s\", \"last_updated\": \"%s\", \"action\": \"lock\"}", rack.ID, user.ID, user.Name, user.Email, user.Phone, rack.UpdatedAt.Format("2006-01-02 15:04:05")))

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
	var req = UnlockRequest{RackID: 1, UserName: "Erick", UserEmail: "erick@gmail.com", UserPhone: "1231231234"}

	// Decrypt data.
	aesKey := []byte("12345678901234567890123456789012")
	var requestBody struct {
		Data string `json:"data"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "POST Unlock: Invalid JSON in request body", http.StatusBadRequest)
		return
	}

	requestSections := strings.Split(string(requestBody.Data), ":")
	if len(requestSections) != 2 {
		http.Error(w, "POST Unlock: Incorrect section count in request payload", http.StatusBadRequest)
		return
	}

	ivBase64 := requestSections[0]
	encryptedData := requestSections[1]
	decrypted, err := decryptData(encryptedData, ivBase64, aesKey)
	if err != nil {
		http.Error(w, "POST Unlock: Invalid encryption of request payload", http.StatusBadRequest)
		return
	}
	log.Printf("Decrypted data: %s", decrypted)

	// Decode decrypted JSON.
	err = json.Unmarshal([]byte(decrypted), &req)
	if err != nil {
		http.Error(w, "POST Unlock: Invalid JSON in decrypted payload", http.StatusBadRequest)
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

	// Update rack occupancy.
	res := db.Model(&models.Rack{}).Where("id = ?", rack.ID).Update("curr_user_id", 0)
	if res.Error != nil || res.RowsAffected == 0 {
		http.Error(w, "POST Unlock: Failed updating rack in DB", http.StatusInternalServerError)
		return
	}

	NotifyWebSocketClients(fmt.Sprintf("{\"rack_id\": %d, \"user_id\": %d, \"user_name\": \"\", \"user_email\": \"\", \"user_phone\": \"\", \"updated_at\": \"\", \"action\": \"unlock\"}", rack.ID, 0))

	response := UnlockResponse{
		RackID:        rack.ID,
		UserID:        user.ID,
		UnlockSuccess: true,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Fetch all racks to display on the frontend.
func GetRacksHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var racks []models.Rack
	if err := db.Find(&racks).Error; err != nil {
		http.Error(w, "Failed to fetch racks", http.StatusInternalServerError)
		return
	}

	var allRackDetails []RackDetails
	for _, rack := range racks {
		if rack.CurrUserID == 0 {
			allRackDetails = append(allRackDetails, RackDetails{RackID: rack.ID})
			continue
		}
		user, err := fetchUserByID(db, rack.CurrUserID)
		if err != nil {
			http.Error(w, "Get Racks: Failed to fetch user", http.StatusInternalServerError)
			return
		}
		rackDetails := RackDetails{
			RackID:      rack.ID,
			UserID:      user.ID,
			UserName:    user.Name,
			UserEmail:   user.Email,
			UserPhone:   user.Phone,
			LastUpdated: rack.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		allRackDetails = append(allRackDetails, rackDetails)
	}

	var getRackResponse = GetRacksResponse{Racks: allRackDetails}

	json.NewEncoder(w).Encode(getRackResponse)
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
