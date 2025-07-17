package models

import "time"

type Book struct {
	ID          int       `db:"id"`
	Title       string    `db:"title"`
	Author      string    `db:"author"`
	Description string    `db:"description"`
	Thumbnail   string    `db:"thumbnail_url"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
