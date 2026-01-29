package templates

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
	"mailmate/internal/models"
)

// ParsedTemplateFile represents the separated subject and body from a template file.
type ParsedTemplateFile struct {
	Subject string
	Body    string
	// Default recipients from template frontmatter
	To  string
	Cc  string
	Bcc string
}

// ParseTemplateFile reads a template file, extracts the frontmatter (if any),
// and returns the parsed subject and body.
func ParseTemplateFile(path string) (*ParsedTemplateFile, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Check for frontmatter start
	if !bytes.HasPrefix(content, []byte("---")) {
		return &ParsedTemplateFile{
			Body: string(content),
		}, nil
	}

	// Find the end of the frontmatter
	// Look for \n--- or \r\n---
	endIndex := -1
	offset := 3 // Skip initial ---

	if idx := bytes.Index(content[offset:], []byte("\n---")); idx != -1 {
		endIndex = offset + idx
	} else if idx := bytes.Index(content[offset:], []byte("\r\n---")); idx != -1 {
		endIndex = offset + idx
	}

	if endIndex == -1 {
		// Started with --- but no closing --- found.
		return nil, fmt.Errorf("malformed frontmatter: missing closing '---'")
	}

	yamlData := content[offset:endIndex]

	// Body starts after the closing ---
	// The closing sequence is \n--- or \r\n---. We need to skip that plus potentially the newline after it.
	bodyStart := endIndex + 4 // \n--- is 4 chars
	if content[endIndex] == '\r' {
		bodyStart = endIndex + 5 // \r\n--- is 5 chars
	}

	// Skip potential newline after the second ---
	if bodyStart < len(content) && content[bodyStart] == '\r' {
		bodyStart++
	}
	if bodyStart < len(content) && content[bodyStart] == '\n' {
		bodyStart++
	}

	var meta struct {
		Subject string `yaml:"subject"`
		To      string `yaml:"to"`
		Cc      string `yaml:"cc"`
		Bcc     string `yaml:"bcc"`
	}
	if err := yaml.Unmarshal(yamlData, &meta); err != nil {
		return nil, fmt.Errorf("parsing frontmatter yaml: %w", err)
	}

	return &ParsedTemplateFile{
		Subject: meta.Subject,
		Body:    string(content[bodyStart:]),
		To:      meta.To,
		Cc:      meta.Cc,
		Bcc:     meta.Bcc,
	}, nil
}

// ParseTemplate reads a template file and extracts variables and their filters.
// It parses both the frontmatter Subject and the Body.
func ParseTemplate(path string) ([]models.TemplateVariable, error) {
	parsed, err := ParseTemplateFile(path)
	if err != nil {
		return nil, err
	}

	// Combine subject, recipients, and body for variable scanning
	combined := parsed.Subject + "\n" + parsed.To + "\n" + parsed.Cc + "\n" + parsed.Bcc + "\n" + parsed.Body

	// Regex to find {{ VariableName | filters... }}
	// Captures:
	// 1. Variable Name (PascalCase)
	// 2. The rest of the tag content (potential filters)
	varRegex := regexp.MustCompile(`{{\s*([A-Z][a-zA-Z0-9]*)(.*?)\s*}}`)

	matches := varRegex.FindAllStringSubmatch(combined, -1)

	var variables []models.TemplateVariable
	seen := make(map[string]bool)

	for _, match := range matches {
		name := match[1]
		filterStr := match[2]

		if seen[name] {
			continue
		}
		seen[name] = true

		tv := models.TemplateVariable{
			Name: name,
		}

		// Parse filters if present
		if strings.TrimSpace(filterStr) != "" {
			tv.Filters = parseFilters(filterStr)
		}

		variables = append(variables, tv)
	}

	return variables, nil
}

// parseFilters extracts filters from a string like "| type:'date' | default:'now'"
func parseFilters(raw string) []models.TemplateFilter {
	var filters []models.TemplateFilter

	// Regex matches: | name (: arg)?
	// Captures:
	// 1. Filter Name
	// 2. Double quoted arg content
	// 3. Single quoted arg content
	// 4. Unquoted arg content
	filterRegex := regexp.MustCompile(`\|\s*([a-zA-Z0-9_]+)(?:\s*:\s*(?:"([^"]*)"|'([^']*)'|([^|\s]+)))?`)

	matches := filterRegex.FindAllStringSubmatch(raw, -1)

	for _, m := range matches {
		fName := m[1]
		fArg := ""

		if m[2] != "" {
			fArg = m[2]
		} else if m[3] != "" {
			fArg = m[3]
		} else if m[4] != "" {
			fArg = m[4]
		}

		filters = append(filters, models.TemplateFilter{
			Name: fName,
			Arg:  fArg,
		})
	}

	return filters
}
