package app

import (
	"fmt"
	"os"
)

// Exit codes
const (
	ExitSuccess = 0
	ExitError   = 1
)

// Exit prints the error message to stderr and exits with status 1
func Exit(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(ExitError)
}

// ExitWithAction prints the error message and a suggested action to stderr, then exits with status 1
func ExitWithAction(err error, action string) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	fmt.Fprintf(os.Stderr, "Suggested Action: %s\n", action)
	os.Exit(ExitError)
}
