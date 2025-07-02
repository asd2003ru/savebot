package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleUpdate(update tgbotapi.Update) {

	if update.CallbackQuery != nil {
		b.handleCallback(update.CallbackQuery)
		return
	}

	if update.Message == nil { // ignore non-message updates
		return
	}

	msg := update.Message

	// This is command
	if msg.IsCommand() {
		b.handleCommand(msg)
		return
	}

	if _, ok := b.users[msg.From.ID]; !ok {
		//b.api.Send(tgbotapi.NewMessage(msg.Chat.ID, "You are not registered. Please contact the bot administrator and send you ID: "+string(msg.From.ID)))
		return
	}

	// No chat groups
	if !msg.Chat.IsPrivate() {
		return
	}

	if err := b.saveMessage(msg); err != nil {
		b.log.Error().Err(err).Msgf("Failed to save message")
		b.api.Send(tgbotapi.NewMessage(msg.Chat.ID, "🚨 Failed to save message"))
		return
	}

	b.api.Send(tgbotapi.NewMessage(msg.Chat.ID, "✅ Message saved!"))

}
