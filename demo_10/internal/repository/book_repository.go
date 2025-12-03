package repository

import (
	"database/sql"
	"errors"

	"book-service/internal/models"
)

type BookRepository struct {
	db *sql.DB
}

func NewBookRepository(db *sql.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) CreateBook(book *models.CreateBookRequest) (*models.Book, error) {
	query := `
		INSERT INTO books (title, author, isbn, pages, published, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, title, author, isbn, pages, published, created_at, updated_at
	`

	var createdBook models.Book
	err := r.db.QueryRow(
		query,
		book.Title,
		book.Author,
		book.ISBN,
		book.Pages,
		book.Published,
	).Scan(
		&createdBook.ID,
		&createdBook.Title,
		&createdBook.Author,
		&createdBook.ISBN,
		&createdBook.Pages,
		&createdBook.Published,
		&createdBook.CreatedAt,
		&createdBook.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &createdBook, nil
}

func (r *BookRepository) GetBookByID(id int) (*models.Book, error) {
	query := `
		SELECT id, title, author, isbn, pages, published, created_at, updated_at
		FROM books WHERE id = $1
	`

	var book models.Book
	err := r.db.QueryRow(query, id).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.ISBN,
		&book.Pages,
		&book.Published,
		&book.CreatedAt,
		&book.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("book not found")
		}
		return nil, err
	}

	return &book, nil
}

func (r *BookRepository) GetAllBooks() ([]models.Book, error) {
	query := `
		SELECT id, title, author, isbn, pages, published, created_at, updated_at
		FROM books ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Author,
			&book.ISBN,
			&book.Pages,
			&book.Published,
			&book.CreatedAt,
			&book.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	return books, rows.Err()
}

func (r *BookRepository) UpdateBook(id int, book *models.UpdateBookRequest) (*models.Book, error) {
	// Get current book first
	current, err := r.GetBookByID(id)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if book.Title != nil {
		current.Title = *book.Title
	}
	if book.Author != nil {
		current.Author = *book.Author
	}
	if book.ISBN != nil {
		current.ISBN = *book.ISBN
	}
	if book.Pages != nil {
		current.Pages = *book.Pages
	}
	if book.Published != nil {
		current.Published = *book.Published
	}

	query := `
		UPDATE books
		SET title = $1, author = $2, isbn = $3, pages = $4, published = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
		RETURNING id, title, author, isbn, pages, published, created_at, updated_at
	`

	var updatedBook models.Book
	err = r.db.QueryRow(
		query,
		current.Title,
		current.Author,
		current.ISBN,
		current.Pages,
		current.Published,
		id,
	).Scan(
		&updatedBook.ID,
		&updatedBook.Title,
		&updatedBook.Author,
		&updatedBook.ISBN,
		&updatedBook.Pages,
		&updatedBook.Published,
		&updatedBook.CreatedAt,
		&updatedBook.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &updatedBook, nil
}

func (r *BookRepository) DeleteBook(id int) error {
	query := `DELETE FROM books WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("book not found")
	}

	return nil
}
