package bot

import (
	"fmt"
	"savebot/internal/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleCommand(msg *tgbotapi.Message) {
	if msg == nil || msg.Text == "" {
		return
	}

	log := b.log.WithFields(logger.Fields{"chat_id": msg.Chat.ID, "from": msg.From.UserName, "command": msg.Command()})

	switch msg.Command() {
	case "start":
		if _, ok := b.users[msg.From.ID]; ok {
			log.Debug("Regidstred user send start command")
			b.api.Send(tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Welcome %s!", msg.From.UserName)))
		} else {
			log.Info("User want to register")
			b.api.Send(tgbotapi.NewMessage(msg.Chat.ID,
				fmt.Sprintf("Hello %s!\n You are not registered.\n Please contact the bot administrator and send you ID: %d.",
					msg.From.UserName, msg.From.ID)))
		}

	}
}
