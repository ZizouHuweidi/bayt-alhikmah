// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package gen

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

const addUser = `-- name: AddUser :one
INSERT INTO users (username, email, password_hash, first_name, last_name)
VALUES ($1, $2, $3, $4, $5)
RETURNING user_id, username, email, password_hash, first_name, last_name, date_joined, last_login, is_active
`

type AddUserParams struct {
	Username     string
	Email        string
	PasswordHash string
	FirstName    pgtype.Text
	LastName     pgtype.Text
}

// Add a new user
func (q *Queries) AddUser(ctx context.Context, arg AddUserParams) (User, error) {
	row := q.db.QueryRow(ctx, addUser,
		arg.Username,
		arg.Email,
		arg.PasswordHash,
		arg.FirstName,
		arg.LastName,
	)
	var i User
	err := row.Scan(
		&i.UserID,
		&i.Username,
		&i.Email,
		&i.PasswordHash,
		&i.FirstName,
		&i.LastName,
		&i.DateJoined,
		&i.LastLogin,
		&i.IsActive,
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

const findBooksWithNoRatings = `-- name: FindBooksWithNoRatings :many
SELECT b.title, b.isbn
FROM books b
LEFT JOIN user_books ub ON b.book_id = ub.book_id
WHERE ub.book_id IS NULL
`

type FindBooksWithNoRatingsRow struct {
	Title string
	Isbn  pgtype.Text
}

// Find books that have no ratings
func (q *Queries) FindBooksWithNoRatings(ctx context.Context) ([]FindBooksWithNoRatingsRow, error) {
	rows, err := q.db.Query(ctx, findBooksWithNoRatings)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindBooksWithNoRatingsRow
	for rows.Next() {
		var i FindBooksWithNoRatingsRow
		if err := rows.Scan(&i.Title, &i.Isbn); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findInactiveUsers = `-- name: FindInactiveUsers :many
SELECT user_id, username, email, last_login
FROM users
WHERE last_login < NOW() - INTERVAL '30 days'
   OR last_login IS NULL
`

type FindInactiveUsersRow struct {
	UserID    int32
	Username  string
	Email     string
	LastLogin pgtype.Timestamp
}

// Find users who have not logged in for the last 30 days
func (q *Queries) FindInactiveUsers(ctx context.Context) ([]FindInactiveUsersRow, error) {
	rows, err := q.db.Query(ctx, findInactiveUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindInactiveUsersRow
	for rows.Next() {
		var i FindInactiveUsersRow
		if err := rows.Scan(
			&i.UserID,
			&i.Username,
			&i.Email,
			&i.LastLogin,
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

const getActiveUsers = `-- name: GetActiveUsers :many
SELECT user_id, username, email, first_name, last_name
FROM users
WHERE is_active = TRUE
`

type GetActiveUsersRow struct {
	UserID    int32
	Username  string
	Email     string
	FirstName pgtype.Text
	LastName  pgtype.Text
}

// Retrieve all active users
func (q *Queries) GetActiveUsers(ctx context.Context) ([]GetActiveUsersRow, error) {
	rows, err := q.db.Query(ctx, getActiveUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetActiveUsersRow
	for rows.Next() {
		var i GetActiveUsersRow
		if err := rows.Scan(
			&i.UserID,
			&i.Username,
			&i.Email,
			&i.FirstName,
			&i.LastName,
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

const getAverageRatingByGenre = `-- name: GetAverageRatingByGenre :many
SELECT b.genre, AVG(ub.rating) as average_rating
FROM books b
JOIN user_books ub ON b.book_id = ub.book_id
GROUP BY b.genre
ORDER BY average_rating DESC
`

type GetAverageRatingByGenreRow struct {
	Genre         pgtype.Text
	AverageRating float64
}

// Get the average rating for each genre
func (q *Queries) GetAverageRatingByGenre(ctx context.Context) ([]GetAverageRatingByGenreRow, error) {
	rows, err := q.db.Query(ctx, getAverageRatingByGenre)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAverageRatingByGenreRow
	for rows.Next() {
		var i GetAverageRatingByGenreRow
		if err := rows.Scan(&i.Genre, &i.AverageRating); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getBookReviewers = `-- name: GetBookReviewers :many
SELECT u.username, ub.rating, ub.review
FROM users u
JOIN user_books ub ON u.user_id = ub.user_id
WHERE ub.book_id = $1 AND ub.review IS NOT NULL
`

type GetBookReviewersRow struct {
	Username string
	Rating   pgtype.Int4
	Review   pgtype.Text
}

// Get a list of users who have reviewed a specific book
func (q *Queries) GetBookReviewers(ctx context.Context, bookID int32) ([]GetBookReviewersRow, error) {
	rows, err := q.db.Query(ctx, getBookReviewers, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetBookReviewersRow
	for rows.Next() {
		var i GetBookReviewersRow
		if err := rows.Scan(&i.Username, &i.Rating, &i.Review); err != nil {
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

const getMostActiveUsers = `-- name: GetMostActiveUsers :many
SELECT u.username, COUNT(ub.book_id) as books_read
FROM users u
JOIN user_books ub ON u.user_id = ub.user_id
WHERE ub.status = 'Completed'
GROUP BY u.user_id, u.username
ORDER BY books_read DESC
LIMIT 5
`

type GetMostActiveUsersRow struct {
	Username  string
	BooksRead int64
}

// Find the most active users (users who have read the most books)
func (q *Queries) GetMostActiveUsers(ctx context.Context) ([]GetMostActiveUsersRow, error) {
	rows, err := q.db.Query(ctx, getMostActiveUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetMostActiveUsersRow
	for rows.Next() {
		var i GetMostActiveUsersRow
		if err := rows.Scan(&i.Username, &i.BooksRead); err != nil {
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

const getTopAuthors = `-- name: GetTopAuthors :many
SELECT a.first_name, a.last_name, COUNT(ba.book_id) as book_count
FROM authors a
JOIN book_authors ba ON a.author_id = ba.author_id
GROUP BY a.author_id, a.first_name, a.last_name
ORDER BY book_count DESC
LIMIT 5
`

type GetTopAuthorsRow struct {
	FirstName string
	LastName  string
	BookCount int64
}

// Find authors with the most books in the library
func (q *Queries) GetTopAuthors(ctx context.Context) ([]GetTopAuthorsRow, error) {
	rows, err := q.db.Query(ctx, getTopAuthors)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTopAuthorsRow
	for rows.Next() {
		var i GetTopAuthorsRow
		if err := rows.Scan(&i.FirstName, &i.LastName, &i.BookCount); err != nil {
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

const getUserReadingList = `-- name: GetUserReadingList :many
SELECT b.title, b.isbn, ub.start_date
FROM books b
JOIN user_books ub ON b.book_id = ub.book_id
WHERE ub.user_id = $1 AND ub.status = 'Reading'
`

type GetUserReadingListRow struct {
	Title     string
	Isbn      pgtype.Text
	StartDate pgtype.Date
}

// Get a list of all books in a user's 'Reading' list
func (q *Queries) GetUserReadingList(ctx context.Context, userID int32) ([]GetUserReadingListRow, error) {
	rows, err := q.db.Query(ctx, getUserReadingList, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUserReadingListRow
	for rows.Next() {
		var i GetUserReadingListRow
		if err := rows.Scan(&i.Title, &i.Isbn, &i.StartDate); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUserLastLogin = `-- name: UpdateUserLastLogin :exec
UPDATE users
SET last_login = NOW()
WHERE user_id = $1
`

// Update user's last login time
func (q *Queries) UpdateUserLastLogin(ctx context.Context, userID int32) error {
	_, err := q.db.Exec(ctx, updateUserLastLogin, userID)
	return err
}