package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/imvalerio/transcribeit/internal/utils"
)

type TranscribeHandler struct {
	log *log.Logger
}

func NewTranscribeHandler(logger *log.Logger) *TranscribeHandler {
	return &TranscribeHandler{
		log: logger,
	}
}

func (th *TranscribeHandler) UploadAudio(w http.ResponseWriter, r *http.Request) {
	// r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		th.log.Printf("Error Retrieving the File: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid file"})
		return
	}
	defer file.Close()
	th.log.Printf("Uploaded File: %+v\n", handler.Filename)
	th.log.Printf("File Size: %+v\n", handler.Size)
	th.log.Printf("MIME Header: %+v\n", handler.Header)

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	//
	uuid := uuid.New().String()
	lindex := strings.LastIndex(handler.Filename, ".")
	ext := handler.Filename[lindex+1:]
	tempFile, err := os.CreateTemp("temp-audios", fmt.Sprintf("%s-*."+ext, uuid))
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		th.log.Println("Error creating temporary file:", err)
		return
	}
	defer tempFile.Close()

	audioID := strings.TrimPrefix(tempFile.Name(), "temp-audios\\")

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		th.log.Printf("Error reading file: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid file"})
		return
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	absPath, err := filepath.Abs(tempFile.Name())
	if err != nil {
		th.log.Printf("Error getting absolute path: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}

	err = utils.PushToQueue("upload-audio", absPath)
	if err != nil {
		th.log.Printf("Error pushing to queue: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"audio_id": audioID})
}

func (th *TranscribeHandler) GetTranscription(w http.ResponseWriter, r *http.Request) {
	audioID := chi.URLParam(r, "id")
	if audioID == "" {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Missing audio_id"})
		return
	}

	if !strings.Contains(audioID, ".") {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "Invalid audio_id"})
		return
	}

	fileName := strings.Split(audioID, ".")[0]
	file, err := os.Open("transcriptions/" + fileName + ".txt")
	if err != nil {
		th.log.Printf("Error opening file: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "Internal server error"})
		return
	}

	defer file.Close()

	buf := make([]byte, 1024)
	var builder strings.Builder

	for {
		n, err := file.Read(buf)
		if n > 0 {
			chunk := buf[:n]
			builder.Write(chunk)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	text := builder.String()

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"transcription": text})
}
