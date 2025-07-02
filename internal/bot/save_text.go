package bot

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const textFileDir = "Text"

func (b *Bot) saveTextMsg(msg *tgbotapi.Message, mediaFile string) error {

	if msg.Caption == "" && msg.Text == "" {
		return nil
	}

	dir := path.Join(b.users[msg.From.ID], textFileDir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	fullpath := path.Join(dir, time.Now().Format("20060102_150405")+".md")
	f, err := os.Create(fullpath)
	if err != nil {
		return err
	}
	defer f.Close()

	// Construct message for write

	var content strings.Builder

	// HEAD
	from := "From: Unknown"
	if msg.ForwardFrom != nil {
		from = "Forwarded from: @" + msg.ForwardFrom.UserName
	} else if msg.From != nil {
		from = "From: @" + msg.From.UserName
	}
	content.WriteString("# Message from " + from + "\n")
	content.WriteString("Date: " + msg.Time().Format("2006-01-02 15:04:05") + "\n\n")

	// BODY
	if mediaFile != "" {
		mediaFile = path.Join("media", filepath.Base(mediaFile))
		content.WriteString(fmt.Sprintf("![Attachment](%s)\n", mediaFile))
	}

	if msg.Caption != "" {
		content.WriteString(msg.Caption + "\n\n")
	} else if msg.Text != "" {
		content.WriteString(msg.Text + "\n\n")
	} else {
		content.WriteString("Empty message\n\n")
	}

	_, err = f.WriteString(content.String())
	return err

}
