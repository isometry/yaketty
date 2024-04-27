package utils

import (
	"log/slog"
	"os"
	"strings"
)

func IsFilePath(filePath string) bool {
	if strings.ContainsAny(filePath, " \t\n") {
		// only accept paths without whitespace
		return false
	}
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		slog.Debug("path does not exist", slog.String("path", filePath))
		return false
	}
	slog.Debug("path exists", slog.String("path", filePath), slog.Bool("isDir", fileInfo.IsDir()))
	return !fileInfo.IsDir()
}
