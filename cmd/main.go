package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Ranil-Valeev/voice-to-text/internal/handler"
	"github.com/Ranil-Valeev/voice-to-text/internal/storage"
)

func main() {
	db, err := storage.New()
	if err != nil {
		log.Fatalf("база данных: %v", err)
	}

	if err := db.Migrate(); err != nil {
		log.Fatalf("миграция: %v", err)
	}

	log.Println("база данных подключена")

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	mux.HandleFunc("/transcribe", handler.Transcribe(db))
	mux.HandleFunc("/history", handler.History(db))

	log.Println("сервер запущен на http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
