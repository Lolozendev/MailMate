package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"mailmate/internal/mailer/outlookole"
	"mailmate/internal/models"
	"mailmate/internal/runner"
)

func main() {
	// Pre-process args to handle --kv without value
	// If --kv is present without a value, insert an empty string
	var kvExplicitlyProvided bool
	for i := 1; i < len(os.Args); i++ {
		if os.Args[i] == "--kv" || os.Args[i] == "-kv" {
			kvExplicitlyProvided = true
			// Check if next arg exists and doesn't start with "-"
			if i+1 >= len(os.Args) || strings.HasPrefix(os.Args[i+1], "-") {
				// No value provided, insert empty string
				args := make([]string, 0, len(os.Args)+1)
				args = append(args, os.Args[:i+1]...)
				args = append(args, "")
				args = append(args, os.Args[i+1:]...)
				os.Args = args
			}
			break
		}
	}

	// Parse CLI flags
	noPreview := flag.Bool("no-preview", false, "Skip the HTML preview step and open Outlook directly")
	template := flag.String("template", "", "Template path (skip TUI selection)")
	to := flag.String("to", "", "Recipient email address")
	cc := flag.String("cc", "", "Carbon copy recipient email address")
	bcc := flag.String("bcc", "", "Blind carbon copy recipient email address")
	kv := flag.String("kv", "", "Key-value pairs for template variables (key1='value';key2='value2')")
	flag.Parse()

	// Determine if --kv flag was explicitly provided
	var kvPtr *string
	if kvExplicitlyProvided {
		kvPtr = kv
	}

	options := models.Options{
		NoPreview: *noPreview,
		Template:  *template,
		To:        *to,
		Cc:        *cc,
		Bcc:       *bcc,
		KV:        kvPtr,
	}

	// Initialize dependencies
	sender := outlookole.NewSender()

	// Run the application
	if err := runner.Run(sender, options); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
