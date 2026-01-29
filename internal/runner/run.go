package runner

import (
	"fmt"
	"path/filepath"

	"mailmate/internal/kv"
	"mailmate/internal/mailer"
	"mailmate/internal/models"
	"mailmate/internal/templates"
	"mailmate/internal/tui"
)

// displayRequiredVariables prints the list of required template variables
func displayRequiredVariables(vars []models.TemplateVariable) {
	fmt.Println("Template variables required:")
	for _, v := range vars {
		filterInfo := ""
		if len(v.Filters) > 0 {
			filterNames := make([]string, len(v.Filters))
			for i, f := range v.Filters {
				if f.Arg != "" {
					filterNames[i] = fmt.Sprintf("%s:%s", f.Name, f.Arg)
				} else {
					filterNames[i] = f.Name
				}
			}
			filterInfo = fmt.Sprintf(" (filters: %s)", filterNames[0])
			if len(filterNames) > 1 {
				for i := 1; i < len(filterNames); i++ {
					filterInfo = filterInfo[:len(filterInfo)-1] + ", " + filterNames[i] + ")"
				}
			}
		}
		fmt.Printf("  - %s%s\n", v.Name, filterInfo)
	}
	fmt.Println("\nUsage: --kv \"key1='value1';key2='value2'\"")
}

// Run executes the main application flow:
// 1. Scan templates
// 2. Select template
// 3. Parse variables
// 4. Collect user input
// 5. Render template
// 6. Send draft (via Outlook)
func Run(sender mailer.EmailSender, options models.Options) error {
	// 1. Scan templates
	// Assuming "templates" directory is in the current working directory
	tmpls, err := templates.ScanTemplates("templates")
	if err != nil {
		return fmt.Errorf("scanning templates: %w", err)
	}

	// 2. Select template
	var selected *models.TemplateRef
	if options.Template != nil {
		// --template flag was provided
		if *options.Template == "" {
			// --template flag provided but empty: list available templates and exit
			fmt.Println("Available templates:")
			for _, tmpl := range tmpls {
				fmt.Printf("  - %s\n", tmpl.Path)
			}
			fmt.Println("\nUsage: --template <path>")
			return nil
		}
		
		// CLI template selection: find matching template by path
		// Normalize paths for comparison (handles / vs \ on Windows)
		normalizedInput := filepath.Clean(*options.Template)
		found := false
		for i := range tmpls {
			if filepath.Clean(tmpls[i].Path) == normalizedInput {
				selected = &tmpls[i]
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("template not found: %s", *options.Template)
		}
	} else {
		// --template flag not provided: use TUI template selection
		selected, err = tui.SelectTemplate(tmpls)
		if err != nil {
			return fmt.Errorf("selecting template: %w", err)
		}
	}

	// 3. Parse variables
	vars, err := templates.ParseTemplate(selected.Path)
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}

	// 4. Collect user input
	var input *models.UserInput
	
	// Check if --kv flag was provided
	if options.KV != nil {
		// Flag was provided - check if it's empty
		if *options.KV == "" {
			// --kv flag provided but empty: show required variables and exit
			displayRequiredVariables(vars)
			return nil
		}
		
		// --kv flag provided with values: parse and use them
		kvValues, err := kv.Parse(*options.KV)
		if err != nil {
			return fmt.Errorf("parsing key-value pairs: %w", err)
		}

		// Validate values against template variables
		if err := kv.ValidateValues(kvValues, vars); err != nil {
			fmt.Printf("Error: %v\n\n", err)
			displayRequiredVariables(vars)
			return fmt.Errorf("validation failed")
		}

		input = &models.UserInput{
			Values: kvValues,
		}
	} else {
		// --kv flag not provided: use TUI to collect input
		input, err = tui.CollectUserInput(vars)
		if err != nil {
			return fmt.Errorf("collecting input: %w", err)
		}
	}

	// Handle attachments
	var attachments []string
	for _, v := range vars {
		for _, f := range v.Filters {
			if f.Name == "type" && f.Arg == "filepath" {
				fullPath, ok := input.Values[v.Name]
				if ok && fullPath != "" {
					absPath, err := filepath.Abs(fullPath)
					if err != nil {
						return fmt.Errorf("resolving absolute path for %s: %w", fullPath, err)
					}
					attachments = append(attachments, absPath)
				}
			}
		}
	}

	// 5. Render template
	rendered, err := templates.RenderTemplate(selected.Path, input.Values)
	if err != nil {
		return fmt.Errorf("rendering template: %w", err)
	}

	// TODO: T019 - Implement preview screen here
	// if !options.NoPreview { ... }

	// 6. Send draft
	draft := models.DraftEmail{
		To:          options.To,
		Cc:          options.Cc,
		Bcc:         options.Bcc,
		Subject:     rendered.Subject,
		HTMLBody:    rendered.HTML,
		Attachments: attachments,
	}

	if err := sender.Send(draft); err != nil {
		return fmt.Errorf("sending draft: %w", err)
	}

	return nil
}
