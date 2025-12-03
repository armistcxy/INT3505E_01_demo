package service

import (
	"book-service/internal/models"
	"book-service/internal/repository"
)

type BookService struct {
	repo *repository.BookRepository
}

func NewBookService(repo *repository.BookRepository) *BookService {
	return &BookService{repo: repo}
}

func (s *BookService) CreateBook(book *models.CreateBookRequest) (*models.Book, error) {
	return s.repo.CreateBook(book)
}

func (s *BookService) GetBookByID(id int) (*models.Book, error) {
	return s.repo.GetBookByID(id)
}

func (s *BookService) GetAllBooks() ([]models.Book, error) {
	return s.repo.GetAllBooks()
}

func (s *BookService) UpdateBook(id int, book *models.UpdateBookRequest) (*models.Book, error) {
	return s.repo.UpdateBook(id, book)
}

func (s *BookService) DeleteBook(id int) error {
	return s.repo.DeleteBook(id)
}
