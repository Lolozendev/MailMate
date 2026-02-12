package outlookole

import (
	"fmt"
	"mailmate/internal/mailer"
	"mailmate/internal/models"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

// OutlookSender implements mailer.EmailSender using Outlook OLE automation.
type OutlookSender struct{}

// NewSender creates a new instance of OutlookSender.
func NewSender() mailer.EmailSender {
	return &OutlookSender{}
}

// Send creates a new email in Outlook and saves it to drafts.
// Note: Despite the method name "Send" (from the interface), this implementation
// only saves to drafts as per requirements, it does not actually send the email.
func (s *OutlookSender) Send(draft models.DraftEmail) error {
	// Initialize COM library
	// internal/runner/run.go runs this in the main goroutine so this is safe.
	if err := ole.CoInitialize(0); err != nil {
		return fmt.Errorf("OLE initialization failed: %w", err)
	}
	defer ole.CoUninitialize()

	// Create Outlook.Application object
	unknown, err := oleutil.CreateObject("Outlook.Application")
	if err != nil {
		return fmt.Errorf("failed to create Outlook.Application object. Is Outlook installed?: %w", err)
	}
	defer unknown.Release()

	outlook, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return fmt.Errorf("failed to query Outlook interface: %w", err)
	}
	defer outlook.Release()

	// CreateItem(0) creates a new MailItem
	// 0 = olMailItem
	item, err := outlook.CallMethod("CreateItem", 0)
	if err != nil {
		return fmt.Errorf("failed to create mail item: %w", err)
	}
	mailItem := item.ToIDispatch()
	defer mailItem.Release()

	// Set email properties
	if draft.To != "" {
		if _, err := oleutil.PutProperty(mailItem, "To", draft.To); err != nil {
			return fmt.Errorf("failed to set mail property %q: %w", "To", err)
		}
	}
	if draft.Cc != "" {
		if _, err := oleutil.PutProperty(mailItem, "CC", draft.Cc); err != nil {
			return fmt.Errorf("failed to set mail property %q: %w", "CC", err)
		}
	}
	if draft.Bcc != "" {
		if _, err := oleutil.PutProperty(mailItem, "BCC", draft.Bcc); err != nil {
			return fmt.Errorf("failed to set mail property %q: %w", "BCC", err)
		}
	}

	if _, err := oleutil.PutProperty(mailItem, "Subject", draft.Subject); err != nil {
		return fmt.Errorf("failed to set mail property %q: %w", "Subject", err)
	}
	if _, err := oleutil.PutProperty(mailItem, "HTMLBody", draft.HTMLBody); err != nil {
		return fmt.Errorf("failed to set mail property %q: %w", "HTMLBody", err)
	}

	// Add attachments
	if len(draft.Attachments) > 0 {
		attachments := oleutil.MustGetProperty(mailItem, "Attachments").ToIDispatch()
		// No need to release attachments IDispatch as it is managed by the parent object?
		// Actually it's better to be safe, but MustGetProperty might not return an owned reference in the same way.
		// Let's stick to standard oleutil usage.
		// However, Go-OLE docs/examples often show direct method calls on the item for simple properties,
		// but Attachments is a collection.

		for _, path := range draft.Attachments {
			// Attachments.Add(Source, Type, Position, DisplayName)
			// Source is required. Others are optional.
			_, err := attachments.CallMethod("Add", path)
			if err != nil {
				// We log/return error but maybe we should continue adding others?
				// For now let's be strict.
				return fmt.Errorf("failed to add attachment %s: %w", path, err)
			}
		}
	}

	// Save the email to Drafts folder
	if _, err := mailItem.CallMethod("Save"); err != nil {
		return fmt.Errorf("failed to save email to drafts: %w", err)
	}

	// Optionally display the email window so the user can see it
	// if _, err := mailItem.CallMethod("Display"); err != nil {
	// 	 return fmt.Errorf("failed to display email: %w", err)
	// }

	return nil
}
