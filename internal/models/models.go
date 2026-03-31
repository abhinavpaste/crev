package models

import "time"

type User struct {
	ID        string    `db:"id"         json:"id"`
	Username  string    `db:"username"   json:"username"`
	Email     string    `db:"email"      json:"email"`
	Password  string    `db:"password"   json:"-"` // never serialized
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type Snippet struct {
	ID         string    `db:"id"          json:"id"`
	OwnerID    string    `db:"owner_id"    json:"owner_id"`
	Title      string    `db:"title"       json:"title"`
	Code       string    `db:"code"        json:"code"`
	Language   string    `db:"language"    json:"language"`
	ShareToken string    `db:"share_token" json:"share_token"`
	CreatedAt  time.Time `db:"created_at"  json:"created_at"`
}

type Comment struct {
	ID         string    `db:"id"          json:"id"`
	SnippetID  string    `db:"snippet_id"  json:"snippet_id"`
	AuthorID   string    `db:"author_id"   json:"author_id"`
	LineNumber int       `db:"line_number" json:"line_number"`
	Body       string    `db:"body"        json:"body"`
	CreatedAt  time.Time `db:"created_at"  json:"created_at"`
}
