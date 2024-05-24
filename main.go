package main

import (
	"database/sql"
	"fmt"
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

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello, world\n")
	})

	mux.HandleFunc("POST /v1/users", cfg.handlerCreateUsers)

	mux.HandleFunc("/v1/readiness", cfg.handlerReadiness)
	mux.HandleFunc("/v1/error", cfg.handlerError)

	log.Fatal(srv.ListenAndServe())
}
