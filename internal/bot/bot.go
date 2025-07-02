package bot

import (
	"context"
	"savebot/internal/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api   *tgbotapi.BotAPI
	log   logger.ILogger
	users map[int64]string // Map of user IDs to usernames
}

// NewBot creates a new instance of the bot with the provided configuration, database, and file server
func NewBot(users map[int64]string, token string, log logger.ILogger) (*Bot, error) {
	botAPI, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	botAPI.Debug = false
	return &Bot{
		log:   log,
		api:   botAPI,
		users: users,
	}, nil
}

// Start initializes the bot and starts listening for updates
func (b *Bot) Start(ctx context.Context) error {

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.api.GetUpdatesChan(u)

	b.log.Info("Bot @%s started, wait for messages.", b.api.Self.UserName)

	for {
		select {
		case <-ctx.Done():
			return nil
		case update := <-updates:
			if update.Message == nil {
				continue
			}
			go b.handleUpdate(update)
		}

	}
}
