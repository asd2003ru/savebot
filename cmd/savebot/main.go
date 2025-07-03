package main

import (
	"context"
	"os"
	"os/signal"
	"savebot/internal/app"
	"savebot/internal/config"
	"savebot/internal/logger"
)

var (
	version = "dev" // default value; will be overwritten during build
	commit  = "none"
	date    = "unknown"
)

func main() {

	log := logger.NewLogger(logger.JSONType, logger.InfoLevel)

	log.Info("Bot started (version %s, commit %s, date %s)", version, commit, date)

	cfg, err := config.NewConfig()
	if err != nil {
		log.Error(err, "Failed to read config")
		os.Exit(1)
	}

	// Graceful shutdown handling
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	if err := app.Run(ctx, log, cfg); err != nil {
		log.Error(err, "Failed to run bot")
		stop()
		os.Exit(1)
	}
	<-ctx.Done()
	stop()
	log.Info("Bot shutting down")

}
