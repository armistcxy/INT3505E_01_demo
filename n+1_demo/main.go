package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cast"
)

var pool *pgxpool.Pool

func main() {
	dsn := "postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable"

	var err error
	pool, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		panic(err)
	}

	defer pool.Close()

	loadSchema()
	cleanUp()
	bootstrapData()

	http.HandleFunc("GET /api/v1/users/posts", ListUsersWithPostsV1)
	http.HandleFunc("GET /api/v2/users/posts", ListUsersWithPostsV2)

	log.Println("Ready to serve on http://localhost:8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Post struct {
	ID      int    `json:"id"`
	UserID  int    `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Demo n + 1 queries
func ListUsersWithPostsV1(w http.ResponseWriter, r *http.Request) {
	results, err := NPlus1Query()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Demo solving n + 1 queries with eager loading
func ListUsersWithPostsV2(w http.ResponseWriter, r *http.Request) {
	results := SolveNPlus1()
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type UserWithPosts struct {
	User  User   `json:"user"`
	Posts []Post `json:"posts"`
}

// n+1 queries
func NPlus1Query() ([]UserWithPosts, error) {
	results := make([]UserWithPosts, 0)
	rows, err := pool.Query(context.Background(), "SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		rows.Scan(&user.ID, &user.Name)

		posts := make([]Post, 0)
		rows2, err := pool.Query(context.Background(), "SELECT * FROM posts WHERE user_id = $1", user.ID)
		if err != nil {
			return nil, err
		}

		for rows2.Next() {
			var post Post
			err := rows2.Scan(&post.ID, &post.UserID, &post.Title, &post.Content)
			if err != nil {
				return nil, err
			}
			posts = append(posts, post)
		}

		rows2.Close()
		results = append(results, UserWithPosts{User: user, Posts: posts})
	}

	return results, nil
}

// solve n+1 query by using eager loading
func SolveNPlus1() []UserWithPosts {
	results := []UserWithPosts{}
	rows, err := pool.Query(context.Background(), "SELECT * FROM users")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var user User
		rows.Scan(&user.ID, &user.Name)
		users = append(users, user)
	}

	userIDs := []int{}
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
	}
	rows2, err := pool.Query(context.Background(), "SELECT * FROM posts WHERE user_id = ANY($1)", userIDs)
	if err != nil {
		panic(err)
	}
	defer rows2.Close()

	posts := []Post{}
	for rows2.Next() {
		var post Post
		rows2.Scan(&post.ID, &post.UserID, &post.Title, &post.Content)
		posts = append(posts, post)
	}

	for _, user := range users {
		userPosts := []Post{}
		for _, post := range posts {
			if post.UserID == user.ID {
				userPosts = append(userPosts, post)
			}
		}
		results = append(results, UserWithPosts{User: user, Posts: userPosts})
	}

	return results
}

func loadSchema() {
	// Load schema if not exists
	schema := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS posts (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES users(id),
			title TEXT NOT NULL,
			content TEXT NOT NULL
		);
	`
	_, _ = pool.Exec(context.Background(), schema)
}

func bootstrapData() {

	// create 100 users
	users := []User{}
	for i := 0; i < 100; i++ {
		users = append(users, User{ID: i + 1, Name: "User " + cast.ToString(i+1)})
	}

	batch := &pgx.Batch{}
	for _, user := range users {
		batch.Queue("INSERT INTO users (id, name) VALUES ($1, $2)", user.ID, user.Name)
	}
	pool.SendBatch(context.Background(), batch)

	// create 1000 posts
	posts := []Post{}
	for i := 0; i < 1000; i++ {
		posts = append(posts, Post{ID: i + 1, UserID: (i % 100) + 1, Title: "Post " + cast.ToString(i+1), Content: "Content " + cast.ToString(i+1)})
	}

	batch = &pgx.Batch{}
	for _, post := range posts {
		batch.Queue("INSERT INTO posts (id, user_id, title, content) VALUES ($1, $2, $3, $4)", post.ID, post.UserID, post.Title, post.Content)
	}
	pool.SendBatch(context.Background(), batch)
}

func cleanUp() {
	_, _ = pool.Exec(context.Background(), "DELETE FROM users")
	_, _ = pool.Exec(context.Background(), "DELETE FROM posts")
}
