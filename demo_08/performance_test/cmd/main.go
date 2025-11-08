package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	pft "demo_08/performance_test"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
)

func main() {
	rc := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	defer rc.Close()
	pool, err := pgxpool.New(context.Background(), "postgres://postgres:postgres@localhost:5432/postgres")
	if err != nil {
		panic(err)
	}
	bootstrapData(pool)

	bookQuerier := pft.NewBookQuerier(pool, rc)

	http.HandleFunc("/api/v1/books/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		book, err := bookQuerier.GetBookWithoutCache(r.Context(), cast.ToInt(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(book); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/api/v2/books/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		book, err := bookQuerier.GetBookWithCache(r.Context(), cast.ToInt(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(book); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	log.Println("Server is listening on port 8080")
	http.ListenAndServe(":8080", nil)
}

func bootstrapData(pool *pgxpool.Pool) {
	_, err := pool.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS books (id INT, name TEXT)")
	if err != nil {
		panic(err)
	}

	_, err = pool.Exec(context.Background(), "TRUNCATE TABLE books")
	if err != nil {
		panic(err)
	}

	values := [][]interface{}{}
	for i := 0; i < 10000; i++ {
		values = append(values, []interface{}{i + 1, "Book " + cast.ToString(i+1)})
	}

	_, err = pool.CopyFrom(context.Background(), []string{"books"}, []string{"id", "name"}, pgx.CopyFromRows(values))
	if err != nil {
		panic(err)
	}
}
