package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/arjablc/chirpy/internal/auth"
	"github.com/arjablc/chirpy/internal/database"
)

func (cfg *Config) refreshToken(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Inside the refresh Token method")
	type tokenResponse struct {
		Token string `json:"token"`
	}

	fmt.Println("Getting bearer")
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		fmt.Println("error getting bearer")
		fmt.Print(err)
		errorResponse(w, "Unauthorized", 401)
		return
	}

	fmt.Println("Getting refresh token from db")
	token, err := cfg.db.GetRefreshTokenByToken(r.Context(), refreshToken)
	if err != nil {
		fmt.Println("error getting refersh token from db")
		fmt.Print(err)
		errorResponse(w, "Unauthorized", 401)
		return
	}
	fmt.Println("Checking revoked")
	if token.RevokedAt.Valid {
		fmt.Print("Invalid token (expired)")
		errorResponse(w, "Unauthorized", 401)
		return
	}
	fmt.Println("Getting user for token")
	user, err := cfg.db.GetUserByRefreshToken(r.Context(), refreshToken)
	if err != nil {
		fmt.Println("error getting user from db")
		fmt.Print(err)
		errorResponse(w, "Unauthorized", 401)
		return
	}
	fmt.Println("Creating a token")
	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		fmt.Println("error creating jwt token")
		fmt.Print(err)
		errorResponse(w, "Unauthorized", 401)
		return
	}
	respondJSON(w, 200, tokenResponse{
		Token: accessToken,
	})

}

func (cfg *Config) revokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		fmt.Print(err)
		errorResponse(w, "Unauthorized", 401)
	}
	revokeParms := database.RevokeRefreshTokenByTokenParams{
		Token:     refreshToken,
		RevokedAt: sql.NullTime{Valid: true, Time: time.Now()},
	}
	err = cfg.db.RevokeRefreshTokenByToken(r.Context(), revokeParms)
	if err != nil {
		fmt.Print(err)
		errorResponse(w, "Unauthorized", 401)
	}
	w.WriteHeader(204)
}
