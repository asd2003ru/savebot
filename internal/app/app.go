package app

import (
	"context"
	"fmt"
	"savebot/internal/bot"
	"savebot/internal/config"
	"savebot/internal/logger"
)

func Run(ctx context.Context, log logger.ILogger, config *config.Config) error {

	// Create work directories for each user
	log.Info("Create, if not exist, work directories for %d users", len(config.Users))
	for chatID, home := range config.Users {

		if wd, err := MakeUserDirs(home); err != nil {
			return fmt.Errorf("failed to create work directory for user %d: %w", chatID, err)
		} else {
			log.Info("Work directory '%s' created for user '%d'", wd, chatID)
			config.Users[chatID] = wd // Update home path in config
		}
	}

	botInstance, err := bot.NewBot(config.Users, config.BotToken, log)
	if err != nil {
		return fmt.Errorf("failed to create bot instance: %w", err)
	}

	if err := botInstance.Start(ctx); err != nil {
		return fmt.Errorf("failed to start bot: %w", err)
	}

	return nil
}
