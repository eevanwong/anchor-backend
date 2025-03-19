package api

import (
	"anchor-backend/handlers"
	"net/http"

	"gorm.io/gorm"
)

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from frontend
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	handlers.HandleWebSocket(w, r)
}

func RegisterRoutes(db *gorm.DB) {
	http.HandleFunc("/api/lock", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.LockHandler(w, r, db)
	}))
	http.HandleFunc("/api/unlock", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.UnlockHandler(w, r, db)
	}))
	http.HandleFunc("/api/racks", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlers.GetRacksHandler(w, r, db)
	}))

	http.HandleFunc("/ws", wsHandler)
}
