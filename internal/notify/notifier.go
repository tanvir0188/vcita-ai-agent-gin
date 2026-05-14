// Package notify handles staff alerts for escalations and refill reminders.
// Alert emails never contain PHI in the subject line.
// Body content is kept minimal and contains only what staff need to act.
package notify

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	gomail "gopkg.in/gomail.v2"
)

// Urgency level for staff alerts
type Urgency int

const (
	UrgencyLow    Urgency = iota
	UrgencyMedium Urgency = iota
	UrgencyHigh   Urgency = iota
)

// Alert is the data needed to send a staff notification.
type Alert struct {
	ClientID string  // Only ID — not name/contact (staff can look up in vcita)
	Reason   string  // Human-readable reason — no PHI
	Urgency  Urgency
}

// Notifier sends alerts to staff via email.
type Notifier struct {
	dialer     *gomail.Dialer
	fromAddr   string
	staffEmail string
	logger     *zap.Logger
}

// NewNotifier creates a Notifier from SMTP config.
func NewNotifier(host string, port int, username, password, from, staffEmail string, logger *zap.Logger) *Notifier {
	d := gomail.NewDialer(host, port, username, password)
	return &Notifier{
		dialer:     d,
		fromAddr:   from,
		staffEmail: staffEmail,
		logger:     logger,
	}
}

// AlertStaff sends an email alert to staff.
// Subject never contains PHI. Body contains client_id only (staff use vcita to look up patient).
func (n *Notifier) AlertStaff(ctx context.Context, alert Alert) error {
	urgencyLabel := urgencyLabel(alert.Urgency)
	subject := fmt.Sprintf("[%s] Patient Action Required — vcita AI Agent", urgencyLabel)

	body := fmt.Sprintf(`
A patient interaction requires your attention.

Client ID: %s
Reason: %s
Urgency: %s

Please log into vcita to review this patient's conversation and take appropriate action.

This is an automated message from the vcita AI Agent. Do not reply to this email.
`, alert.ClientID, alert.Reason, urgencyLabel)

	m := gomail.NewMessage()
	m.SetHeader("From", n.fromAddr)
	m.SetHeader("To", n.staffEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	if err := n.dialer.DialAndSend(m); err != nil {
		n.logger.Error("failed to send staff alert email",
			zap.String("urgency", urgencyLabel),
			// Do NOT log client_id here — associate via audit log
			zap.Error(err),
		)
		return fmt.Errorf("notify: failed to send alert: %w", err)
	}

	n.logger.Info("staff alert sent",
		zap.String("urgency", urgencyLabel),
		zap.String("reason", alert.Reason),
	)
	return nil
}

// SendRefillReminder sends a refill reminder alert to staff.
// Medication details are not included — staff use vcita to review.
func (n *Notifier) SendRefillReminder(ctx context.Context, clientID string) error {
	return n.AlertStaff(ctx, Alert{
		ClientID: clientID,
		Reason:   "Medication refill due within 7 days — outreach recommended",
		Urgency:  UrgencyMedium,
	})
}

func urgencyLabel(u Urgency) string {
	switch u {
	case UrgencyHigh:
		return "HIGH"
	case UrgencyMedium:
		return "MEDIUM"
	default:
		return "LOW"
	}
}
