package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/imvalerio/transcribeit/internal/app"
)

func SetupRoutes(a *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", app.HealthCheck)
	r.Post("/transcribe", a.TranscribeHandler.UploadAudio)
	r.Get("/transcribe/{id}", a.TranscribeHandler.GetTranscription)

	return r
}
