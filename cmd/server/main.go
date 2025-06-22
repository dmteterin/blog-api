package main

import (
	"blog-api/internal/handler"
	"blog-api/internal/service"
	"blog-api/internal/storage"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
)

const port = ":8080"
const shouldSeed = true

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	postStorage := storage.NewInMemoryStorage(logger, shouldSeed)
	postService := service.NewPostService(postStorage)
	postHandler := handler.NewPostHandler(postService)

	srv := &http.Server{
		Addr:    port,
		Handler: handler.NewRouter(postHandler, logger),
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Info().Msgf("Server starting on port %v...", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("Server failed")
		}
	}()

	<-ctx.Done()
	logger.Info().Msg("Shutdown signal received, shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error().Err(err).Msg("Server shutdown failed")
	} else {
		logger.Info().Msg("Server gracefully stopped")
	}
}
