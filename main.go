package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
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
	res.Header().Add("Content-Type", "text/plain; charset=utf-8")
	data := []byte(fmt.Sprintf("Hits: %d", config.fileserverHits.Load()))
	res.Write(data)
}
func (config *apiConfig) resetHanlder(res http.ResponseWriter, req *http.Request) {
	config.fileserverHits.Store(0)
	res.WriteHeader(200)
}

func main() {
	servMux := http.ServeMux{}
	cfg := apiConfig{}
	port := 8080
	server := http.Server{Handler: &servMux, Addr: ":8080"}
	fsHandler := http.FileServer(http.Dir("."))
	servMux.Handle("/app/", http.StripPrefix("/app", cfg.middlewareMetricInc(fsHandler)))
	servMux.HandleFunc("GET /healthz", hanldeHealthStatus)
	servMux.HandleFunc("GET /metrics", cfg.metricHandler)
	servMux.HandleFunc("POST /reset", cfg.resetHanlder)
	fmt.Println("Serving and listening on port", port)
	server.ListenAndServe()
}
