package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Ranil-Valeev/voice-to-text/internal/handler"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})

	mux.HandleFunc("/transcribe", handler.Transcribe)

	log.Println("Сервер запущен на http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
