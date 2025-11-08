package performancetest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Book struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type BookQuerier struct {
	db *pgxpool.Pool
	rc *redis.Client
}

func NewBookQuerier(db *pgxpool.Pool, rc *redis.Client) *BookQuerier {
	return &BookQuerier{db: db, rc: rc}
}

func (b *BookQuerier) GetBookWithoutCache(ctx context.Context, id int) (*Book, error) {
	row := b.db.QueryRow(ctx, "SELECT id, name FROM books WHERE id = $1", id)
	book := &Book{}
	if err := row.Scan(&book.ID, &book.Name); err != nil {
		log.Println("Error scanning row:", err)
		return nil, err
	}

	return book, nil
}

func (b *BookQuerier) GetBookWithCache(ctx context.Context, id int) (*Book, error) {
	rmd := b.rc.Get(ctx, fmt.Sprintf("book:%d", id))
	if rmd.Err() == nil {
		book := &Book{}
		if err := json.Unmarshal([]byte(rmd.Val()), book); err != nil {
			return nil, err
		}
		return book, nil
	}

	row := b.db.QueryRow(ctx, "SELECT id, name FROM books WHERE id = $1", id)
	book := &Book{}
	if err := row.Scan(&book.ID, &book.Name); err != nil {
		return nil, err
	}

	b.rc.Set(ctx, fmt.Sprintf("book:%d", id), book, 0)

	return book, nil
}
