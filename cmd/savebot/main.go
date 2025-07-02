package main

import (
	"context"
	"savebot/internal/app"
	"savebot/internal/config"
	"io"
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		os.Stderr.WriteString("Failed to load config: " + err.Error() + "\n")
		os.Exit(1)
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Open log file for appending, create if not exists
	logFile, err := os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		os.Stderr.WriteString("Failed to open log file " + cfg.LogFile)
		os.Exit(1)
	}
	defer logFile.Close()

	zerolog.TimeFieldFormat = "2006-01-02T15:04:05.000Z"
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: multiWriter, NoColor: true, TimeFormat: time.RFC3339}).With().Logger()

	log.Info().Msg("Bot starting...")

	// Graceful shutdown handling
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	if err := app.Run(ctx, &log.Logger, cfg); err != nil {
		log.Error().Err(err).Msg("Failed to run bot")
		os.Exit(1)
	}
	<-ctx.Done()
	stop()
	log.Info().Msg("Bot shutting down")

}
