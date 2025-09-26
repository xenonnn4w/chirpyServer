package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func (cfg *apiConfig) handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type Chirp struct {
		Body string `json:"body"`
	}
	type ChirpError struct {
		Error string `json:"error"`
	}
	type validChirp struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	var chirp Chirp
	err := decoder.Decode(&chirp)
	if err != nil {
		log.Printf("Error decoding JSON: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChirpError{Error: "Error decoding JSON"})
	}
	if len(chirp.Body) > 140 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ChirpError{Error: "Chirp is too long"})
	}

	w.WriteHeader(http.StatusOK)
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(chirp.Body, " ")

	for i, word := range words {
		if contains(badWords, strings.ToLower(word)) {
			words[i] = "****"
		}
	}
	cleanedChirp := strings.Join(words, " ")
	json.NewEncoder(w).Encode(validChirp{CleanedBody: cleanedChirp})
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
