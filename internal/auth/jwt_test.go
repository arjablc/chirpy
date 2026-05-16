package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWTRoundTripsUserID(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"

	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT returned an error: %v", err)
	}

	gotUserID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT returned an error: %v", err)
	}

	if gotUserID != userID {
		t.Fatalf("ValidateJWT returned user ID %s, want %s", gotUserID, userID)
	}
}

func TestValidateJWTRejectsExpiredToken(t *testing.T) {
	token, err := MakeJWT(uuid.New(), "test-secret", -time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT returned an error: %v", err)
	}

	if _, err := ValidateJWT(token, "test-secret"); err == nil {
		t.Fatal("ValidateJWT returned no error for an expired token")
	}
}

func TestValidateJWTRejectsWrongSecret(t *testing.T) {
	token, err := MakeJWT(uuid.New(), "correct-secret", time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT returned an error: %v", err)
	}

	if _, err := ValidateJWT(token, "wrong-secret"); err == nil {
		t.Fatal("ValidateJWT returned no error for a token signed with a different secret")
	}
}
