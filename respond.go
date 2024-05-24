package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, status int, msg string) {
	log.Printf("responding with status code %v: %v\n", status, msg)

	payload := struct {
		Error string `json:"error"`
	}{
		Error: msg,
	}

	respondWithJSON(w, status, payload)
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Add("Content-type", "application/json")

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("could not marshal reponse: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	w.Write(data)
}
