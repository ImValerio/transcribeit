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

```bash
git clone https://github.com/imvalerio/transcribeit.git
cd transcribeit
```

### Start Kafka

```bash
docker-compose up -d
```

### Run the backend

Create .env file and define two variables:

- **VOICE_RECOGNITION_FOLDER**, where the script to run is located
- **PYTHON_EXEC_PATH**, abs path of the python which will run the script

Ex:

```bash
VOICE_RECOGNITION_FOLDER=C:\Users\USER\Documents\transcribeit\voice-recognition
PYTHON_EXEC_PATH=C:\Users\USER\Documents\transcribeit\voice-recognition\venv\Scripts\python.exe
```

Run the web server

```bash
go run main.go
```

### Run the frontend:

```bash
cd frontend
```

Install dependencies:

```bash
npm install
```

Run NextJs locally:

```bash
npm run dev
```

## Python script

It could be used stand-alone, takes in input some arguments:

1. input file path (required)
2. output file path (required)
3. model type (optional, default: turbo)
4. language (optional)
