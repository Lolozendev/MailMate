package tui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/huh"

	"mailmate/internal/app"
)

// CollectUserInput prompts the user for any variables defined in the selected template.
func CollectUserInput(variables []app.TemplateVariable) (*app.UserInput, error) {
	// If there are no variables, return empty input immediately
	if len(variables) == 0 {
		return &app.UserInput{
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

	return &app.UserInput{
		Values: finalValues,
	}, nil
}

// createValidator returns a validation function based on the provided filters.
func createValidator(filters []app.TemplateFilter) func(string) error {
	return func(str string) error {
		// Required check (variables are implicitly required for now unless we add an optional filter later)
		if strings.TrimSpace(str) == "" {
			return fmt.Errorf("value is required")
		}

		for _, f := range filters {
			switch f.Name {
			case "int":
				if _, err := strconv.Atoi(str); err != nil {
					return fmt.Errorf("must be an integer")
				}
			case "type":
				if f.Arg == "date" {
					// validate YYYY-MM-DD
					// TODO : make this a switch case for easier adding of future types
					if _, err := time.Parse("02-01-2006", str); err != nil {
						return fmt.Errorf("must be a date (DD-MM-YYYY)")
					}
				}
			}
		}
		return nil
	}
}

// getHint returns a placeholder string based on filters
// TODO : make this a switch case for easier adding of future types
func getHint(filters []app.TemplateFilter) string {
	for _, f := range filters {
		if f.Name == "type" && f.Arg == "date" {
			return "DD-MM-YYYY"
		}
	}
	return ""
}
