package bot

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"savebot/internal/logger"
	"time"

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

	log := b.log.WithFields(logger.Fields{"chat_id": msg.Chat.ID, "from": msg.From.UserName, "message_id": msg.MessageID})

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
		log.Warn("Unknown or unsupported message type")
		return fmt.Errorf("unknown or unsupported message type")
	}

	// Save media to disk
	if contentType != ContentTypeText {
		p, err := b.saveFile(fileID, b.users[msg.From.ID], filename, contentType)
		if err != nil {
			log.Error(err, "Failed to save %s", filename)
			return err
		}
		log.Info("File saved %s", p)
	}
	// Save text
	if contentType == ContentTypeText || contentType == ContentTypeMixed {
		p, err := b.saveTextMsg(msg, filename)
		if err != nil {
			log.Error(err, "Failed to save text message")
			return err
		}
		log.Info("Text saved to file %s", p)
	}

	return nil
}

// Save attachment to disk
// Returns full path
func (b *Bot) saveFile(fileID string, filepath, filename string, contentType ContentType) (string, error) {
	fileConfig := tgbotapi.FileConfig{FileID: fileID}
	file, err := b.api.GetFile(fileConfig)
	if err != nil {
		b.log.Error(err, "Failed to get file %s", fileID)
		return "", err
	}
	fileURL := file.Link(b.api.Token)
	u, err := b.api.GetFileDirectURL(fileID)
	if err != nil {
		b.log.Error(err, "Failed to get file %s", fileID)
		return "", err
	}
	fileURL = u

	// Download
	resp, err := http.Get(fileURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		b.log.Error(err, "Failed to download file %s status code %d", fileURL, resp.StatusCode)
		return "", fmt.Errorf("failed to download file %s status code %d", fileURL, resp.StatusCode)
	}
	defer resp.Body.Close()

	// Filename may be empty
	if filename == "" {
		filename = time.Now().Format("20060102_150405")
	}

	// Directory by type
	var dir string
	switch contentType {
	case ContentTypeImage:
		dir = path.Join(filepath, "Images")
		filename = filename + ".jpg"
	case ContentTypeDocument:
		dir = path.Join(filepath, "Documents")
	case ContentTypeAudio:
		dir = path.Join(filepath, "Audio")
		filename = filename + ".ogg"
	case ContentTypeVideo:
		dir = path.Join(filepath, "Video")
		filename = filename + ".mp4"
	case ContentTypeVoice:
		dir = path.Join(filepath, "Voice")
		filename = filename + ".ogg"
	case ContentTypeMixed:
		dir = path.Join(filepath, "Text", "media")

	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		b.log.Error(err, "Failed to create directory %s", dir)
		return "", err
	}

	fullpath := path.Join(dir, filename)
	outFile, err := os.Create(fullpath)
	if err != nil {
		b.log.Error(err, "Failed to create file %s", fullpath)
		return "", err
	}
	defer outFile.Close()

	if _, err := io.Copy(outFile, resp.Body); err != nil {
		b.log.Error(err, "Failed save data to file %s", fullpath)
		return "", err
	}

	return fullpath, nil
}
