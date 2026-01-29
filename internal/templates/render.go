package templates

import (
	"fmt"
	"strconv"
	"strings"

	"mailmate/internal/app"

	"github.com/flosch/pongo2/v6"
)

// init registers the custom filters with pongo2.
func init() {
	// Register the "type" filter.
	// Usage: {{ Variable | type:"date" }}
	pongo2.RegisterFilter("type", filterType)

	// Register the "int" filter.
	// Usage: {{ Variable | int }}
	pongo2.RegisterFilter("int", filterInt)
}

// filterType implements the "type" filter which validates/formats values based on the argument.
// Currently supported types: "date", "filepath".
func filterType(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	val := in.String()
	typ := param.String()

	switch typ {
	case "date":
		// For the TUI form, we just pass through the value as a string.
		// Validation happens during form input collection.
		// This filter is mainly a marker for the form generator.
		return pongo2.AsValue(val), nil
	case "filepath":
		// Similar to date, this is a marker for the form generator.
		// Basic validation could be added here if needed during rendering.
		if strings.TrimSpace(val) == "" {
			return nil, &pongo2.Error{
				Sender:    "filter:type",
				OrigError: fmt.Errorf("filepath cannot be empty"),
			}
		}
		return pongo2.AsValue(val), nil
	default:
		return nil, &pongo2.Error{
			Sender:    "filter:type",
			OrigError: fmt.Errorf("unknown type: %s", typ),
		}
	}
}

// filterInt implements the "int" filter which ensures the value is an integer.
func filterInt(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	// Attempt to convert to integer.
	// If it's already an integer/number, pongo2 handles it.
	// If it's a string representation of a number, we try to convert.
	// Since in.Integer() doesn't return an error, we do manual conversion for strictness.
	s := in.String()
	i, err := strconv.Atoi(s)
	if err != nil {
		return nil, &pongo2.Error{
			Sender:    "filter:int",
			OrigError: fmt.Errorf("value is not an integer: %v", s),
		}
	}
	return pongo2.AsValue(i), nil
}

// passthroughFilter implements a pass-through filter that returns the input value unchanged.
// This is useful when a filter definition is needed for the TUI (to trigger specific form behaviors)
// but no transformation is required during the actual template rendering.
func passthroughFilter(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return pongo2.AsValue(in), nil
}

// RenderTemplate renders the template at the given path using the provided variables.
func RenderTemplate(tmplPath string, variables map[string]string) (*app.RenderedTemplate, error) {
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

	return &app.RenderedTemplate{
		Subject: subjectOut,
		HTML:    bodyOut,
	}, nil
}

//TODO : maybe I can build a package to generate filters for render in same object has form checker for tui, it would export the required function that iwll be imported by form.go and render.go
