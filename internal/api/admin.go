package api

import (
	"fmt"
	"net/http"
)

func (cfg *Config) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(res, req)
	})
}

func (cfg *Config) metricsHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	fmt.Fprintf(res, `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())
}

func (cfg *Config) resetHandler(res http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		res.WriteHeader(403)
		return
	}
	err := cfg.db.ResetUsers(req.Context())
	if err != nil {
		res.WriteHeader(500)
		response := []byte("Failed Users Reset")
		res.Write(response)
		return
	}
	res.WriteHeader(200)
}
