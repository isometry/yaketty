package library

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

// embeddedFS will be set by the root package
// We use a variable here that will be initialized externally
var embeddedFS embed.FS

// SetEmbeddedFS allows the root package to set the embedded filesystem
func SetEmbeddedFS(fs embed.FS) {
	embeddedFS = fs
}

const (
	LibraryTypePersona  = "personas"
	LibraryTypeScenario = "scenarios"
)

// ReadFile attempts to read a file from the library with local override support.
// It first checks the local filesystem at libraryPath/filename.
// If the file doesn't exist locally, it falls back to the embedded filesystem.
// Returns the file contents and any error encountered.
func ReadFile(libraryPath, filename string) ([]byte, error) {
	// Try local filesystem first (for user overrides)
	localPath := filepath.Join(libraryPath, filename)
	if data, err := os.ReadFile(localPath); err == nil {
		slog.Debug("loaded file from local filesystem", slog.String("path", localPath))
		return data, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		// If error is not "file not found", return it
		return nil, fmt.Errorf("error reading local file %s: %w", localPath, err)
	}

	// Fall back to embedded filesystem
	embeddedPath := filepath.Join(libraryPath, filename)
	data, err := embeddedFS.ReadFile(embeddedPath)
	if err != nil {
		return nil, fmt.Errorf("file not found in local or embedded filesystem: %s", filename)
	}

	slog.Debug("loaded file from embedded filesystem", slog.String("path", embeddedPath))
	return data, nil
}

// FileExists checks if a file exists in either local filesystem or embedded filesystem.
// Returns true if the file exists in either location.
func FileExists(libraryPath, filename string) bool {
	// Check local filesystem first
	localPath := filepath.Join(libraryPath, filename)
	if _, err := os.Stat(localPath); err == nil {
		return true
	}

	// Check embedded filesystem
	embeddedPath := filepath.Join(libraryPath, filename)
	if _, err := embeddedFS.Open(embeddedPath); err == nil {
		return true
	}

	return false
}

// ListPersonas returns a list of all embedded persona filenames (without .yaml extension).
func ListPersonas() ([]string, error) {
	return listFiles(LibraryTypePersona)
}

// ListScenarios returns a list of all embedded scenario filenames (without .yaml extension).
func ListScenarios() ([]string, error) {
	return listFiles(LibraryTypeScenario)
}

// listFiles returns a list of all .yaml files in the specified embedded directory.
func listFiles(libraryType string) ([]string, error) {
	entries, err := fs.ReadDir(embeddedFS, libraryType)
	if err != nil {
		return nil, fmt.Errorf("error reading embedded directory %s: %w", libraryType, err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".yaml") {
			// Remove .yaml extension for cleaner output
			name := strings.TrimSuffix(entry.Name(), ".yaml")
			files = append(files, name)
		}
	}

	return files, nil
}

// IsDirectPath returns true if the path should be treated as a direct file path.
// Rules:
// - Is absolute path → direct path
// - Contains path separator → direct path
// - Has .yaml extension AND file exists → direct path
// - Otherwise → library reference
func IsDirectPath(path string) bool {
	// Absolute paths or paths with separators are always direct
	if filepath.IsAbs(path) || strings.Contains(path, string(filepath.Separator)) {
		return true
	}

	// If has .yaml extension, check if file exists locally
	// If it exists, it's a direct path. If not, it might be a library reference (backward compat)
	if strings.HasSuffix(path, ".yaml") {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}

	return false
}

// ReadFileOrPath reads from a full file path.
// The path can be:
// - A direct filesystem path (passed through to os.ReadFile)
// - A library-constructed path like "personas/biden.yaml" (tries local then embedded)
func ReadFileOrPath(fullPath string) ([]byte, error) {
	// First try to read as a direct path from filesystem
	if data, err := os.ReadFile(fullPath); err == nil {
		slog.Debug("loaded file from direct path", slog.String("path", fullPath))
		return data, nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("error reading file %s: %w", fullPath, err)
	}

	// If local file doesn't exist, try embedded filesystem
	data, err := embeddedFS.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("file not found in local or embedded filesystem: %s", fullPath)
	}

	slog.Debug("loaded file from embedded filesystem", slog.String("path", fullPath))
	return data, nil
}

// GetEmbeddedFile retrieves embedded file content by name (without checking local filesystem).
// Automatically adds .yaml extension if not present.
func GetEmbeddedFile(libraryType, name string) ([]byte, error) {
	filename := name
	if !strings.HasSuffix(filename, ".yaml") {
		filename = filename + ".yaml"
	}

	embeddedPath := filepath.Join(libraryType, filename)
	data, err := embeddedFS.ReadFile(embeddedPath)
	if err != nil {
		return nil, fmt.Errorf("embedded file not found: %s", name)
	}
	return data, nil
}
