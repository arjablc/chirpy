package main

import "net/http"

func hanldeHealthStatus(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(200)
	responseBody := []byte("OK")
	res.Write(responseBody)
}
