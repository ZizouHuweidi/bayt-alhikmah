-- Users table
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    date_joined DATE NOT NULL DEFAULT CURRENT_DATE,
    last_login TIMESTAMP,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

-- Authors table
CREATE TABLE authors (
    author_id SERIAL PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    birth_date DATE,
    nationality VARCHAR(50),
    biography TEXT
);

-- Books table
CREATE TABLE books (
    book_id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    isbn VARCHAR(13) UNIQUE,
    publication_date DATE,
    genre VARCHAR(50),
    language VARCHAR(50),
    description TEXT,
    cover_image_url VARCHAR(255)
);

-- BookAuthors table (for many-to-many relationship between Books and Authors)
CREATE TABLE book_authors (
    book_id INTEGER NOT NULL REFERENCES books(book_id) ON DELETE CASCADE,
    author_id INTEGER NOT NULL REFERENCES authors(author_id) ON DELETE CASCADE,
    PRIMARY KEY (book_id, author_id)
);

-- UserBooks table (for tracking user's book interactions)
CREATE TABLE user_books (
    user_id INTEGER NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    book_id INTEGER NOT NULL REFERENCES books(book_id) ON DELETE CASCADE,
    status TEXT NOT NULL CHECK (status IN ('Reading', 'Completed', 'Wishlist')),
    start_date DATE,
    end_date DATE,
    rating INTEGER CHECK (rating >= 1 AND rating <= 5),
    review TEXT,
    PRIMARY KEY (user_id, book_id)
);

-- Create an index on the Books table for faster searches
CREATE INDEX idx_book_title ON books(title);

-- Create an index on the Authors table for faster searches
CREATE INDEX idx_author_name ON authors(last_name, first_name);
