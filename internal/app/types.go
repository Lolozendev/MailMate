package app

// TemplateRef represents a template file found in the templates directory.
type TemplateRef struct {
	// Name is the display name of the template (e.g. "invitation.html").
	Name string
	// Path is the relative or absolute path to the template file.
	Path string
}

// TemplateFilter represents a filter applied to a variable in the template.
type TemplateFilter struct {
	// Name is the name of the filter (e.g. "type", "int").
	Name string
	// Arg is the optional argument provided to the filter.
	Arg string
}

// TemplateVariable represents a variable detected in the template.
type TemplateVariable struct {
	// Name is the PascalCase identifier of the variable.
	Name string
	// Filters is the list of filters applied to this variable.
	Filters []TemplateFilter
}

// UserInput represents the values collected from the user via the TUI.
type UserInput struct {
	// Values maps variable names (TemplateVariable.Name) to user-provided values.
	Values map[string]string
}

// RenderedTemplate represents the result of rendering a template.
type RenderedTemplate struct {
	// Subject is the rendered subject line.
	Subject string
	// HTML is the final HTML string to be used as the email body.
	HTML string
}

// DraftEmail represents the email draft command to be sent to the sender.
type DraftEmail struct {
	// To is the recipient email address.
	To string
	// Cc is the carbon copy recipient email address.
	Cc string
	// Bcc is the blind carbon copy recipient email address.
	Bcc string
	// Subject is the email subject.
	Subject string
	// HTMLBody is the HTML content of the email.
	HTMLBody string
}

// Options represents the application configuration options.
type Options struct {
	// NoPreview indicates whether to skip the HTML preview step.
	NoPreview bool
	// To is the default recipient email address from flags.
	To string
	// Cc is the default carbon copy recipient email address from flags.
	Cc string
	// Bcc is the default blind carbon copy recipient email address from flags.
	Bcc string
}
