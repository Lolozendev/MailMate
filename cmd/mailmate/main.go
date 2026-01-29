package main

import (
	"flag"
	"fmt"
	"os"

	"mailmate/internal/app"
	"mailmate/internal/mailer/outlookole"
	"mailmate/internal/runner"
)

func main() {
	// Parse CLI flags
	noPreview := flag.Bool("no-preview", false, "Skip the HTML preview step and open Outlook directly")
	to := flag.String("to", "", "Recipient email address")
	cc := flag.String("cc", "", "Carbon copy recipient email address")
	bcc := flag.String("bcc", "", "Blind carbon copy recipient email address")
	flag.Parse()

	options := app.Options{
		NoPreview: *noPreview,
		To:        *to,
		Cc:        *cc,
		Bcc:       *bcc,
	}

	// Initialize dependencies
	sender := outlookole.NewSender()

	// Run the application
	if err := runner.Run(sender, options); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
