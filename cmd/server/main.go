package main

import (
	"log"
	"net/http"
	"os"

	"github.com/abhinavpaste/crev/internal/db"
	"github.com/abhinavpaste/crev/internal/handler"
	"github.com/abhinavpaste/crev/internal/middleware"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system env")
	} else {
		log.Println("loaded .env:", os.Getenv("DATABASE_URL"))
	}

	db.Init()

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	mux.HandleFunc("/register", handler.Register)
	mux.HandleFunc("/login", handler.Login)
	mux.HandleFunc("/snippets/", handler.GetSnippet)
	mux.HandleFunc("/snippets", middleware.Auth(handler.CreateSnippet))
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	mux.HandleFunc("/snippets/{id}/comments", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			middleware.Auth(handler.CreateComment)(w, r)
		} else {
			handler.GetComments(w, r)
		}
	})

	log.Println("crev running on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
