package mailer

import (
	"mailmate/internal/app"
)

// EmailSender defines the interface for sending emails.
// It decouples the application logic from the specific email sending implementation (e.g. Outlook OLE).
type EmailSender interface {
	Send(draft app.DraftEmail) error
}
