package models

import "time"

type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type Snippet struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Language  string    `json:"language"`
	CreatedAt time.Time `json:"created_at"`
}

type Comment struct {
	ID         string    `json:"id"`
	SnippetID  string    `json:"snippet_id"`
	UserID     string    `json:"user_id"`
	LineNumber int       `json:"line_number"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}
