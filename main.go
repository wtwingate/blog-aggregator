package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/wtwingate/blog-aggregator/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	conn := os.Getenv("CONN")

	db, err := sql.Open("postgres", conn)
	if err != nil {
		log.Fatal(err)
	}

	dbQueries := database.New(db)
	cfg := apiConfig{
		DB: dbQueries,
	}

	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.HandleFunc("POST /v1/users", cfg.handlerCreateUsers)
	mux.HandleFunc("GET /v1/users", cfg.middlewareAuth(cfg.handlerGetUsers))

	mux.HandleFunc("POST /v1/feeds", cfg.middlewareAuth(cfg.handlerCreateFeeds))
	mux.HandleFunc("GET /v1/feeds", cfg.handlerGetFeeds)

	mux.HandleFunc("POST /v1/feed_follows", cfg.middlewareAuth(cfg.handlerCreateFeedFollows))
	mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", cfg.middlewareAuth(cfg.handlerDeleteFeedFollows))

	mux.HandleFunc("/v1/readiness", cfg.handlerReadiness)

	mux.HandleFunc("/v1/error", cfg.handlerError)

	log.Fatal(srv.ListenAndServe())
}
