package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondJSON(res http.ResponseWriter, statusCode int, payload any) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statusCode)
	if err := json.NewEncoder(res).Encode(payload); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func errorResponse(res http.ResponseWriter, errorString string, errorCode int) {
	type errorRes struct {
		Error string `json:"error"`
	}
	respondJSON(res, errorCode, errorRes{Error: errorString})
}
