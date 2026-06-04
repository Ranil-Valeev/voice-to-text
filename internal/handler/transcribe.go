package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Ranil-Valeev/voice-to-text/internal/whisper"
)

type response struct {
	Text  string `json:"text,omitempty"`
	Error string `json:"error,omitempty"`
}

func Transcribe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(response{Error: "только POST"})
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 100<<20)

	file, header, err := r.FormFile("audio")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response{Error: "не удалось прочитать файл: " + err.Error()})
		return
	}
	defer file.Close()

	// Читаем громкость из формы, по умолчанию 1.0
	volume := 1.0
	if v := r.FormValue("volume"); v != "" {
		if parsed, err := strconv.ParseFloat(v, 64); err == nil && parsed > 0 {
			volume = parsed
		}
	}

	log.Printf("получен файл: %s (%.2f MB), усиление: %.1fx", header.Filename, float64(header.Size)/1024/1024, volume)

	tmpInput := filepath.Join("tmp", header.Filename)
	dst, err := os.Create(tmpInput)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response{Error: "не удалось сохранить файл"})
		return
	}
	defer os.Remove(tmpInput)
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response{Error: "ошибка записи файла"})
		return
	}
	dst.Close()

	text, err := whisper.Transcribe(tmpInput, volume)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response{Error: fmt.Sprintf("ошибка транскрибации: %s", err.Error())})
		return
	}

	json.NewEncoder(w).Encode(response{Text: text})
}
