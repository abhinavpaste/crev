package main

import (
	"context"
	"log"
	"os"

	"github.com/abhinavpaste/crev/internal/handler"
	"github.com/abhinavpaste/crev/internal/middleware"
	"github.com/abhinavpaste/crev/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	db, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("cannot connect to db: %v", err)
	}
	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		log.Fatalf("db ping failed: %v", err)
	}
	log.Println("connected to database")

	s := store.New(db)

	authHandler := handler.NewAuthHandler(s)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// public routes
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	// protected routes (snippet + comment handlers go here)
	protected := r.Group("/")
	protected.Use(middleware.AuthRequired())
	_ = protected

	port := os.Getenv("PORT")
	log.Printf("server starting on :%s", port)
	log.Fatal(r.Run(":" + port))
}
