
-- name: GetActiveUsers :many
-- Retrieve all active users
SELECT user_id, username, email, first_name, last_name
FROM users
WHERE is_active = TRUE;


-- name: FindInactiveUsers :many
-- Find users who have not logged in for the last 30 days
SELECT user_id, username, email, last_login
FROM users
WHERE last_login < NOW() - INTERVAL '30 days'
   OR last_login IS NULL;

-- name: GetUserReadingList :many
-- Get a list of all books in a user's 'Reading' list
SELECT b.title, b.isbn, ub.start_date
FROM books b
JOIN user_books ub ON b.book_id = ub.book_id
WHERE ub.user_id = $1 AND ub.status = 'Reading';

-- name: GetTopAuthors :many
-- Find authors with the most books in the library
SELECT a.first_name, a.last_name, COUNT(ba.book_id) as book_count
FROM authors a
JOIN book_authors ba ON a.author_id = ba.author_id
GROUP BY a.author_id, a.first_name, a.last_name
ORDER BY book_count DESC
LIMIT 5;

-- name: GetAverageRatingByGenre :many
-- Get the average rating for each genre
SELECT b.genre, AVG(ub.rating) as average_rating
FROM books b
JOIN user_books ub ON b.book_id = ub.book_id
GROUP BY b.genre
ORDER BY average_rating DESC;

-- name: FindBooksWithNoRatings :many
-- Find books that have no ratings
SELECT b.title, b.isbn
FROM books b
LEFT JOIN user_books ub ON b.book_id = ub.book_id
WHERE ub.book_id IS NULL;

-- name: GetBookReviewers :many
-- Get a list of users who have reviewed a specific book
SELECT u.username, ub.rating, ub.review
FROM users u
JOIN user_books ub ON u.user_id = ub.user_id
WHERE ub.book_id = $1 AND ub.review IS NOT NULL;

-- name: GetMostActiveUsers :many
-- Find the most active users (users who have read the most books)
SELECT u.username, COUNT(ub.book_id) as books_read
FROM users u
JOIN user_books ub ON u.user_id = ub.user_id
WHERE ub.status = 'Completed'
GROUP BY u.user_id, u.username
ORDER BY books_read DESC
LIMIT 5;


-- name: AddUser :one
-- Add a new user
INSERT INTO users (username, email, password_hash, first_name, last_name)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateUserLastLogin :exec
-- Update user's last login time
UPDATE users
SET last_login = NOW()
WHERE user_id = $1;

