package whisper

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Transcribe(inputPath string, volume float64) (string, error) {
	wavPath := strings.TrimSuffix(inputPath, filepath.Ext(inputPath)) + ".wav"
	defer os.Remove(wavPath)

	ffmpeg := exec.Command("ffmpeg",
		"-y",
		"-i", inputPath,
		"-ar", "16000",
		"-ac", "1",
		"-af", fmt.Sprintf("volume=%.1f", volume),
		"-f", "wav",
		wavPath,
	)

	if out, err := ffmpeg.CombinedOutput(); err != nil {
		return "", fmt.Errorf("ffmpeg: %s", string(out))
	}

	modelPath := os.Getenv("WHISPER_MODEL")
	if modelPath == "" {
		modelPath = "models/ggml-base.bin"
	}

	whisperBin := os.Getenv("WHISPER_BIN")
	if whisperBin == "" {
		whisperBin = "whisper.cpp/main"
	}

	cmd := exec.Command(whisperBin,
		"--model", modelPath,
		"--language", "ru",
		"--output-txt",
		"--no-prints",
		wavPath,
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("whisper err=%v output=%q", err, string(out))
	}

	txtPath := wavPath + ".txt"
	defer os.Remove(txtPath)

	result, err := os.ReadFile(txtPath)
	if err != nil {
		return "", fmt.Errorf("не удалось прочитать результат: %w", err)
	}

	return strings.TrimSpace(string(result)), nil
}
