package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"

	"mailmate/internal/models"
	"mailmate/internal/validator"
)

// CollectUserInput prompts the user for any variables defined in the selected template.
func CollectUserInput(variables []models.TemplateVariable) (*models.UserInput, error) {
	// If there are no variables, return empty input immediately
	if len(variables) == 0 {
		return &models.UserInput{
			Values: make(map[string]string),
		}, nil
	}

	// Create dynamic fields for variables
	variableValues := make(map[string]*string)
	var variableFields []huh.Field

	for _, v := range variables {
		valPtr := new(string)
		variableValues[v.Name] = valPtr

		input := huh.NewInput().
			Title(v.Name).
			Value(valPtr)

		// Apply basic validation based on filters
		input.Validate(createValidator(v.Filters))

		// Add placeholder hint based on filters
		if hint := getHint(v.Filters); hint != "" {
			input.Placeholder(hint)
		}

		variableFields = append(variableFields, input)
	}

	variableGroup := huh.NewGroup(variableFields...).Title("Template Variables")

	form := huh.NewForm(variableGroup)
	err := form.Run()
	if err != nil {
		return nil, fmt.Errorf("form cancelled/error: %w", err)
	}

	// Collect values
	finalValues := make(map[string]string)
	for name, ptr := range variableValues {
		finalValues[name] = *ptr
	}

	return &models.UserInput{
		Values: finalValues,
	}, nil
}

// createValidator returns a validation function based on the provided filters.
func createValidator(filters []models.TemplateFilter) func(string) error {
	return func(str string) error {
		// Required check (variables are implicitly required for now unless we add an optional filter later)
		if strings.TrimSpace(str) == "" {
			return fmt.Errorf("value is required")
		}

		for _, f := range filters {
			switch f.Name {
			case "int":
				if _, err := validator.ValidateInt(str); err != nil {
					return fmt.Errorf("must be an integer")
				}
			case "type":
				switch f.Arg {
				case "date":
					// validate YYYY-MM-DD
					if _, err := validator.ValidateDate(str); err != nil {
						return fmt.Errorf("must be a date (DD-MM-YYYY)")
					}
				case "filepath":
					if err := validator.ValidateFileExists(str); err != nil {
						return fmt.Errorf("file does not exist")
					}
				}
			}
		}
		return nil
	}
}

// getHint returns a placeholder string based on filters
// TODO : make this a switch case for easier adding of future types
func getHint(filters []models.TemplateFilter) string {
	for _, f := range filters {
		if f.Name == "type" {
			switch f.Arg {
			case "date":
				return "DD-MM-YYYY"
			case "filepath":
				return "/path/to/file"
			}
		}
	}
	return ""
}
