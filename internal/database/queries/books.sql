-- name: GetBooks :many
-- Get all books
SELECT * FROM books;

-- name: FindBooksByAuthor :many
-- Find all books by a specific author
SELECT b.title, b.isbn, b.publication_date
FROM books b
JOIN book_authors ba ON b.book_id = ba.book_id
JOIN authors a ON ba.author_id = a.author_id
WHERE a.last_name = $1 AND a.first_name = $2;

-- name: GetTopRatedBooks :many
-- Get the top 10 highest-rated books
SELECT b.title, b.isbn, AVG(ub.rating) as average_rating
FROM books b
JOIN user_books ub ON b.book_id = ub.book_id
GROUP BY b.book_id, b.title, b.isbn
ORDER BY average_rating DESC
LIMIT 10;

-- name: GetRecentBooks :many
-- Get a list of books published in the last year
SELECT title, isbn, publication_date
FROM books
WHERE publication_date >= NOW() - INTERVAL '1 year';

-- name: GetBooksWithMultipleAuthors :many
-- Find books with multiple authors
SELECT b.title, STRING_AGG(CONCAT(a.first_name, ' ', a.last_name), ', ') as authors
FROM books b
JOIN book_authors ba ON b.book_id = ba.book_id
JOIN authors a ON ba.author_id = a.author_id
GROUP BY b.book_id, b.title
HAVING COUNT(DISTINCT a.author_id) > 1;


-- name: AddBook :one
-- Add a new book
INSERT INTO books (title, isbn, publication_date, genre, language, description, cover_image_url)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: AddBookRating :one
-- Add or update a book rating
INSERT INTO user_books (user_id, book_id, status, rating, review)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (user_id, book_id) 
DO UPDATE SET status = EXCLUDED.status, rating = EXCLUDED.rating, review = EXCLUDED.review
RETURNING *;
