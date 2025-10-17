package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

const validToken = "my-secret-token"

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if !strings.HasPrefix(authHeader, "Bearer ") || strings.TrimPrefix(authHeader, "Bearer ") != validToken {
			http.Error(w, "Unauthorized - Missing or invalid token", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{
		"message": "Access granted! Here is your profile.",
	}
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/profile", authMiddleware(profileHandler))

	http.ListenAndServe(":8080", nil)
}
