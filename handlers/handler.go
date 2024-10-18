package handlers

import (
	"encoding/json"
	"net/http"
)

func DummyHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        handleGetResource(w, r)
    } else {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func handleGetResource(w http.ResponseWriter, r *http.Request) {
    response := map[string]string{"message": "Hello, World!"}
    jsonResponse(w, response, http.StatusOK)
}

func jsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(data)
}
