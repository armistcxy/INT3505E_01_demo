package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"go.uber.org/ratelimit"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	rl := ratelimit.New(5)

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello! Time: %s\n", time.Now().Format(time.RFC3339))
	})

	mux.HandleFunc("GET /books", HandleGetBooks)
	mux.HandleFunc("POST /books", HandleAddBooks)

	handler := loggingMiddleware(rateLimitMiddleware(mux, rl))

	log.Printf("Listening on :%s ...", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(body))

		log.Printf("[LOG] %s %s Body: %s", r.Method, r.URL.Path, string(body))
		next.ServeHTTP(w, r)
	})
}

func rateLimitMiddleware(next http.Handler, rl ratelimit.Limiter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rl.Take()
		next.ServeHTTP(w, r)
	})
}

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

func HandleGetBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func HandleAddBooks(w http.ResponseWriter, r *http.Request) {
	var b Book
	json.NewDecoder(r.Body).Decode(&b)

	books = append(books, b)
	json.NewEncoder(w).Encode(map[string]string{
		"msg": fmt.Sprintf("Add book %s successfully", b.Title),
	})
}

var books = []Book{
	{ID: 1, Title: "The Pragmatic Programmer", Author: "Andrew Hunt"},
	{ID: 2, Title: "Clean Code", Author: "Robert C. Martin"},
	{ID: 3, Title: "Designing Data-Intensive Applications", Author: "Martin Kleppmann"},
}
