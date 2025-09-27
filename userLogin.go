package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/xenonnn4w/chirpy/internal/auth"
	"github.com/xenonnn4w/chirpy/internal/database"
)

func (apiCfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	type User struct {
		ID           uuid.UUID `json:"id"`
		Email        string    `json:"email"`
		Password     string    `json:"password"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var p parameters
	err := decoder.Decode(&p)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := apiCfg.db.GetUserEmail(r.Context(), p.Email)
	if err != nil {
		fmt.Println("Error getting user:", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	match, err := auth.CheckPasswordHash(p.Password, user.Password)
	if err != nil {
		fmt.Println("Error checking password:", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !match {
		fmt.Println("Password does not match")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Create access token (1 hour)
	accessToken, err := auth.MakeJWT(user.ID, apiCfg.JWTSecret, time.Hour)
	if err != nil {
		fmt.Println("Error making JWT:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create refresh token (60 days)
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		fmt.Println("Error making refresh token:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Store refresh token in database
	_, err = apiCfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour), // 60 days
	})
	if err != nil {
		fmt.Println("Error creating refresh token:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        accessToken,
		RefreshToken: refreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (apiCfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type RefreshResponse struct {
		Token string `json:"token"`
	}

	// Get refresh token from header
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting bearer token: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get user from refresh token
	user, err := apiCfg.db.GetUserFromRefreshToken(r.Context(), tokenString)
	if err != nil {
		log.Printf("Error getting user from refresh token: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Create new access token
	accessToken, err := auth.MakeJWT(user.ID, apiCfg.JWTSecret, time.Hour)
	if err != nil {
		log.Printf("Error making JWT: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := RefreshResponse{
		Token: accessToken,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (apiCfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	// Get refresh token from header
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting bearer token: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Revoke the token
	err = apiCfg.db.RevokeRefreshToken(r.Context(), tokenString)
	if err != nil {
		log.Printf("Error revoking refresh token: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
