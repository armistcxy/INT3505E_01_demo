package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type StudentV1 struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func listStudentsVersion1(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	students := []StudentV1{
		{ID: 1, Name: "Alice", Age: 10},
		{ID: 2, Name: "Bob", Age: 11},
		{ID: 3, Name: "Charlie", Age: 12},
	}

	_ = json.NewEncoder(w).Encode(students)
}

type StudentV2 struct {
	StudentID int    `json:"student_id"`
	FullName  string `json:"full_name"`
	Age       int    `json:"age"`
	Class     string `json:"class"`
}

func listStudentsVersion2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	students := []StudentV2{
		{StudentID: 1, FullName: "Alice", Age: 10, Class: "A"},
		{StudentID: 2, FullName: "Bob", Age: 11, Class: "B"},
		{StudentID: 3, FullName: "Charlie", Age: 12, Class: "C"},
	}

	_ = json.NewEncoder(w).Encode(students)
}

func main() {
	http.HandleFunc("/v1/students", listStudentsVersion1)
	http.HandleFunc("/v2/students", listStudentsVersion2)
	log.Println("Server is listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
