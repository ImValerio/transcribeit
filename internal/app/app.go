package app

import (
	"log"
	"net/http"
	"os"

	"github.com/imvalerio/transcribeit/internal/api"
	"github.com/imvalerio/transcribeit/internal/utils"
)

type Application struct {
	Log               *log.Logger
	TranscribeHandler *api.TranscribeHandler
}

func NewApplication() (*Application, error) {
	dirs := []string{"temp-audios", "transcriptions"}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			log.Fatalf("Error creating directory %s: %v", dir, err)
		}
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	return &Application{
		Log:               logger,
		TranscribeHandler: api.NewTranscribeHandler(logger),
	}, nil
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"health": "ok"})
}
