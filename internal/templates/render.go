package templates

import (
	"fmt"

	"mailmate/internal/models"
	"mailmate/internal/validator"

	"github.com/flosch/pongo2/v6"
)

// init registers the custom filters with pongo2.
func init() {
	// Register the "type" filter.
	// Usage: {{ Variable | type:"date" }}
	if err := pongo2.RegisterFilter("type", filterType); err != nil {
		panic(fmt.Errorf(`failed to register pongo2 filter %q: %w`, "type", err))
	}

	// Register the "int" filter.
	// Usage: {{ Variable | int }}
	if err := pongo2.RegisterFilter("int", filterInt); err != nil {
		panic(fmt.Errorf(`failed to register pongo2 filter %q: %w`, "int", err))
	}
}

// filterType implements the "type" filter which validates/formats values based on the argument.
// Currently supported types: "date", "filepath".
func filterType(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	val := in.String()
	typ := param.String()

	switch typ {
	case "date":
		// Check validity but return string as is for template
		if _, err := validator.ValidateDate(val); err != nil {
			return nil, &pongo2.Error{Sender: "filter:type", OrigError: err}
		}
		return pongo2.AsValue(val), nil
	case "filepath":
		if _, err := validator.ValidateFilepath(val); err != nil {
			return nil, &pongo2.Error{Sender: "filter:type", OrigError: err}
		}
		// For display purposes in the email body, we only want the filename, not the full path.
		return pongo2.AsValue(validator.GetFilename(val)), nil
	default:
		return nil, &pongo2.Error{
			Sender:    "filter:type",
			OrigError: fmt.Errorf("unknown type: %s", typ),
		}
	}
}

// filterInt implements the "int" filter which ensures the value is an integer.
func filterInt(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	s := in.String()
	i, err := validator.ValidateInt(s)
	if err != nil {
		return nil, &pongo2.Error{
			Sender:    "filter:int",
			OrigError: err,
		}
	}
	return pongo2.AsValue(i), nil
}

// passthroughFilter implements a pass-through filter that returns the input value unchanged.
// This is useful when a filter definition is needed for the TUI (to trigger specific form behaviors)
// but no transformation is required during the actual template rendering.
/* func passthroughFilter(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return pongo2.AsValue(in), nil
} */

// RenderTemplate renders the template at the given path using the provided variables.
func RenderTemplate(tmplPath string, variables map[string]string) (*models.RenderedTemplate, error) {
	// Parse the template file to separate frontmatter (subject) and body.
	parsed, err := ParseTemplateFile(tmplPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template file %q: %w", tmplPath, err)
	}

	// Convert map[string]string to pongo2.Context (map[string]interface{})
	ctx := pongo2.Context{}
	for k, v := range variables {
		ctx[k] = v
	}

	// 1. Render the Body
	// We use FromString because we have already read and stripped the frontmatter.
	bodyTpl, err := pongo2.FromString(parsed.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template body for %q: %w", tmplPath, err)
	}

	bodyOut, err := bodyTpl.Execute(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to render template body for %q: %w", tmplPath, err)
	}

	// 2. Render the Subject
	// The subject might also contain variables.
	subjectTpl, err := pongo2.FromString(parsed.Subject)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template subject for %q: %w", tmplPath, err)
	}

	subjectOut, err := subjectTpl.Execute(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to render template subject for %q: %w", tmplPath, err)
	}

	// 3. Render the recipients (To, Cc, Bcc)
	// These might also contain variables.
	var toOut, ccOut, bccOut string

	if parsed.To != "" {
		toTpl, err := pongo2.FromString(parsed.To)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template 'to' field for %q: %w", tmplPath, err)
		}
		toOut, err = toTpl.Execute(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to render template 'to' field for %q: %w", tmplPath, err)
		}
	}

	if parsed.Cc != "" {
		ccTpl, err := pongo2.FromString(parsed.Cc)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template 'cc' field for %q: %w", tmplPath, err)
		}
		ccOut, err = ccTpl.Execute(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to render template 'cc' field for %q: %w", tmplPath, err)
		}
	}

	if parsed.Bcc != "" {
		bccTpl, err := pongo2.FromString(parsed.Bcc)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template 'bcc' field for %q: %w", tmplPath, err)
		}
		bccOut, err = bccTpl.Execute(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to render template 'bcc' field for %q: %w", tmplPath, err)
		}
	}

	return &models.RenderedTemplate{
		Subject: subjectOut,
		HTML:    bodyOut,
		To:      toOut,
		Cc:      ccOut,
		Bcc:     bccOut,
	}, nil
}

//TODO : maybe I can build a package to generate filters for render in same object has form checker for tui, it would export the required function that iwll be imported by form.go and render.go
