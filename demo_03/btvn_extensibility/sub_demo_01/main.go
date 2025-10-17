package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Student struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Grade int    `json:"grade,omitempty"`
	Email string `json:"email,omitempty"`
}

// Fake data
var students = []Student{
	{ID: 1, Name: "Alice", Age: 10, Grade: 10, Email: "alice@gmail.com"},
	{ID: 2, Name: "Bob", Age: 11, Grade: 11, Email: "bob@gmail.com"},
	{ID: 3, Name: "Charlie", Age: 12, Grade: 12, Email: "charlie@gmail.com"},
}

func listStudentsBasic(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	basicStudents := []Student{}
	for _, s := range students {
		basicStudents = append(basicStudents, Student{
			ID:   s.ID,
			Name: s.Name,
			Age:  s.Age,
		})
	}

	_ = json.NewEncoder(w).Encode(basicStudents)
}

func listStudentsUpgraded(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(students)
}

func main() {
	// http.HandleFunc("/students", listStudentsBasic)
	http.HandleFunc("/students", listStudentsUpgraded)
	log.Println("Server is listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
