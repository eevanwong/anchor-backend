package api

import (
	"anchor-backend/handlers"
	"net/http"
)

func RegisterRoutes() {
    http.HandleFunc("/api/dummy", handlers.DummyHandler)
}
