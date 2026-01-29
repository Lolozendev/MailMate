package kv

import (
	"fmt"
	"strings"

	"mailmate/internal/models"
	"mailmate/internal/validator"
)

// ValidateValues validates key-value pairs against template variables.
// It checks that all required variables are present and that values match their expected types.
func ValidateValues(kvValues map[string]string, variables []models.TemplateVariable) error {
	// Create a map of variable names for quick lookup
	varMap := make(map[string]models.TemplateVariable)
	for _, v := range variables {
		varMap[v.Name] = v
	}

	// Check that all provided keys exist in the template
	for key := range kvValues {
		if _, exists := varMap[key]; !exists {
			return fmt.Errorf("unknown variable: %s", key)
		}
	}

	// Validate each variable
	for _, v := range variables {
		value, exists := kvValues[v.Name]
		
		// Check if required variable is missing
		if !exists || strings.TrimSpace(value) == "" {
			return fmt.Errorf("variable %s is required", v.Name)
		}

		// Apply filter-based validation
		if err := validateWithFilters(value, v.Filters, v.Name); err != nil {
			return err
		}
	}

	return nil
}

// validateWithFilters applies validation based on template filters
func validateWithFilters(value string, filters []models.TemplateFilter, varName string) error {
	for _, f := range filters {
		switch f.Name {
		case "int":
			if _, err := validator.ValidateInt(value); err != nil {
				return fmt.Errorf("variable %s: %w", varName, err)
			}
		case "type":
			switch f.Arg {
			case "date":
				if _, err := validator.ValidateDate(value); err != nil {
					return fmt.Errorf("variable %s: %w", varName, err)
				}
			case "filepath":
				if err := validator.ValidateFileExists(value); err != nil {
					return fmt.Errorf("variable %s: %w", varName, err)
				}
			}
		}
	}
	return nil
}
