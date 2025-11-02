# **TranscribeIT**

Local web app to transcribe audio in text. Uses _Kafka_ to manage audio transcriptions.

## 🚀 Features

- 🎧 Audio-to-text transcription using local model
- ⚙️ Scalable architecture powered by Apache Kafka
- 🖥️ Web interface built with Next.js
- 🐍 Python transcription model inference (Whisper)
- 💬 Multi-language support

## Dependencies

- python: version 3.12
  - ffmpeg: version 5/6/7
- go: version 1.23 (backend)
- nextjs (frontend)

## 🧪 Setup

Clone the repository & open it

`git clone https://github.com/imvalerio/transcribeit.git`

`cd transcribeit`

Start Kafka:

`docker-compose up -d`

Run the backend:

`go run main.go`

Run the frontend:

`cd frontend`

Install dependencies:

`npm install`

Run NextJs locally:

`npm run dev`

## Python script

Could be used stand-alone, takes in input some arguments:

1. input file path (required)
2. output file path (required)
3. model type (optional, default: turbo)
4. language (optional)
