package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"velvet-archive-api/internal/db"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	// Load env variables
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	// Initialize the database connection pool
	pool, err := db.InitPool(dbURL)
	if err != nil {
		log.Fatalf("Failed to bind database: %v", err)
	}
	defer pool.Close()

	// Wrap pool inside SQLC's generated query interface
	store := db.New(pool)
	_ = store

	// Initialize the Chi router
	r := chi.NewRouter()

	// Essential middleware (Logging, Recoverer catches panics so the server doesn't crash)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// CORS config
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// API v1 routing group
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			// Check if database is accessible during health checks
			err := pool.Ping(context.Background())
			if err != nil {
				http.Error(w, "Database connection unhealthy", http.StatusInternalServerError)
				return
			}
			w.Write([]byte("The Velvet Archive API and Database are operational."))
		})
	})

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Starting server on %s", addr)

	err = http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
