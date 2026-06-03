# voice-to-text

REST API на Go для транскрибации аудиофайлов в текст. Работает локально, без внешних сервисов — используется [whisper.cpp](https://github.com/ggerganov/whisper.cpp).

## Стек

- Go — HTTP-сервер
- whisper.cpp — локальная модель распознавания речи
- ffmpeg — конвертация аудио в нужный формат

## Запуск

### 1. Зависимости

- [Go 1.21+](https://go.dev/dl/)
- [ffmpeg](https://ffmpeg.org/download.html) — добавить в PATH
- [whisper.cpp](https://github.com/ggerganov/whisper.cpp) — собрать из исходников

```bash
git clone https://github.com/ggerganov/whisper.cpp
cd whisper.cpp
make
bash models/download-ggml-model.sh base
```

### 2. Переменные окружения

```bash
export WHISPER_BIN=./whisper.cpp/build/bin/Release/whisper-cli.exe
export WHISPER_MODEL=./whisper.cpp/models/ggml-base.bin
```

### 3. Запуск сервера

```bash
go run cmd/main.go
# Сервер запущен на http://localhost:8080
```

## Использование

**POST /transcribe** — отправить аудиофайл, получить текст

```bash
curl -X POST http://localhost:8080/transcribe \
  -F "audio=@запись.m4a"
```

Ответ:

```json
{
  "text": "распознанный текст"
}
```

Поддерживаемые форматы: m4a, mov, mp3, wav, ogg и любые другие которые поддерживает ffmpeg.

**GET /health** — проверка что сервер живой

```bash
curl http://localhost:8080/health
```

## Модели

| Модель | Размер | Точность |
|--------|--------|----------|
| tiny   | 75MB   | низкая   |
| base   | 141MB  | средняя  |
| small  | 466MB  | хорошая  |

По умолчанию используется `base`. Для лучшей точности `small`
