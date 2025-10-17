package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type Loan struct {
	BookID       int        `json:"book_id"`
	BookTitle    string     `json:"book_title"`
	CheckoutDate time.Time  `json:"checkout_date"`
	LoanID       int        `json:"loan_id"`
	ReturnDate   *time.Time `json:"return_date"` // có thể null
	UserID       int        `json:"user_id"`
	UserName     string     `json:"user_name"`
}

func getLoan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	loan := Loan{
		BookID:       1,
		BookTitle:    "Learn Python",
		CheckoutDate: time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC),
		LoanID:       1,
		ReturnDate:   nil,
		UserID:       2,
		UserName:     "Bao Nguyen",
	}

	json.NewEncoder(w).Encode(loan)
}

func main() {
	http.HandleFunc("/loan", getLoan)
	http.ListenAndServe(":8080", nil)
}
