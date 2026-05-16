package auth

import "testing"

func TestHashPasswordMatchesOriginalPassword(t *testing.T) {
	password := "correct-horse-battery-staple"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned an error: %v", err)
	}

	if hash == password {
		t.Fatal("HashPassword returned the plaintext password")
	}

	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash returned an error: %v", err)
	}
	if !match {
		t.Fatal("CheckPasswordHash returned false for the original password")
	}
}

func TestCheckPasswordHashRejectsWrongPassword(t *testing.T) {
	hash, err := HashPassword("correct-password")
	if err != nil {
		t.Fatalf("HashPassword returned an error: %v", err)
	}

	match, err := CheckPasswordHash("wrong-password", hash)
	if err != nil {
		t.Fatalf("CheckPasswordHash returned an error: %v", err)
	}
	if match {
		t.Fatal("CheckPasswordHash returned true for the wrong password")
	}
}
