package bot

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ContentType int

const (
	ContentTypeUnknown ContentType = iota
	ContentTypeText
	ContentTypeImage
	ContentTypeDocument
	ContentTypeAudio
	ContentTypeVideo
	ContentTypeVoice
	ContentTypeMixed
)

func (b *Bot) saveMessage(msg *tgbotapi.Message) error {

	var fileID string
	var filename string
	var contentType ContentType
	if msg.Caption != "" || msg.Text != "" {
		contentType = ContentTypeText
	}

	dt := msg.Time().Format("20060102-150405") + fmt.Sprintf("-%d", msg.MessageID)

	if len(msg.Photo) > 0 {

		fileID = msg.Photo[len(msg.Photo)-1].FileID
		filename = dt + ".jpg"

		if contentType == ContentTypeText {
			contentType = ContentTypeMixed
		} else {
			contentType = ContentTypeImage
		}

	} else if msg.Document != nil {
		fileID = msg.Document.FileID
		filename = msg.Document.FileName
		if contentType == ContentTypeText {
			contentType = ContentTypeMixed
		} else {
			contentType = ContentTypeDocument
		}
	} else if msg.Audio != nil {
		fileID = msg.Audio.FileID
		filename = msg.Audio.FileName
		if contentType == ContentTypeText {
			contentType = ContentTypeMixed
		} else {
			contentType = ContentTypeAudio
		}
	} else if msg.Video != nil {
		fileID = msg.Video.FileID
		filename = msg.Video.FileName
		if contentType == ContentTypeText {
			contentType = ContentTypeMixed
		} else {
			contentType = ContentTypeVideo
		}
	} else if msg.Voice != nil {
		fileID = msg.Voice.FileID
		filename = dt + ".ogg"
		if contentType == ContentTypeText {
			contentType = ContentTypeMixed
		} else {
			contentType = ContentTypeVoice
		}
	}

	if contentType == ContentTypeUnknown {

		return fmt.Errorf("unknown content type")
	}

	// Save media to disk
	if contentType != ContentTypeText {
		if err := b.saveFile(fileID, b.users[msg.From.ID], filename, contentType); err != nil {
			b.log.Error().Err(err).Msgf("Failed to save image")
			return err
		}
		b.log.Info().Str("user", msg.From.UserName).Int64("chat_id", msg.Chat.ID).Msgf("Saved %s", filename)
	}
	// Save text
	if contentType == ContentTypeText || contentType == ContentTypeMixed {
		if err := b.saveTextMsg(msg, filename); err != nil {
			b.log.Error().Err(err).Msgf("Failed to save text")
			return err
		}
		b.log.Info().Str("user", msg.From.UserName).Int64("chat_id", msg.Chat.ID).Msgf("Saved %s", filename)
	}

	return nil
}

// Save attachment to disk
func (b *Bot) saveFile(fileID string, filepath, filename string, contentType ContentType) error {
	fileConfig := tgbotapi.FileConfig{FileID: fileID}
	file, _ := b.api.GetFile(fileConfig)
	fileURL := file.Link(b.api.Token)

	// Download
	resp, err := http.Get(fileURL)
	if err != nil {
		b.log.Error().Err(err).Msgf("Failed to download %s", fileURL)
		return err
	}
	defer resp.Body.Close()

	// Directory by type
	var dir string
	switch contentType {
	case ContentTypeImage:
		dir = path.Join(filepath, "Images")
	case ContentTypeDocument:
		dir = path.Join(filepath, "Documents")
	case ContentTypeAudio:
		dir = path.Join(filepath, "Audio")
	case ContentTypeVideo:
		dir = path.Join(filepath, "Video")
	case ContentTypeVoice:
		dir = path.Join(filepath, "Voice")
	case ContentTypeMixed:
		dir = path.Join(filepath, "Text", "media")
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		b.log.Error().Err(err).Msgf("Failed to create directory")
		return err
	}

	fullpath := path.Join(dir, filename)
	outFile, err := os.Create(fullpath)
	if err != nil {
		b.log.Error().Err(err).Msgf("Failed to create file %s", fullpath)
		return err
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, resp.Body); err != nil {
		b.log.Error().Err(err).Msgf("Failed to save data")
		return err
	}

	return nil
}
