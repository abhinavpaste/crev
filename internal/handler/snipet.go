package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/abhinavpaste/crev/internal/db"
	"github.com/abhinavpaste/crev/internal/middleware"
	"github.com/abhinavpaste/crev/internal/models"
)

type snippetRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Language string `json:"language"`
}

func CreateSnippet(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(string)

	var req snippetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	var snippet models.Snippet
	err := db.DB.QueryRow(
		`INSERT INTO snippets (user_id, title, content, language) VALUES ($1, $2, $3, $4) RETURNING id, user_id, title, content, language, created_at`,
		userID, req.Title, req.Content, req.Language,
	).Scan(&snippet.ID, &snippet.UserID, &snippet.Title, &snippet.Content, &snippet.Language, &snippet.CreatedAt)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(snippet)
}

func GetSnippet(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/snippets/")
	if id == "" {
		http.Error(w, "missing snippet id", http.StatusBadRequest)
		return
	}

	var snippet models.Snippet
	err := db.DB.QueryRow(
		`SELECT id, user_id, title, content, language, created_at FROM snippets WHERE id = $1`,
		id,
	).Scan(&snippet.ID, &snippet.UserID, &snippet.Title, &snippet.Content, &snippet.Language, &snippet.CreatedAt)
	if err != nil {
		http.Error(w, "snippet not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(snippet)
}
