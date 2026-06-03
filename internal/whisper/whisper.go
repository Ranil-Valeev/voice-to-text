package whisper

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Transcribe принимает путь к аудиофайлу (любой формат),
// конвертирует его в wav через ffmpeg и запускает whisper.cpp
func Transcribe(inputPath string) (string, error) {
	wavPath := strings.TrimSuffix(inputPath, filepath.Ext(inputPath)) + ".wav"
	defer os.Remove(wavPath)

	// Шаг 1: конвертируем в wav 16kHz моно — именно такой формат нужен whisper
	ffmpeg := exec.Command("ffmpeg",
		"-y", // перезаписать если существует
		"-i", inputPath,
		"-ar", "16000", // частота дискретизации 16kHz
		"-ac", "2", // моно
		"-f", "wav",
		wavPath,
	)

	if out, err := ffmpeg.CombinedOutput(); err != nil {
		return "", fmt.Errorf("whisper err=%v output=%q", err, string(out))
	}

	// Шаг 2: запускаем whisper.cpp
	// --model  — путь к модели (положи рядом с проектом)
	// --output-txt — сохранить результат в файл .txt
	// --language ru — русский язык (поменяй если нужен другой)
	modelPath := os.Getenv("WHISPER_MODEL")
	if modelPath == "" {
		modelPath = "models/ggml-base.bin" // путь по умолчанию
	}

	whisperBin := os.Getenv("WHISPER_BIN")
	if whisperBin == "" {
		whisperBin = "whisper.cpp/main" // путь по умолчанию
	}

	cmd := exec.Command(whisperBin,
		"--model", modelPath,
		"--language", "ru",
		"--output-txt",
		"--no-prints",
		wavPath,
	)

	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("whisper: %s", string(out))
	}

	// whisper.cpp сохраняет результат в файл wavPath.txt
	txtPath := wavPath + ".txt"
	defer os.Remove(txtPath)

	result, err := os.ReadFile(txtPath)
	if err != nil {
		return "", fmt.Errorf("не удалось прочитать результат: %w", err)
	}

	return strings.TrimSpace(string(result)), nil
}
