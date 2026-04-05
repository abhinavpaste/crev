package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/abhinavpaste/crev/internal/db"
	"github.com/abhinavpaste/crev/internal/middleware"
	"github.com/abhinavpaste/crev/internal/models"
)

type commentRequest struct {
	LineNumber int    `json:"line_number"`
	Content    string `json:"content"`
}

func CreateComment(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	snippetID := strings.TrimPrefix(r.URL.Path, "/snippets/")
	snippetID = strings.TrimSuffix(snippetID, "/comments")

	var req commentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	var comment models.Comment
	err := db.DB.QueryRow(
		`INSERT INTO comments (snippet_id, user_id, line_number, content) VALUES ($1, $2, $3, $4) RETURNING id, snippet_id, user_id, line_number, content, created_at`,
		snippetID, userID, req.LineNumber, req.Content,
	).Scan(&comment.ID, &comment.SnippetID, &comment.UserID, &comment.LineNumber, &comment.Content, &comment.CreatedAt)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

func GetComments(w http.ResponseWriter, r *http.Request) {
	snippetID := strings.TrimPrefix(r.URL.Path, "/snippets/")
	snippetID = strings.TrimSuffix(snippetID, "/comments")

	rows, err := db.DB.Query(
		`SELECT id, snippet_id, user_id, line_number, content, created_at FROM comments WHERE snippet_id = $1 ORDER BY line_number`,
		snippetID,
	)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	comments := []models.Comment{}
	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(&c.ID, &c.SnippetID, &c.UserID, &c.LineNumber, &c.Content, &c.CreatedAt); err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		comments = append(comments, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}
