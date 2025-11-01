package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/cors"
	"github.com/imvalerio/transcribeit/internal/app"
	"github.com/imvalerio/transcribeit/internal/consumer"
	"github.com/imvalerio/transcribeit/internal/routes"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	app, error := app.NewApplication()
	if error != nil {
		panic(error)
	}

	r := routes.SetupRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	brokerAddr := os.Getenv("BROKER_ADDRESS")
	if brokerAddr == "" {
		brokerAddr = "localhost:9092"
	}

	env := os.Getenv("APP_ENV") // e.g. "dev", "prod"
	if env == "" {
		env = "dev"
	}
	var handler http.Handler = r

	if env == "dev" {
		c := cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"*"},
			AllowCredentials: true,
		})
		handler = c.Handler(r)
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: handler,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := consumer.NewConsumer("upload-audio", brokerAddr, app.Log).Start(ctx); err != nil {
			app.Log.Printf("Consumer upload-audio stopped: %v", err)
		}
	}()
	go func() {
		if err := consumer.NewConsumer("completed-transcriptions", brokerAddr, app.Log).Start(ctx); err != nil {
			app.Log.Printf("Consumer upload-audio stopped: %v", err)
		}
	}()

	go func() {
		<-sigCh
		app.Log.Println("Shutting down...")
		cancel()
		server.Shutdown(context.Background())
	}()

	app.Log.Printf("Starting server on port %s", server.Addr)
	server.ListenAndServe()

}
