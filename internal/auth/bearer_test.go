package auth

import (
	"net/http"
	"testing"
)

func TestGetBearerTokenReturnsToken(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer test-token")

	token, err := GetBearerToken(headers)
	if err != nil {
		t.Fatalf("GetBearerToken returned an error: %v", err)
	}

	if token != "test-token" {
		t.Fatalf("GetBearerToken returned %q, want %q", token, "test-token")
	}
}

func TestGetBearerTokenTrimsTokenWhitespace(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer   test-token   ")

	token, err := GetBearerToken(headers)
	if err != nil {
		t.Fatalf("GetBearerToken returned an error: %v", err)
	}

	if token != "test-token" {
		t.Fatalf("GetBearerToken returned %q, want %q", token, "test-token")
	}
}

func TestGetBearerTokenRejectsMissingAuthorizationHeader(t *testing.T) {
	if _, err := GetBearerToken(http.Header{}); err == nil {
		t.Fatal("GetBearerToken returned no error for missing Authorization header")
	}
}

func TestGetBearerTokenRejectsMalformedAuthorizationHeader(t *testing.T) {
	tests := []struct {
		name       string
		authHeader string
	}{
		{
			name:       "missing bearer scheme",
			authHeader: "test-token",
		},
		{
			name:       "wrong scheme",
			authHeader: "Basic test-token",
		},
		{
			name:       "empty bearer token",
			authHeader: "Bearer ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := http.Header{}
			headers.Set("Authorization", tt.authHeader)

			if _, err := GetBearerToken(headers); err == nil {
				t.Fatal("GetBearerToken returned no error for malformed Authorization header")
			}
		})
	}
}
