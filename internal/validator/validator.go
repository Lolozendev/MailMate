package validator

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// ValidateInt checks if the value represents a valid integer.
// It returns the integer value if valid, or an error if not.
func ValidateInt(value string) (int, error) {
	i, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("value %q is not a valid integer", value)
	}
	return i, nil
}

// ValidateDate checks if the value matches the expected date format (DD-MM-YYYY).
// It returns the time.Time object if valid, or an error if not.
func ValidateDate(value string) (time.Time, error) {
	t, err := time.Parse("02-01-2006", value)
	if err != nil {
		return time.Time{}, fmt.Errorf("value %q is not a valid date (expected DD-MM-YYYY)", value)
	}
	return t, nil
}

// ValidateFilepath checks if the value is a non-empty filepath.
// Note: It does NOT check if the file exists on disk, as templates might be rendered
// without the file strictly needing to exist at that moment (e.g. for preview),
// or it might be a relative path resolved later.
// However, if strict existence check is needed, use ValidateFileExists.
func ValidateFilepath(value string) (string, error) {
	if strings.TrimSpace(value) == "" {
		return "", fmt.Errorf("filepath cannot be empty")
	}
	// Basic sanity check, maybe ensure it's not just whitespace
	return value, nil
}

// ValidateFileExists checks if the file actually exists on the filesystem.
func ValidateFileExists(value string) error {
	if _, err := os.Stat(value); os.IsNotExist(err) {
		return fmt.Errorf("file %q does not exist", value)
	}
	return nil
}

// GetFilename returns the base name of a filepath.
// Useful for displaying just the filename in templates instead of full path.
func GetFilename(value string) string {
	return filepath.Base(value)
}
