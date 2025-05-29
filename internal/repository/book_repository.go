package repository

import (
	"context"
	"database/sql"

	"github.com/zizouhuweidi/bayt-alhikmah/internal/models"
	dbgen "github.com/zizouhuweidi/bayt-alhikmah/internal/repository/db"
)

type BookRepository interface {
	ListBooks(ctx context.Context, titleFilter string) ([]models.Book, error)
	GetBook(ctx context.Context, id int) (*models.Book, error)
	CreateBook(ctx context.Context, book models.Book) (*models.Book, error)
	UpdateBook(ctx context.Context, id int, book models.Book) (*models.Book, error)
	DeleteBook(ctx context.Context, id int) error
}

type bookRepo struct {
	q *dbgen.Queries
}

func NewBookRepository(db *sql.DB) BookRepository {
	return &bookRepo{
		q: dbgen.New(db),
	}
}

func (r *bookRepo) ListBooks(ctx context.Context, titleFilter string) ([]models.Book, error) {
	books, err := r.q.ListBooks(ctx, titleFilter)
	if err != nil {
		return nil, err
	}

	result := make([]models.Book, len(books))
	for i, book := range books {
		result[i] = models.Book{
			ID:          int(book.ID),
			Title:       book.Title,
			Author:      book.Author,
			Description: book.Description.String,
			Thumbnail:   book.ThumbnailUrl.String,
			CreatedAt:   book.CreatedAt,
			UpdatedAt:   book.UpdatedAt,
		}
	}
	return result, nil
}

func (r *bookRepo) GetBook(ctx context.Context, id int) (*models.Book, error) {
	book, err := r.q.GetBook(ctx, int32(id))
	if err != nil {
		return nil, err
	}

	return &models.Book{
		ID:          int(book.ID),
		Title:       book.Title,
		Author:      book.Author,
		Description: book.Description.String,
		Thumbnail:   book.ThumbnailUrl.String,
		CreatedAt:   book.CreatedAt,
		UpdatedAt:   book.UpdatedAt,
	}, nil
}

func (r *bookRepo) CreateBook(ctx context.Context, book models.Book) (*models.Book, error) {
	params := dbgen.CreateBookParams{
		Title:        book.Title,
		Author:       book.Author,
		Description:  sql.NullString{String: book.Description, Valid: book.Description != ""},
		ThumbnailUrl: sql.NullString{String: book.Thumbnail, Valid: book.Thumbnail != ""},
	}

	result, err := r.q.CreateBook(ctx, params)
	if err != nil {
		return nil, err
	}

	return &models.Book{
		ID:          int(result.ID),
		Title:       result.Title,
		Author:      result.Author,
		Description: result.Description.String,
		Thumbnail:   result.ThumbnailUrl.String,
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
	}, nil
}

func (r *bookRepo) UpdateBook(ctx context.Context, id int, book models.Book) (*models.Book, error) {
	params := dbgen.UpdateBookParams{
		ID:           int32(id),
		Title:        book.Title,
		Author:       book.Author,
		Description:  sql.NullString{String: book.Description, Valid: book.Description != ""},
		ThumbnailUrl: sql.NullString{String: book.Thumbnail, Valid: book.Thumbnail != ""},
	}

	result, err := r.q.UpdateBook(ctx, params)
	if err != nil {
		return nil, err
	}

	return &models.Book{
		ID:          int(result.ID),
		Title:       result.Title,
		Author:      result.Author,
		Description: result.Description.String,
		Thumbnail:   result.ThumbnailUrl.String,
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
	}, nil
}

func (r *bookRepo) DeleteBook(ctx context.Context, id int) error {
	return r.q.DeleteBook(ctx, int32(id))
}
