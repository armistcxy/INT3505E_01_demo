package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"book-service/internal/handler"
	"book-service/internal/repository"
	"book-service/internal/service"
	"book-service/pkg/database"
	"book-service/pkg/middlewares"
)

func main() {
	// Database configuration
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "postgres"
	}
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "bookdb"
	}

	// Connect to database
	log.Println("Connecting to database...")
	db, err := database.NewConnection(dbUser, dbPassword, dbHost, dbPort, dbName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize schema
	log.Println("Initializing schema...")
	if err := database.InitializeSchema(db); err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}

	// Initialize repository, service, and handler
	repo := repository.NewBookRepository(db)
	svc := service.NewBookService(repo)
	bookHandler := handler.NewBookHandler(svc)

	// Setup routes
	r := mux.NewRouter()

	// Rate limiting middleware
	rateLimitMiddleware := middlewares.NewRateLimitMiddleware(50)

	// Book routes
	r.HandleFunc("/api/books", bookHandler.CreateBook).Methods("POST")
	r.HandleFunc("/api/books", bookHandler.GetAllBooks).Methods("GET")
	r.Handle("/api/books/{id}", rateLimitMiddleware(http.HandlerFunc(bookHandler.GetBook))).Methods("GET")
	r.HandleFunc("/api/books/{id}", bookHandler.UpdateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", bookHandler.DeleteBook).Methods("DELETE")

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status":"ok"}`)
	}).Methods("GET")

	// Prometheus metrics
	r.Handle("/metrics", promhttp.Handler())

	// Start server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
