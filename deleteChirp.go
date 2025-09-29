package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/xenonnn4w/chirpy/internal/auth"
	"github.com/xenonnn4w/chirpy/internal/database"
)

func (apiCfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting bearer token: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userID, err := auth.ValidateJWT(token, apiCfg.JWTSecret)
	if err != nil {
		log.Printf("Error Validating JWT: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	chirpIDstr := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDstr)
	if err != nil {
		log.Printf("Error parsing chirp ID: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	chirp, err := apiCfg.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		log.Printf("Error getting the chirp by ID: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if chirp.UserID != userID {
		log.Printf("User %s attempted to delete chirp %s owned by %s", userID, chirpID, chirp.UserID)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err = apiCfg.db.DeleteChirp(r.Context(), database.DeleteChirpParams{
		ID:     chirpID,
		UserID: userID,
	})
	if err != nil {
		log.Printf("Error deleting chirp: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
