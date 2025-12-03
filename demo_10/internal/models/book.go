package models

import "time"

type Book struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	ISBN      string    `json:"isbn"`
	Pages     int       `json:"pages"`
	Published time.Time `json:"published"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateBookRequest struct {
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	ISBN      string    `json:"isbn"`
	Pages     int       `json:"pages"`
	Published time.Time `json:"published"`
}

type UpdateBookRequest struct {
	Title     *string    `json:"title,omitempty"`
	Author    *string    `json:"author,omitempty"`
	ISBN      *string    `json:"isbn,omitempty"`
	Pages     *int       `json:"pages,omitempty"`
	Published *time.Time `json:"published,omitempty"`
}
