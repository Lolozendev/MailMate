package templates

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"mailmate/internal/models"
)

// ScanTemplates searches for .html files in the specified directory.
// It returns a list of found templates or an error if the directory is missing or empty.
func ScanTemplates(dir string) ([]models.TemplateRef, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("templates directory not found at %q: %w", dir, err)
		}
		return nil, fmt.Errorf("failed to read templates directory: %w", err)
	}

	var templates []models.TemplateRef
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".html") {
			templates = append(templates, models.TemplateRef{
				Name: entry.Name(),
				Path: filepath.Join(dir, entry.Name()),
			})
		}
	}

	if len(templates) == 0 {
		return nil, fmt.Errorf("no templates found in %q", dir)
	}

	return templates, nil
}
