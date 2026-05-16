package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/arjablc/chirpy/internal/auth"
	"github.com/arjablc/chirpy/internal/database"
	"github.com/google/uuid"
)

type userReqPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type userResPayload struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (C *Config) createUser(res http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		errorResponse(res, "Failed to read body", 500)
		return
	}
	var requestBody userReqPayload
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		errorResponse(res, "Invalid request body", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(requestBody.Email) == "" {
		errorResponse(res, "Email is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(requestBody.Password) == "" {
		errorResponse(res, "Password is required", http.StatusBadRequest)
		return
	}
	hashed, err := auth.HashPassword(requestBody.Password)
	if err != nil {
		errorResponse(res, "Failed to hash", 500)
		return
	}

	userDto := database.CreateUserParams{
		Email:          requestBody.Email,
		HashedPassword: hashed,
	}
	user, err := C.db.CreateUser(req.Context(), userDto)
	if err != nil {
		errorResponse(res, "Failed to create User", 500)
		return
	}
	respondJSON(res, 201, userResPayload{
		Id:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}

func (cfg *Config) login(w http.ResponseWriter, r *http.Request) {
	const invalidCredMsg = "Incorrect email or password"
	body, err := io.ReadAll(r.Body)
	if err != nil {
		errorResponse(w, "Failed to read body", 500)
		return
	}
	var reqPayload userReqPayload
	err = json.Unmarshal(body, &reqPayload)
	if err != nil {
		errorResponse(w, "Failed to unmarshal", 500)
		return
	}
	user, err := cfg.db.UserByEmail(r.Context(), reqPayload.Email)
	if err != nil {
		errorResponse(w, invalidCredMsg, 401)
		return
	}
	match, err := auth.CheckPasswordHash(reqPayload.Password, user.HashedPassword)
	if !match || err != nil {
		errorResponse(w, invalidCredMsg, 401)
		return
	}
	respondJSON(w, 200, userResPayload{
		Id:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})

}
