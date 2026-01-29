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

		// Apply filter-based validation using centralized validator
		if err := validator.ApplyFilters(value, v.Filters); err != nil {
			return fmt.Errorf("variable %s: %w", v.Name, err)
		}
	}

	return nil
}
