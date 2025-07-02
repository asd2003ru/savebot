package app

import (
	"context"
	"fmt"
	"savebot/internal/bot"
	"savebot/internal/config"

	"github.com/rs/zerolog"
)

func Run(ctx context.Context, log *zerolog.Logger, config *config.Config) error {

	// Create work directories for each user
	log.Info().Msgf("Creating work directories for %d users", len(config.Users))
	for chatID, home := range config.Users {

		if wd, err := MakeUserDirs(home); err != nil {
			return fmt.Errorf("failed to create work directory for user %d: %w", chatID, err)
		} else {
			log.Info().Msgf("Work directory '%s' created for user '%d'", wd, chatID)
			config.Users[chatID] = wd // Update home path in config
		}
	}

	botInstance, err := bot.NewBot(config.Users, config.BotToken, log)
	if err != nil {
		return fmt.Errorf("failed to create bot instance: %w", err)
	}

	log.Info().Msg("Starting bot...")
	if err := botInstance.Start(ctx); err != nil {
		return fmt.Errorf("failed to start bot: %w", err)
	}

	return nil
}
