package api

import (
	"anchor-backend/handlers"
	"net/http"

	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB) {
	http.HandleFunc("/api/lock", func(w http.ResponseWriter, r *http.Request) {
		handlers.LockHandler(w, r, db)
	})
	http.HandleFunc("/api/unlock", func(w http.ResponseWriter, r *http.Request) {
		handlers.UnlockHandler(w, r, db)
	})
}
