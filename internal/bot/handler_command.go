package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleCommand(msg *tgbotapi.Message) {
	if msg == nil || msg.Text == "" {
		return
	}

	switch msg.Command() {
	case "start":
		b.log.Info().Msgf("User %s (%d) started the bot", msg.From.UserName, msg.From.ID)
		if _, ok := b.users[msg.From.ID]; ok {
			b.api.Send(tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Welcome %s!", msg.From.UserName)))
		} else {
			b.api.Send(tgbotapi.NewMessage(msg.Chat.ID,
				fmt.Sprintf("Hello %s!\n You are not registered.\n Please contact the bot administrator and send you ID: %d.",
					msg.From.UserName, msg.From.ID)))
		}

	}
}
