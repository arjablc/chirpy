package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeaderVal := headers.Get("Authorization")
	if authHeaderVal == "" {
		return "", errors.New("No auth headers")
	}
	split := strings.Fields(authHeaderVal)
	if len(split) != 2 || split[0] != "Bearer" {
		return "", errors.New("Malformed auth header")
	}
	return split[1], nil

}
