package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/arjablc/chirpy/internal/api"
	"github.com/arjablc/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("failed to load .env: %v", err)
	}

	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	jwtSecret := os.Getenv("JWT_SECRET")
	polkaKey := os.Getenv("POLKA_KEY")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Db connect err: %v", err)
		panic("failed to open db")
	}
	defer db.Close()
	dbQueries := database.New(db)
	cfg := api.NewConfig(dbQueries, platform, jwtSecret, polkaKey)

	port := "8080"
	server := http.Server{
		Handler: cfg.Router("."),
		Addr:    ":" + port,
	}

	fmt.Println("Serving and listening on port", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
