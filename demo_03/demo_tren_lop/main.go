package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type ErrorDetail struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

type ErrorPayload struct {
	Error ErrorDetail `json:"error"`
}

func writeError(w http.ResponseWriter, code int, status, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(ErrorPayload{Error: ErrorDetail{
		Code:    code,
		Message: message,
		Status:  status,
	}})
}

func writeJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

type Movie struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Genre string `json:"genre"`
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Verified bool   `json:"verified"`
}

type Rating struct {
	ID      int    `json:"id"`
	UserID  int    `json:"userId"`
	MovieID int    `json:"movieId"`
	Score   int    `json:"score"`
	Comment string `json:"comment,omitempty"`
}

var (
	movies       = map[int]Movie{}
	users        = map[int]User{}
	ratings      = map[int]Rating{}
	nextMovieID  = 1
	nextUserID   = 1
	nextRatingID = 1
)

func seedData() {
	m1 := Movie{ID: nextMovieID, Name: "Inception", Genre: "sci-fi"}
	nextMovieID++
	m2 := Movie{ID: nextMovieID, Name: "Titanic", Genre: "romance"}
	nextMovieID++
	movies[m1.ID] = m1
	movies[m2.ID] = m2

	u1 := User{ID: nextUserID, Name: "Alice", Email: "alice@example.com", Verified: false}
	nextUserID++
	u2 := User{ID: nextUserID, Name: "Bob", Email: "bob@example.com", Verified: true}
	nextUserID++
	users[u1.ID] = u1
	users[u2.ID] = u2

	r1 := Rating{ID: nextRatingID, UserID: u2.ID, MovieID: m1.ID, Score: 5, Comment: "Mind-blowing!"}
	nextRatingID++
	ratings[r1.ID] = r1
}

func main() {
	seedData()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/v1", func(v1 chi.Router) {
		// Movies CRUD
		v1.Get("/movies", listMovies)
		v1.Post("/movies", createMovie)
		v1.Get("/movies/{id}", getMovie)
		v1.Patch("/movies/{id}", updateMovie)
		v1.Delete("/movies/{id}", deleteMovie)

		// Custom verb on collection: import
		v1.Post("/movies:import", importMovies)

		// Users
		v1.Get("/users", listUsers)
		v1.Post("/users", createUser)
		v1.Get("/users/{id}", getUser)

		// Custom verb on instance: sendVerificationEmail
		v1.Post("/users/{id}:sendVerificationEmail", sendVerificationEmail)

		// Nested within users: ratings
		v1.Get("/users/{id}/ratings", listRatingsOfUser)
		v1.Post("/users/{id}/ratings", createRatingForUser)

		// Ratings direct access
		v1.Get("/ratings/{id}", getRating)
		v1.Patch("/ratings/{id}", updateRating)
		v1.Delete("/ratings/{id}", deleteRating)
	})

	addr := ":8080"
	log.Printf("HTTP server listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}

// ========= Movies =========

func listMovies(w http.ResponseWriter, r *http.Request) {
	search := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("search")))
	genre := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("genre")))
	out := make([]Movie, 0, len(movies))
	for _, m := range movies {
		if search != "" && !strings.Contains(strings.ToLower(m.Name), search) {
			continue
		}
		if genre != "" && strings.ToLower(m.Genre) != genre {
			continue
		}
		out = append(out, m)
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":  out,
		"count": len(out),
	})
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ARGUMENT", fmt.Sprintf("Invalid movie id '%s'.", idStr))
		return
	}
	m, ok := movies[id]
	if !ok {
		writeError(w, http.StatusNotFound, "NOT_FOUND", fmt.Sprintf("Movie with name 'movies/%d' not found.", id))
		return
	}
	writeJSON(w, http.StatusOK, m)
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name  string `json:"name"`
		Genre string `json:"genre"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ARGUMENT", "Invalid JSON body.")
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	req.Genre = strings.TrimSpace(req.Genre)
	if req.Name == "" || req.Genre == "" {
		writeError(w, http.StatusBadRequest, "INVALID_ARGUMENT", "Both 'name' and 'genre' are required.")
		return
	}
	m := Movie{
		ID:    nextMovieID,
		Name:  req.Name,
		Genre: strings.ToLower(req.Genre),
	}
	nextMovieID++
	movies[m.ID] = m
	writeJSON(w, http.StatusCreated, m)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ARGUMENT", fmt.Sprintf("Invalid movie id '%s'.", idStr))
		return
	}
	m, ok := movies[id]
	if !ok {
		writeError(w, http.StatusNotFound, "NOT_FOUND", fmt.Sprintf("Movie with name 'movies/%d' not found.", id))
		return
	}
	var patch struct {
		Name  *string `json:"name,omitempty"`
		Genre *string `json:"genre,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ARGUMENT", "Invalid JSON body.")
		return
	}
	if patch.Name != nil {
		m.Name = strings.TrimSpace(*patch.Name)
	}
	if patch.Genre != nil {
		m.Genre = strings.ToLower(strings.TrimSpace(*patch.Genre))
	}
	movies[id] = m
	writeJSON(w, http.StatusOK, m)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ARGUMENT", fmt.Sprintf("Invalid movie id '%s'.", idStr))
		return
	}
	if _, ok := movies[id]; !ok {
		writeError(w, http.StatusNotFound, "NOT_FOUND", fmt.Sprintf("Movie with name 'movies/%d' not found.", id))
		return
	}
	delete(movies, id)
	w.WriteHeader(http.StatusNoContent)
}

// importMovies supports JSON array of movies or CSV file upload via multipart/form-data (file field "file")
func importMovies(w http.ResponseWriter, r *http.Request) {
	ct := r.Header.Get("Content-Type")
	if strings.HasPrefix(ct, "application/json") {
		var in []struct {
			Name  string `json:"name"`
			Genre string `json:"genre"`
		}
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_ARGUMENT", "Invalid JSON body.")
			return
		}
		imported := make([]Movie, 0, len(in))
		for _, row := range in {
			name := strings.TrimSpace(row.Name)
			genre := strings.ToLower(strings.TrimSpace(row.Genre))
			if name == "" || genre == "" {
				writeError(w, http.StatusBadRequest, "INVALID_ARGUMENT", "All rows must have name and genre.")
				return
			}
			m := Movie{ID: nextMovieID, Name: name, Genre: genre}
			nextMovieID++
			movies[m.ID] = m
			imported = append(imported, m)
		}
		writeJSON(w, http.StatusCreated, map[string]interface{}{
			"data":  imported,
			"count": len(imported),
		})
		return
	}

	if strings.HasPrefix(ct, "multipart/form-data") {
		if err := r.ParseMultipartForm(5 << 20); err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_ARGUMENT", "Failed to parse multipart form.")
			return
		}
		file, _, err := r.FormFile("file")
		if err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_ARGUMENT", "File field 'file' is required.")
			return
		}
		defer file.Close()
		reader := csv.NewReader(file)
		imported := []Movie{}
		for {
			rec, err := reader.Read()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				writeError(w, http.StatusBadRequest, "INVALID_ARGUMENT", "Invalid CSV file.")
				return
			}
			if len(rec) < 2 {
				writeError(w, http.StatusBadRequest, "INVALID_ARGUMENT", "CSV must have at least two columns: name,genre.")
				return
			}
			name := strings.TrimSpace(rec[0])
			genre := strings.ToLower(strings.TrimSpace(rec[1]))
			if name == "" || genre == "" {
				writeError(w, http.StatusBadRequest, "INVALID_ARGUMENT", "CSV rows must have non-empty name and genre.")
				return
			}
			m := Movie{ID: nextMovieID, Name: name, Genre: genre}
			nextMovieID++
			movies[m.ID] = m
			imported = append(imported, m)
		}
		writeJSON(w, http.StatusCreated, map[string]interface{}{
			"data":  imported,
			"count": len(imported),
		})
		return
	}

	writeError(w, http.StatusUnsupportedMediaType, "INVALID_ARGUMENT", "Content-Type must be application/json or multipart/form-data.")
}

// ========= Users =========

func listUsers(w http.ResponseWriter, r *http.Request) {
	out := make([]User, 0, len(users))
	for _, u := range users {
		out = append(out, u)
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":  out,
		"count": len(out),
	})
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ARGUMENT", "Invalid JSON body.")
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	if req.Name == "" || req.Email == "" {
		writeError(w, http.StatusBadRequest, "INVALID_ARGUMENT", "Both 'name' and 'email' are required.")
		return
	}
	u := User{
		ID:       nextUserID,
		Name:     req.Name,
		Email:    req.Email,
		Verified: false,
	}
	nextUserID++
	users[u.ID] = u
	writeJSON(w, http.StatusCreated, u)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ARGUMENT", fmt.Sprintf("Invalid user id '%s'.", idStr))
		return
	}
	u, ok := users[id]
	if !ok {
		writeError(w, http.StatusNotFound, "NOT_FOUND", fmt.Sprintf("User with name 'users/%d' not found.", id))
		return
	}
	writeJSON(w, http.StatusOK, u)
}

func sendVerificationEmail(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ARGUMENT", fmt.Sprintf("Invalid user id '%s'.", idStr))
		return
	}
	u, ok := users[id]
	if !ok {
		writeError(w, http.StatusNotFound, "NOT_FOUND", fmt.Sprintf("User with name 'users/%d' not found.", id))
		return
	}
	// Simulate sending email
	time.Sleep(100 * time.Millisecond)
	u.Verified = true
	users[id] = u
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Verification email sent.",
		"user":    u,
	})
}

// ========= Ratings (Nested and Direct) =========

// listRatingsOfUser handles GET /v1/users/{id}/ratings
func listRatingsOfUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ARGUMENT", fmt.Sprintf("Invalid user id '%s'.", userIDStr))
		return
	}

	if _, ok := users[userID]; !ok {
		writeError(w, http.StatusNotFound, "NOT_FOUND", fmt.Sprintf("User with name 'users/%d' not found.", userID))
		return
	}

	out := make([]Rating, 0)
	for _, rating := range ratings {
		if rating.UserID == userID {
			out = append(out, rating)
		}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":  out,
		"count": len(out),
	})
}

// createRatingForUser handles POST /v1/users/{id}/ratings
func createRatingForUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ARGUMENT", fmt.Sprintf("Invalid user id '%s'.", userIDStr))
		return
	}

	if _, ok := users[userID]; !ok {
		writeError(w, http.StatusNotFound, "NOT_FOUND", fmt.Sprintf("User with name 'users/%d' not found.", userID))
		return
	}

	var req struct {
		MovieID int    `json:"movieId"`
		Score   int    `json:"score"`
		Comment string `json:"comment"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ARGUMENT", "Invalid JSON body.")
		return
	}

	if _, ok := movies[req.MovieID]; !ok {
		writeError(w, http.StatusBadRequest, "INVALID_ARGUMENT", fmt.Sprintf("Movie with id '%d' does not exist.", req.MovieID))
		return
	}

	if req.Score < 1 || req.Score > 5 {
		writeError(w, http.StatusBadRequest, "INVALID_ARGUMENT", "Score must be between 1 and 5.")
		return
	}

	newRating := Rating{
		ID:      nextRatingID,
		UserID:  userID,
		MovieID: req.MovieID,
		Score:   req.Score,
		Comment: strings.TrimSpace(req.Comment),
	}
	nextRatingID++
	ratings[newRating.ID] = newRating

	writeJSON(w, http.StatusCreated, newRating)
}

// getRating handles GET /v1/ratings/{id}
func getRating(w http.ResponseWriter, r *http.Request) {
	ratingIDStr := chi.URLParam(r, "id")
	ratingID, err := strconv.Atoi(ratingIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ARGUMENT", fmt.Sprintf("Invalid rating id '%s'.", ratingIDStr))
		return
	}

	rating, ok := ratings[ratingID]
	if !ok {
		writeError(w, http.StatusNotFound, "NOT_FOUND", fmt.Sprintf("Rating with name 'ratings/%d' not found.", ratingID))
		return
	}

	writeJSON(w, http.StatusOK, rating)
}

// updateRating handles PATCH /v1/ratings/{id}
func updateRating(w http.ResponseWriter, r *http.Request) {
	ratingIDStr := chi.URLParam(r, "id")
	ratingID, err := strconv.Atoi(ratingIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ARGUMENT", fmt.Sprintf("Invalid rating id '%s'.", ratingIDStr))
		return
	}

	rating, ok := ratings[ratingID]
	if !ok {
		writeError(w, http.StatusNotFound, "NOT_FOUND", fmt.Sprintf("Rating with name 'ratings/%d' not found.", ratingID))
		return
	}

	var patch struct {
		Score   *int    `json:"score,omitempty"`
		Comment *string `json:"comment,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&patch); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ARGUMENT", "Invalid JSON body.")
		return
	}

	if patch.Score != nil {
		if *patch.Score < 1 || *patch.Score > 5 {
			writeError(w, http.StatusBadRequest, "INVALID_ARGUMENT", "Score must be between 1 and 5.")
			return
		}
		rating.Score = *patch.Score
	}

	if patch.Comment != nil {
		rating.Comment = strings.TrimSpace(*patch.Comment)
	}

	ratings[ratingID] = rating
	writeJSON(w, http.StatusOK, rating)
}

// deleteRating handles DELETE /v1/ratings/{id}
func deleteRating(w http.ResponseWriter, r *http.Request) {
	ratingIDStr := chi.URLParam(r, "id")
	ratingID, err := strconv.Atoi(ratingIDStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_ARGUMENT", fmt.Sprintf("Invalid rating id '%s'.", ratingIDStr))
		return
	}

	if _, ok := ratings[ratingID]; !ok {
		writeError(w, http.StatusNotFound, "NOT_FOUND", fmt.Sprintf("Rating with name 'ratings/%d' not found.", ratingID))
		return
	}

	delete(ratings, ratingID)
	w.WriteHeader(http.StatusNoContent)
}
