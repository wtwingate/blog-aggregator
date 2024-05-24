package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")

	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello, world\n")
	})

	mux.HandleFunc("/v1/readiness", handlerReadiness)
	mux.HandleFunc("/v1/error", handlerError)

	log.Fatal(srv.ListenAndServe())
}
