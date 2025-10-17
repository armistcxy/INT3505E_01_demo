package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type Article struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	Updated time.Time `json:"updated"`
}

var article = Article{
	ID:      1,
	Title:   "REST API Cacheable Demo",
	Updated: time.Now(),
}

func getArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.Header().Set("Cache-Control", "public, max-age=10")

	json.NewEncoder(w).Encode(article)
}

func main() {
	http.HandleFunc("/article", getArticle)
	http.ListenAndServe(":8080", nil)
}
