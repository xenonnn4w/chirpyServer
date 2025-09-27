package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestCreateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"

	token, err := MakeJWT(userID, secret, time.Hour*24)
	if err != nil {
		t.Fatalf("Error creating JWT: %v", err)
	}

	if token == "" {
		t.Fatal("Token should not be empty")
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"

	// Create token
	token, err := MakeJWT(userID, secret, time.Hour*24)
	if err != nil {
		t.Fatalf("Error creating JWT: %v", err)
	}

	// Validate token
	validatedUserID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("Error validating JWT: %v", err)
	}

	if validatedUserID != userID {
		t.Fatalf("Expected user ID %v, got %v", userID, validatedUserID)
	}
}

func TestValidateJWT_WrongSecret(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"
	wrongSecret := "wrong-secret"

	// Create token with correct secret
	token, err := MakeJWT(userID, secret, time.Hour*24)
	if err != nil {
		t.Fatalf("Error creating JWT: %v", err)
	}

	// Try to validate with wrong secret
	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Fatal("Expected error when validating with wrong secret")
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"

	// Create token with very short expiration
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(-time.Hour)), // Expired 1 hour ago
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("Error creating expired JWT: %v", err)
	}

	// Try to validate expired token
	_, err = ValidateJWT(tokenString, secret)
	if err == nil {
		t.Fatal("Expected error when validating expired token")
	}
}

func TestGetBearerToken(t *testing.T) {
	// Test valid bearer token
	headers := http.Header{}
	headers.Set("Authorization", "Bearer valid-token-here")

	token, err := GetBearerToken(headers)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if token != "valid-token-here" {
		t.Fatalf("Expected 'valid-token-here', got '%s'", token)
	}
}

func TestGetBearerToken_NoHeader(t *testing.T) {
	headers := http.Header{}

	_, err := GetBearerToken(headers)
	if err == nil {
		t.Fatal("Expected error for missing header")
	}
}

func TestGetBearerToken_InvalidFormat(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "InvalidFormat")

	_, err := GetBearerToken(headers)
	if err == nil {
		t.Fatal("Expected error for invalid format")
	}
}
