package repository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/zizouhuweidi/bayt-alhikmah/internal/models"
)

type BookRepository interface {
	ListBooks(ctx context.Context, titleFilter string) ([]models.Book, error)
	// other methods: GetBook, CreateBook, UpdateBook, DeleteBook...
}

type bookRepo struct {
	db *sqlx.DB
}

func NewBookRepository(db *sqlx.DB) BookRepository {
	return &bookRepo{db: db}
}

func (r *bookRepo) ListBooks(ctx context.Context, titleFilter string) ([]models.Book, error) {
	builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select("id", "title", "author", "description", "thumbnail_url", "created_at", "updated_at").
		From("books")
	if titleFilter != "" {
		builder = builder.Where(sq.Like{"title": fmt.Sprintf("%%%s%%", titleFilter)})
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var books []models.Book
	err = r.db.SelectContext(ctx, &books, sql, args...)
	return books, err
}
