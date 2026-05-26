package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/arjablc/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *Config) polkaWebhookListener(w http.ResponseWriter, r *http.Request) {
	type incomingPayload struct {
		Event string `json:"event"`
		Data  struct {
			UserId string `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		errorResponse(w, "Apikey missing", 401)
		return
	}
	if apiKey != cfg.polkaKey {
		errorResponse(w, "Apikey missing", 401)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		errorResponse(w, "Body not found", 400)
		return
	}
	var bodyPayload incomingPayload
	if err := json.Unmarshal(body, &bodyPayload); err != nil {
		fmt.Println(err)
		errorResponse(w, "Failed Unmarshaling", 500)
		return
	}

	if bodyPayload.Event != "user.upgraded" {
		errorResponse(w, "", 204)
		return
	}
	parsedUid, err := uuid.Parse(bodyPayload.Data.UserId)
	if err != nil {
		errorResponse(w, "User not found", 404)
		return
	}
	_, err = cfg.db.EnableRedById(r.Context(), parsedUid)
	if err != nil {
		errorResponse(w, "User not found", 404)
		return
	}
	type empty struct {
	}
	respondJSON(w, 204, empty{})

}
