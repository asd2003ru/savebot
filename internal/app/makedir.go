package app

import (
	"errors"
	"fmt"
	"os"
	"path"
)

const workDir = "Telegram"

func MakeUserDirs(homeDir string) (string, error) {
	if homeDir == "" {
		return "", errors.New("user home directory name is empty")
	}

	subDirPath := path.Join(homeDir, workDir)
	if err := os.MkdirAll(subDirPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", subDirPath, err)
	}

	return subDirPath, nil
}
