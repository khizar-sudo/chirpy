package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/khizar-sudo/chirpy/internal/database"
)

func Init() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             database.New(db),
		platform:       platform,
	}
	mux := http.NewServeMux()

	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /admin/metrics", cfg.getMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.resetMetrics)

	mux.HandleFunc("GET /api/healthz", healthCheck)
	mux.HandleFunc("POST /api/users", createUser(&cfg))
	mux.HandleFunc("POST /api/login", login(&cfg))
	mux.HandleFunc("POST /api/chirps", createChirp(&cfg))
	mux.HandleFunc("GET /api/chirps", getAllChirps(&cfg))
	mux.HandleFunc("GET /api/chirps/{chirpID}", getChirp(&cfg))

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
	log.Fatal(server.ListenAndServe())
}
