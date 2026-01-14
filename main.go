package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/khizar-sudo/chirpy/internal/config"
	"github.com/khizar-sudo/chirpy/internal/database"
	"github.com/khizar-sudo/chirpy/internal/handlers"
	"github.com/khizar-sudo/chirpy/internal/middleware"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}

	cfg := &config.ApiConfig{
		DB: database.New(db),
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", middleware.MetricsInc(cfg)(http.StripPrefix("/app", http.FileServer(http.Dir("./static")))))
	mux.HandleFunc("GET /api/healthz", handlers.HealthCheck)
	mux.HandleFunc("GET /admin/metrics", handlers.GetMetrics(cfg))
	mux.HandleFunc("POST /admin/reset", handlers.ResetMetrics(cfg))
	mux.HandleFunc("POST /api/validate_chirp", handlers.ValidateChirp)

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
	log.Fatal(server.ListenAndServe())
}
