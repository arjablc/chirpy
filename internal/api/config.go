package api

import (
	"net/http"
	"sync/atomic"

	"github.com/arjablc/chirpy/internal/database"
)

type Config struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

func NewConfig(db *database.Queries, platform string) *Config {
	return &Config{db: db, platform: platform}
}

func (cfg *Config) Router(staticDir string) http.Handler {
	mux := http.NewServeMux()
	fsHandler := http.FileServer(http.Dir(staticDir))

	mux.Handle("/app/", http.StripPrefix("/app", cfg.middlewareMetricsInc(fsHandler)))
	mux.HandleFunc("GET /api/healthz", handleHealthStatus)
	mux.HandleFunc("POST /api/users", cfg.createUser)
	mux.HandleFunc("POST /api/login", cfg.login)
	// chirps
	mux.HandleFunc("POST /api/chirps", cfg.createChirp)
	mux.HandleFunc("GET /api/chirps", cfg.chirpsOrderedByCreatedAt)
	mux.HandleFunc("GET /api/chirps/{id}", cfg.chirpsById)

	mux.HandleFunc("POST /admin/reset", cfg.resetHandler)
	mux.HandleFunc("GET /admin/metrics", cfg.metricsHandler)

	return mux
}
