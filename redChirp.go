package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/xenonnn4w/chirpy/internal/auth"
	"github.com/xenonnn4w/chirpy/internal/database"
)

func (apiCfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		log.Printf("Error getting API key: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if apiKey != apiCfg.polkaWebhookSecret {
		log.Printf("Invalid API key: %s", apiKey)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	var params parameters
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		log.Printf("Error decoding JSON: %v", err)
		w.WriteHeader(http.StatusBadRequest)
	}

	if params.Event != "user.upgraded" {
		log.Printf("Expected user.upgraded event, got %s", params.Event)
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		log.Printf("Error parsing user ID: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("Updating user %s to is_chirpy_red: true", userID)
	user, err := apiCfg.db.UpdateUserIsChirpyRed(r.Context(), database.UpdateUserIsChirpyRedParams{
		IsChirpyRed: sql.NullBool{
			Bool:  true,
			Valid: true,
		},
		ID: userID,
	})
	if err != nil {
		log.Printf("Error updating user is chirpy red: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	log.Printf("User %s is now chirpy red", user.ID)
	w.WriteHeader(http.StatusNoContent)
}
