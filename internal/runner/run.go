package runner

import (
	"fmt"
	"path/filepath"

	"mailmate/internal/mailer"
	"mailmate/internal/models"
	"mailmate/internal/templates"
	"mailmate/internal/tui"
)

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
	if options.Template != "" {
		// CLI template selection: find matching template by path
		found := false
		for i := range tmpls {
			if tmpls[i].Path == options.Template {
				selected = &tmpls[i]
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("template not found: %s", options.Template)
		}
	} else {
		// TUI template selection
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
	input, err := tui.CollectUserInput(vars)
	if err != nil {
		return fmt.Errorf("collecting input: %w", err)
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
