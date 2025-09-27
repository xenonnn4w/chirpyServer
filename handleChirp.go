package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/xenonnn4w/chirpy/internal/database"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string `json:"body"`
		UserID string `json:"user_id"`
	}
	type ChirpResponse struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding JSON %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userID, err := uuid.Parse(params.UserID)
	if err != nil {
		log.Printf("Error parsing user ID %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   params.Body,
		UserID: userID,
	})
	if err != nil {
		log.Printf("Error creating chirp %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	chirpResponse := ChirpResponse{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(chirpResponse)
	log.Printf("Chirp created: %s", chirp.Body)

}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	type Chirp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		log.Printf("Error getting chirps %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	chirpResponses := []Chirp{}
	for _, chirp := range chirps {
		chirpResponses = append(chirpResponses, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chirpResponses)
}

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	type ChirpResponse struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}
	chirpIDstr := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDstr)
	if err != nil {
		log.Printf("Error parsing chirp ID %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	chirp, err := cfg.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		log.Printf("Error getting chirp by ID %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	chirpResponse := ChirpResponse{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chirpResponse)
}
