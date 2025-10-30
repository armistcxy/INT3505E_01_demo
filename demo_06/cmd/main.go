package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func main() {
	http.HandleFunc("/api/login", loginHandler)

	http.HandleFunc("/api/dashboard", authMiddleware(dashboardHandler))

	http.HandleFunc("/api/admin/users", requireRole("admin", adminUsersHandler))
	http.HandleFunc("/api/admin/stats", requireRole("admin", adminStatsHandler))

	log.Println("Server starting on :8080")

	log.Println("POST   /api/login         - Login endpoint")

	log.Printf("\nProtected endpoints (requires authentication):")
	log.Println("GET    /api/dashboard     - Get user dashboard")

	log.Printf("\nAdmin endpoints (requires admin role):")
	log.Println("GET    /api/admin/users   - Get all users")
	log.Println("GET    /api/admin/stats   - Get system statistics")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// User contains user information
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	Role     string `json:"role"`
}

// Claims contains user information and JWT claims
type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// LoginRequest contains login information
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse contains JWT token and user information
type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// Mock database
var users = []User{
	{ID: 1, Username: "admin", Password: "$2a$10$...", Role: "admin"},
	{ID: 2, Username: "user", Password: "$2a$10$...", Role: "user"},
}

func init() {
	// Hash passwords for mock users
	adminPass, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	userPass, _ := bcrypt.GenerateFromPassword([]byte("user123"), bcrypt.DefaultCost)
	users[0].Password = string(adminPass)
	users[1].Password = string(userPass)
}

func generateJWT(user User) (string, error) {
	expirationTime := time.Now().Add(10 * time.Minute)

	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func verifyToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Bearer token format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		claims, err := verifyToken(parts[1])
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		r.Header.Set("X-User-ID", fmt.Sprintf("%d", claims.UserID))
		r.Header.Set("X-User-Role", claims.Role)

		next(w, r)
	}
}

func requireRole(role string, next http.HandlerFunc) http.HandlerFunc {
	return authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		userRole := r.Header.Get("X-User-Role")

		if userRole != role {
			http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
			return
		}

		next(w, r)
	})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Find user
	var foundUser *User
	for i := range users {
		if users[i].Username == req.Username {
			foundUser = &users[i]
			break
		}
	}

	if foundUser == nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate token
	token, err := generateJWT(*foundUser)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Response
	resp := LoginResponse{
		Token: token,
		User:  *foundUser,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")

	dashboard := map[string]interface{}{
		"total_orders":    15,
		"pending_orders":  3,
		"total_spent":     1250.50,
		"recent_activity": []string{"Order #123 placed", "Profile updated", "Password changed"},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": fmt.Sprintf("Dashboard data for user %s", userID),
		"data":    dashboard,
	})
}

func adminUsersHandler(w http.ResponseWriter, r *http.Request) {
	allUsers := []map[string]interface{}{
		{"id": 1, "username": "admin", "role": "admin", "status": "active", "created_at": "2024-01-01"},
		{"id": 2, "username": "user", "role": "user", "status": "active", "created_at": "2024-01-10"},
		{"id": 3, "username": "john_doe", "role": "user", "status": "inactive", "created_at": "2024-02-05"},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":     "Users list retrieved successfully",
		"total_users": len(allUsers),
		"data":        allUsers,
	})
}

func adminStatsHandler(w http.ResponseWriter, r *http.Request) {
	stats := map[string]interface{}{
		"total_users":        156,
		"active_users":       142,
		"total_orders":       1523,
		"revenue_this_month": 45678.90,
		"system_uptime":      "15 days 4 hours",
		"last_backup":        time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "System statistics retrieved successfully",
		"data":    stats,
	})
}
