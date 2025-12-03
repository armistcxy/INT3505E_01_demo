package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"

	"book-service/internal/models"
	"book-service/internal/service"
)

type BookHandler struct {
	service *service.BookService

	// Prometheus metrics
	requestsTotal    prometheus.Counter
	requestsDuration prometheus.Histogram
}

func NewBookHandler(service *service.BookService) *BookHandler {
	requestsTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "book_requests_total",
		Help: "Total number of book requests",
	})

	requestsDuration := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "book_request_duration_seconds",
		Help: "Duration of book requests in seconds",
	})

	prometheus.MustRegister(requestsTotal, requestsDuration)

	return &BookHandler{
		service:          service,
		requestsTotal:    requestsTotal,
		requestsDuration: requestsDuration,
	}
}

func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	defer h.recordMetrics()()

	var req models.CreateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	book, err := h.service.CreateBook(&req)
	if err != nil {
		http.Error(w, "Failed to create book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	defer h.recordMetrics()()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	book, err := h.service.GetBookByID(id)
	if err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	defer h.recordMetrics()()

	books, err := h.service.GetAllBooks()
	if err != nil {
		http.Error(w, "Failed to retrieve books", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	defer h.recordMetrics()()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	var req models.UpdateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	book, err := h.service.UpdateBook(id, &req)
	if err != nil {
		http.Error(w, "Failed to update book", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	defer h.recordMetrics()()

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid book ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteBook(id)
	if err != nil {
		http.Error(w, "Failed to delete book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *BookHandler) recordMetrics() func() {
	timer := prometheus.NewTimer(h.requestsDuration)
	return func() {
		h.requestsTotal.Inc()
		timer.ObserveDuration()
	}
}

func (h *BookHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/api/books", h.GetAllBooks).Methods("GET")
	router.HandleFunc("/api/books", h.CreateBook).Methods("POST")
	router.HandleFunc("/api/books/{id}", h.GetBook).Methods("GET")
	router.HandleFunc("/api/books/{id}", h.UpdateBook).Methods("PUT")
	router.HandleFunc("/api/books/{id}", h.DeleteBook).Methods("DELETE")
}
