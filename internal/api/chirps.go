package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/arjablc/chirpy/internal/auth"
	"github.com/arjablc/chirpy/internal/database"
	"github.com/google/uuid"
)

var profanes = map[string]string{
	"kerfuffle": "****",
	"sharbert":  "****",
	"fornax":    "****",
}

type chirpPayload struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func (C *Config) createChirp(res http.ResponseWriter, req *http.Request) {
	type incomingBody struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		errorResponse(res, "Failed to read body", 500)
		return
	}
	var reqBody incomingBody
	err = json.Unmarshal(body, &reqBody)
	if err != nil {
		errorResponse(res, "Failed to unmarsal body", 500)
		return
	}
	bearerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		fmt.Println("No bearer", err)
		errorResponse(res, "Unauthorized", 401)
		return
	}
	uid, err := auth.ValidateJWT(bearerToken, C.jwtSecret)
	if err != nil {
		fmt.Println("validation error", err)
		errorResponse(res, "Unauthorized", 401)
		return
	}

	createChirpParams := database.CreateChirpParams{Body: reqBody.Body, UserID: uid}
	dbChirp, err := C.db.CreateChirp(req.Context(), createChirpParams)
	if err != nil {
		errorResponse(res, "Failed to create chirp", 500)
		return
	}
	respondJSON(res, 201, chirpPayload{
		Id:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserId:    dbChirp.UserID,
	})
}

func (cfg *Config) chirpsById(resw http.ResponseWriter, req *http.Request) {
	chirpId := req.PathValue("id")
	parsedId, err := uuid.Parse(chirpId)
	if err != nil {
		errorResponse(resw, "Invalid Id", 401)
		return
	}
	chirp, err := cfg.db.GetChirpsById(req.Context(), parsedId)
	if err != nil {
		errorResponse(resw, "Not found", 404)
		return
	}
	respondJSON(resw, 200, chirpPayload{
		Id:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	})

}

func (cfg *Config) chirpsOrderedByCreatedAt(writer http.ResponseWriter, req *http.Request) {
	chirps, err := cfg.db.GetChirpsOrderedBy(req.Context())
	if err != nil {
		errorResponse(writer, "Failed to get Chirps", 500)
		return
	}
	if len(chirps) < 1 {
		errorResponse(writer, "No chirps found", 404)
		return
	}
	resChirps := []chirpPayload{}
	for _, chirp := range chirps {
		resChirps = append(resChirps, chirpPayload{
			Id:        chirp.ID,
			Body:      chirp.Body,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			UserId:    chirp.UserID,
		})
	}
	respondJSON(writer, 200, resChirps)
}

func validateChirp(res http.ResponseWriter, req *http.Request) {
	type incoming struct {
		Body string `json:"body"`
	}
	type validRes struct {
		CleanedBody string `json:"cleaned_body"`
	}
	reqBody := incoming{}
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		errorResponse(res, "Invalid request body", http.StatusBadRequest)
		return
	}

	if utf8.RuneCountInString(reqBody.Body) > 140 {
		errorResponse(res, "Chirp is too long", http.StatusBadRequest)
		return
	}
	respondJSON(res, http.StatusOK, validRes{CleanedBody: cleanBody(reqBody.Body)})
}

func cleanBody(body string) string {
	words := strings.Split(body, " ")
	for index, word := range words {
		if val, ok := profanes[strings.ToLower(word)]; ok {
			words[index] = val
		}
	}
	return strings.Join(words, " ")
}

func (cfg *Config) deleteChirpsById(resw http.ResponseWriter, req *http.Request) {
	bearer, err := auth.GetBearerToken(req.Header)
	if err != nil {
		errorResponse(resw, "Unauthorized", 401)
		return
	}
	uid, err := auth.ValidateJWT(bearer, cfg.jwtSecret)
	if err != nil {
		errorResponse(resw, "Unauthorized", 403)
		return
	}
	chirpId := req.PathValue("id")
	parsedId, err := uuid.Parse(chirpId)
	if err != nil {
		errorResponse(resw, "Invalid Id", 401)
		return
	}
	chirp, err := cfg.db.GetChirpsById(req.Context(), parsedId)
	if err != nil {
		errorResponse(resw, "Not found", 404)
		return
	}
	if chirp.UserID != uid {
		errorResponse(resw, "Unauthorized", 403)
		return
	}
	err = cfg.db.DeleteChirpById(req.Context(), chirp.ID)
	if err != nil {
		errorResponse(resw, "Failed deletion", 500)
		return
	}
	resw.WriteHeader(204)
}
