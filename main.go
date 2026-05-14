package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/arjablc/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
}

func (config *apiConfig) middlewareMetricInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		fmt.Println("From inside the middleware")
		fmt.Println("Before value of atomic:", config.fileserverHits.Load())
		config.fileserverHits.Add(1)
		fmt.Println("After value of atomic:", config.fileserverHits.Load())
		next.ServeHTTP(res, req)
	})
}

func (config *apiConfig) metricHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(200)
	res.Header().Add("Content-Type", "text/html; charset=utf-8")
	data := []byte(fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, config.fileserverHits.Load()))
	res.Write(data)
}
func (config *apiConfig) resetHanlder(res http.ResponseWriter, req *http.Request) {
	config.fileserverHits.Store(0)
	res.WriteHeader(200)
}

func main() {
	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("Failed to open db")
	}
	dbQueries := database.New(db)
	servMux := http.ServeMux{}
	cfg := apiConfig{db: dbQueries}
	port := 8080
	server := http.Server{Handler: &servMux, Addr: ":8080"}
	fsHandler := http.FileServer(http.Dir("."))
	//NOTE: app route
	servMux.Handle("/app/", http.StripPrefix("/app", cfg.middlewareMetricInc(fsHandler)))
	//NOTE: api routes
	servMux.HandleFunc("GET /api/healthz", hanldeHealthStatus)
	servMux.HandleFunc("POST /api/validate_chirp", vallidateChirp)
	//NOTE: admin route
	servMux.HandleFunc("POST /admin/reset", cfg.resetHanlder)
	servMux.HandleFunc("GET /admin/metrics", cfg.metricHandler)
	fmt.Println("Serving and listening on port", port)
	server.ListenAndServe()
}
