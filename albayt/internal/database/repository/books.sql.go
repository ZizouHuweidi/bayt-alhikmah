// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: books.sql

package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const addBook = `-- name: AddBook :one
INSERT INTO books (title, isbn, publication_date, genre, language, description, cover_image_url)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING book_id, title, isbn, publication_date, genre, language, description, cover_image_url
`

type AddBookParams struct {
	Title           string
	Isbn            pgtype.Text
	PublicationDate pgtype.Date
	Genre           pgtype.Text
	Language        pgtype.Text
	Description     pgtype.Text
	CoverImageUrl   pgtype.Text
}

// Add a new book
func (q *Queries) AddBook(ctx context.Context, arg AddBookParams) (Book, error) {
	row := q.db.QueryRow(ctx, addBook,
		arg.Title,
		arg.Isbn,
		arg.PublicationDate,
		arg.Genre,
		arg.Language,
		arg.Description,
		arg.CoverImageUrl,
	)
	var i Book
	err := row.Scan(
		&i.BookID,
		&i.Title,
		&i.Isbn,
		&i.PublicationDate,
		&i.Genre,
		&i.Language,
		&i.Description,
		&i.CoverImageUrl,
	)
	return i, err
}

const addBookRating = `-- name: AddBookRating :one
INSERT INTO user_books (user_id, book_id, status, rating, review)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (user_id, book_id) 
DO UPDATE SET status = EXCLUDED.status, rating = EXCLUDED.rating, review = EXCLUDED.review
RETURNING user_id, book_id, status, start_date, end_date, rating, review
`

type AddBookRatingParams struct {
	UserID int32
	BookID int32
	Status string
	Rating pgtype.Int4
	Review pgtype.Text
}

// Add or update a book rating
func (q *Queries) AddBookRating(ctx context.Context, arg AddBookRatingParams) (UserBook, error) {
	row := q.db.QueryRow(ctx, addBookRating,
		arg.UserID,
		arg.BookID,
		arg.Status,
		arg.Rating,
		arg.Review,
	)
	var i UserBook
	err := row.Scan(
		&i.UserID,
		&i.BookID,
		&i.Status,
		&i.StartDate,
		&i.EndDate,
		&i.Rating,
		&i.Review,
	)
	return i, err
}

const findBooksByAuthor = `-- name: FindBooksByAuthor :many
SELECT b.title, b.isbn, b.publication_date
FROM books b
JOIN book_authors ba ON b.book_id = ba.book_id
JOIN authors a ON ba.author_id = a.author_id
WHERE a.last_name = $1 AND a.first_name = $2
`

type FindBooksByAuthorParams struct {
	LastName  string
	FirstName string
}

type FindBooksByAuthorRow struct {
	Title           string
	Isbn            pgtype.Text
	PublicationDate pgtype.Date
}

// Find all books by a specific author
func (q *Queries) FindBooksByAuthor(ctx context.Context, arg FindBooksByAuthorParams) ([]FindBooksByAuthorRow, error) {
	rows, err := q.db.Query(ctx, findBooksByAuthor, arg.LastName, arg.FirstName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindBooksByAuthorRow
	for rows.Next() {
		var i FindBooksByAuthorRow
		if err := rows.Scan(&i.Title, &i.Isbn, &i.PublicationDate); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getBooks = `-- name: GetBooks :many
SELECT book_id, title, isbn, publication_date, genre, language, description, cover_image_url FROM books
`

// Get all books
func (q *Queries) GetBooks(ctx context.Context) ([]Book, error) {
	rows, err := q.db.Query(ctx, getBooks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Book
	for rows.Next() {
		var i Book
		if err := rows.Scan(
			&i.BookID,
			&i.Title,
			&i.Isbn,
			&i.PublicationDate,
			&i.Genre,
			&i.Language,
			&i.Description,
			&i.CoverImageUrl,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getBooksWithMultipleAuthors = `-- name: GetBooksWithMultipleAuthors :many
SELECT b.title, STRING_AGG(CONCAT(a.first_name, ' ', a.last_name), ', ') as authors
FROM books b
JOIN book_authors ba ON b.book_id = ba.book_id
JOIN authors a ON ba.author_id = a.author_id
GROUP BY b.book_id, b.title
HAVING COUNT(DISTINCT a.author_id) > 1
`

type GetBooksWithMultipleAuthorsRow struct {
	Title   string
	Authors []byte
}

// Find books with multiple authors
func (q *Queries) GetBooksWithMultipleAuthors(ctx context.Context) ([]GetBooksWithMultipleAuthorsRow, error) {
	rows, err := q.db.Query(ctx, getBooksWithMultipleAuthors)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetBooksWithMultipleAuthorsRow
	for rows.Next() {
		var i GetBooksWithMultipleAuthorsRow
		if err := rows.Scan(&i.Title, &i.Authors); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRecentBooks = `-- name: GetRecentBooks :many
SELECT title, isbn, publication_date
FROM books
WHERE publication_date >= NOW() - INTERVAL '1 year'
`

type GetRecentBooksRow struct {
	Title           string
	Isbn            pgtype.Text
	PublicationDate pgtype.Date
}

// Get a list of books published in the last year
func (q *Queries) GetRecentBooks(ctx context.Context) ([]GetRecentBooksRow, error) {
	rows, err := q.db.Query(ctx, getRecentBooks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetRecentBooksRow
	for rows.Next() {
		var i GetRecentBooksRow
		if err := rows.Scan(&i.Title, &i.Isbn, &i.PublicationDate); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTopRatedBooks = `-- name: GetTopRatedBooks :many
SELECT b.title, b.isbn, AVG(ub.rating) as average_rating
FROM books b
JOIN user_books ub ON b.book_id = ub.book_id
GROUP BY b.book_id, b.title, b.isbn
ORDER BY average_rating DESC
LIMIT 10
`

type GetTopRatedBooksRow struct {
	Title         string
	Isbn          pgtype.Text
	AverageRating float64
}

// Get the top 10 highest-rated books
func (q *Queries) GetTopRatedBooks(ctx context.Context) ([]GetTopRatedBooksRow, error) {
	rows, err := q.db.Query(ctx, getTopRatedBooks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTopRatedBooksRow
	for rows.Next() {
		var i GetTopRatedBooksRow
		if err := rows.Scan(&i.Title, &i.Isbn, &i.AverageRating); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
