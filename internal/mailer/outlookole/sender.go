package outlookole

import (
	"fmt"
	"mailmate/internal/app"
	"mailmate/internal/mailer"

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
func (s *OutlookSender) Send(draft app.DraftEmail) error {
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
		oleutil.PutProperty(mailItem, "To", draft.To)
	}
	if draft.Cc != "" {
		oleutil.PutProperty(mailItem, "CC", draft.Cc)
	}
	if draft.Bcc != "" {
		oleutil.PutProperty(mailItem, "BCC", draft.Bcc)
	}

	oleutil.PutProperty(mailItem, "Subject", draft.Subject)
	oleutil.PutProperty(mailItem, "HTMLBody", draft.HTMLBody)

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
