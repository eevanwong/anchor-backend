package main

import (
	"anchor-backend/api"
	"anchor-backend/database"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// invoked before main()
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func main() {
	port, exists := os.LookupEnv("SERVER_PORT")

	if exists {
		log.Printf("Starting server at port %s \n", port)
	}

	db, err := database.MigrateAndSeedDatabase()
	if err != nil {
		log.Fatal("Failed to initialize database: ", err)
	}

	api.RegisterRoutes(db)

	http.HandleFunc("/", handler)
	http.ListenAndServe(":"+port, nil)
}
