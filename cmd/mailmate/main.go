package main

import (
	"flag"
	"fmt"
	"os"

	"mailmate/internal/mailer/outlookole"
	"mailmate/internal/models"
	"mailmate/internal/runner"
)

func main() {
	// Parse CLI flags
	noPreview := flag.Bool("no-preview", false, "Skip the HTML preview step and open Outlook directly")
	template := flag.String("template", "", "Template path (skip TUI selection)")
	to := flag.String("to", "", "Recipient email address")
	cc := flag.String("cc", "", "Carbon copy recipient email address")
	bcc := flag.String("bcc", "", "Blind carbon copy recipient email address")
	kv := flag.String("kv", "", "Key-value pairs for template variables (key1='value';key2='value2')")
	flag.Parse()

	options := models.Options{
		NoPreview: *noPreview,
		Template:  *template,
		To:        *to,
		Cc:        *cc,
		Bcc:       *bcc,
		KV:        *kv,
	}

	// Initialize dependencies
	sender := outlookole.NewSender()

	// Run the application
	if err := runner.Run(sender, options); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
