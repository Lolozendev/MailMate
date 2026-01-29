package kv

import (
	"os"
	"testing"

	"mailmate/internal/models"
)

func TestValidateValues(t *testing.T) {
	// Create a temporary file for filepath testing
	tmpFile, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	tests := []struct {
		name      string
		kvValues  map[string]string
		variables []models.TemplateVariable
		wantErr   bool
		errMsg    string
	}{
		{
			name: "valid simple values",
			kvValues: map[string]string{
				"Name":  "John",
				"Count": "42",
			},
			variables: []models.TemplateVariable{
				{Name: "Name", Filters: []models.TemplateFilter{}},
				{Name: "Count", Filters: []models.TemplateFilter{{Name: "int"}}},
			},
			wantErr: false,
		},
		{
			name: "valid date",
			kvValues: map[string]string{
				"EventDate": "25-01-2026",
			},
			variables: []models.TemplateVariable{
				{Name: "EventDate", Filters: []models.TemplateFilter{{Name: "type", Arg: "date"}}},
			},
			wantErr: false,
		},
		{
			name: "invalid int",
			kvValues: map[string]string{
				"Count": "not-a-number",
			},
			variables: []models.TemplateVariable{
				{Name: "Count", Filters: []models.TemplateFilter{{Name: "int"}}},
			},
			wantErr: true,
			errMsg:  "variable Count:",
		},
		{
			name: "invalid date format",
			kvValues: map[string]string{
				"EventDate": "2026-01-25",
			},
			variables: []models.TemplateVariable{
				{Name: "EventDate", Filters: []models.TemplateFilter{{Name: "type", Arg: "date"}}},
			},
			wantErr: true,
			errMsg:  "variable EventDate:",
		},
		{
			name: "missing required variable",
			kvValues: map[string]string{
				"Name": "John",
			},
			variables: []models.TemplateVariable{
				{Name: "Name", Filters: []models.TemplateFilter{}},
				{Name: "Email", Filters: []models.TemplateFilter{}},
			},
			wantErr: true,
			errMsg:  "variable Email is required",
		},
		{
			name: "unknown variable",
			kvValues: map[string]string{
				"Name":    "John",
				"Unknown": "value",
			},
			variables: []models.TemplateVariable{
				{Name: "Name", Filters: []models.TemplateFilter{}},
			},
			wantErr: true,
			errMsg:  "unknown variable: Unknown",
		},
		{
			name: "valid filepath",
			kvValues: map[string]string{
				"Attachment": tmpFile.Name(),
			},
			variables: []models.TemplateVariable{
				{Name: "Attachment", Filters: []models.TemplateFilter{{Name: "type", Arg: "filepath"}}},
			},
			wantErr: false,
		},
		{
			name: "invalid filepath - file does not exist",
			kvValues: map[string]string{
				"Attachment": "/nonexistent/file.txt",
			},
			variables: []models.TemplateVariable{
				{Name: "Attachment", Filters: []models.TemplateFilter{{Name: "type", Arg: "filepath"}}},
			},
			wantErr: true,
			errMsg:  "variable Attachment:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateValues(tt.kvValues, tt.variables)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && tt.errMsg != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateValues() error = %v, want error containing %q", err, tt.errMsg)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
