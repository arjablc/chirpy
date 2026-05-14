package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"unicode/utf8"
)

var profanes = map[string]string{"kerfuffle": "****", "sharbert": "****", "fornax": "****"}

func hanldeHealthStatus(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(200)
	responseBody := []byte("OK")
	res.Write(responseBody)
}

func vallidateChirp(res http.ResponseWriter, req *http.Request) {
	type incoming struct {
		Body string `json:"body"`
	}
	type validRes struct {
		CleanedBody string `json:"cleaned_body"`
	}
	reqBody := incoming{}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Fatalf("Body read error")
		_errorResponse(res, "Internal Error", 500)
		return
	}
	err = json.Unmarshal(body, &reqBody)
	if err != nil {
		log.Fatalf("Unmarshal error")
		_errorResponse(res, "Internal Error", 500)
		return
	}
	if utf8.RuneCountInString(reqBody.Body) > 140 {
		_errorResponse(res, "Chirp is too long", 400)
		return
	}
	cleanedBody := _cleanBody(reqBody.Body)
	res.WriteHeader(200)
	responseBody := validRes{CleanedBody: cleanedBody}
	response, err := json.Marshal(responseBody)
	if err != nil {
		log.Fatalf("Body marshal error")
	}
	res.Write(response)
	return
}
func _cleanBody(body string) string {
	words := strings.Split(body, " ")
	for index, word := range words {
		lowerWord := strings.ToLower(word)
		if val, ok := profanes[lowerWord]; ok {
			words[index] = val
		}
	}
	return strings.Join(words, " ")
}

func _errorResponse(res http.ResponseWriter, errorString string, errorCode int) {
	type errorRes struct {
		Error string `json:"error"`
	}
	res.WriteHeader(errorCode)
	responseBody := errorRes{Error: "Chirp is too long"}
	response, err := json.Marshal(responseBody)
	if err != nil {
		log.Fatalf("Body marshal error")
	}
	res.Write(response)
	return

}
