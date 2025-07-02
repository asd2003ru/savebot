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
		b.log.Debug("Unregistered user (chat-id '%d', username '%s'), message rejected", msg.Chat.ID, msg.From.UserName)
		return
	}

	// No chat groups
	if !msg.Chat.IsPrivate() {
		b.log.Warn("Message from group (chat-id '%d', username '%s'), message rejected. Only private chats are supported.", msg.Chat.ID, msg.From.UserName)
		return
	}

	if err := b.saveMessage(msg); err != nil {
		//b.log.Error(err, "Failed to save message")
		b.api.Send(tgbotapi.NewMessage(msg.Chat.ID, "🚨 Failed to save message"))
		return
	}

	b.api.Send(tgbotapi.NewMessage(msg.Chat.ID, "✅ Message saved!"))

}
