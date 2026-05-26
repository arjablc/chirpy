package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/arjablc/chirpy/internal/auth"
	"github.com/arjablc/chirpy/internal/database"
	"github.com/google/uuid"
)

const (
	refreshTokenTTL = time.Hour * 24 * 60
)

type userReqPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type userResPayload struct {
	Id          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

type loginResPayload struct {
	Id           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
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
		Id:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}

func (C *Config) updateUser(res http.ResponseWriter, req *http.Request) {
	bearer, err := auth.GetBearerToken(req.Header)
	if err != nil {
		errorResponse(res, "Unauthorized", 401)
		return
	}
	userId, err := auth.ValidateJWT(bearer, C.jwtSecret)
	if err != nil {
		errorResponse(res, "Unauthorized", 401)
		return
	}
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
	hashed, err := auth.HashPassword(requestBody.Password)
	if err != nil {
		errorResponse(res, "Failed to hash", 500)
		return
	}
	updateUserParm := database.UpdateEmailPwByIdParams{
		ID:             userId,
		Email:          requestBody.Email,
		HashedPassword: hashed,
	}
	user, err := C.db.UpdateEmailPwById(req.Context(), updateUserParm)
	if err != nil {
		errorResponse(res, "Failed to create User", 500)
		return
	}
	respondJSON(res, 200, userResPayload{
		Id:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}

func (cfg *Config) loginUser(w http.ResponseWriter, r *http.Request) {
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
	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	refreshToken := auth.MakeRefreshToken()
	refreshTokenParam := database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(refreshTokenTTL),
		RevokedAt: sql.NullTime{},
	}
	_, err = cfg.db.CreateRefreshToken(r.Context(), refreshTokenParam)
	if err != nil {
		fmt.Println(err)
		errorResponse(w, "Couldn't make refresh token", 401)
		return
	}

	if err != nil {
		fmt.Println(err)
		errorResponse(w, "Couldn't make jwt", 401)
		return
	}

	respondJSON(w, 200, loginResPayload{
		Id:           user.ID,
		Email:        user.Email,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Token:        token,
		IsChirpyRed:  user.IsChirpyRed,
		RefreshToken: refreshToken,
	})

}
