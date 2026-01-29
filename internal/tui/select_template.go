package tui

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/huh"

	"mailmate/internal/app"
)

// SelectTemplate prompts the user to select a template from the provided list.
// It returns the selected TemplateRef or an error if selection fails/is cancelled.
func SelectTemplate(templates []app.TemplateRef) (*app.TemplateRef, error) {
	if len(templates) == 0 {
		return nil, errors.New("no templates available to select")
	}

	// Create options for huh.Select
	options := make([]huh.Option[string], len(templates))
	// Map to look up the full struct by Name (since Select returns the value, which we'll use as the Name for uniqueness)
	// Alternatively, we can use the Path as the value. Let's use Path as value to be safe if names collide (though scanner should handle that).
	// But to return the original struct, a map is easy.
	templateMap := make(map[string]app.TemplateRef)

	for i, tmpl := range templates {
		options[i] = huh.NewOption(tmpl.Name, tmpl.Path)
		templateMap[tmpl.Path] = tmpl
	}

	var selectedPath string

	// Create the form with a Select field
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a template").
				Options(options...).
				Value(&selectedPath),
		),
	)

	err := form.Run()
	if err != nil {
		return nil, fmt.Errorf("template selection failed: %w", err)
	}

	selected, ok := templateMap[selectedPath]
	if !ok {
		return nil, errors.New("selected template not found in map")
	}

	return &selected, nil
}
