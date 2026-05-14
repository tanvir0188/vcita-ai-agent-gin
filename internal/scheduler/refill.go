// Package scheduler runs periodic background jobs using robfig/cron.
// The cron library manages goroutine lifecycle internally.
// We never use raw go + time.Sleep for scheduling.
package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"github.com/tanvir0188/vcita-ai-agent/internal/audit"
	"github.com/tanvir0188/vcita-ai-agent/internal/notify"
	"github.com/tanvir0188/vcita-ai-agent/internal/store"
	"github.com/tanvir0188/vcita-ai-agent/internal/vcita"
)

// Scheduler wraps the cron runner and its dependencies.
type Scheduler struct {
	cron     *cron.Cron
	db       *store.DB
	vc       *vcita.APIClient
	notifier *notify.Notifier
	auditor  *audit.Logger
	log      *zap.Logger
}

// New creates a Scheduler. Call Start() to begin running jobs.
func New(db *store.DB, vc *vcita.APIClient, n *notify.Notifier, a *audit.Logger, log *zap.Logger) *Scheduler {
	return &Scheduler{
		cron:     cron.New(cron.WithLogger(cron.PrintfLogger(newCronZapLogger(log)))),
		db:       db,
		vc:       vc,
		notifier: n,
		auditor:  a,
		log:      log,
	}
}

// Start registers jobs and begins the scheduler. Non-blocking.
func (s *Scheduler) Start() error {
	// Daily at 08:00: check and send due refill reminders
	if _, err := s.cron.AddFunc("0 8 * * *", s.checkRefillReminders); err != nil {
		return fmt.Errorf("scheduler: register refill job: %w", err)
	}
	s.cron.Start()
	s.log.Info("scheduler started", zap.String("jobs", "refill_check@08:00"))
	return nil
}

// Stop waits for running jobs to finish then stops the scheduler.
func (s *Scheduler) Stop() {
	ctx := s.cron.Stop()
	select {
	case <-ctx.Done():
	case <-time.After(30 * time.Second):
		s.log.Warn("scheduler: stop timed out")
	}
}

// checkRefillReminders is the cron job body.
// GORM fetches due records; for each one we send a staff alert and patient outreach.
func (s *Scheduler) checkRefillReminders() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	due, err := s.db.GetDueReminders()
	if err != nil {
		s.log.Error("scheduler: GetDueReminders", zap.Error(err))
		return
	}
	s.log.Info("scheduler: refill check", zap.Int("due", len(due)))

	for _, r := range due {
		if err := s.processOne(ctx, r); err != nil {
			s.log.Error("scheduler: processOne",
				zap.Uint("id", r.ID), zap.Error(err))
			// Continue — don't let one failure block the rest
		}
	}
}

func (s *Scheduler) processOne(ctx context.Context, r store.RefillScheduleView) error {
	dateStr := r.RefillDate.Format("January 2, 2006")

	// 1. Staff alert email
	if err := s.notifier.SendRefillReminder(ctx, r.ClientID); err != nil {
		s.log.Warn("scheduler: staff email failed", zap.Error(err))
		// Non-fatal — still attempt patient outreach
	}

	// 2. Patient outreach via vcita
	msg := fmt.Sprintf(
		"Hello! This is a friendly reminder that your prescription for %s may need to be refilled around %s. "+
			"Please contact our office if you have questions or would like to schedule an appointment.",
		r.Medication, dateStr,
	)
	if err := s.vc.SendMessage(ctx, vcita.SendMessageRequest{
		ClientID: r.ClientID,
		Body:     msg,
	}); err != nil {
		return fmt.Errorf("send patient outreach: %w", err)
	}

	// 3. Mark as sent so it doesn't fire again
	if err := s.db.MarkReminderSent(r.ID); err != nil {
		return fmt.Errorf("mark reminder sent: %w", err)
	}

	s.auditor.Log("refill_reminder_sent", r.ClientID, "scheduler",
		fmt.Sprintf("record_id=%d", r.ID))
	return nil
}

// cronZapLogger bridges robfig/cron's Printf-based logger to zap.
type cronZapLogger struct{ log *zap.Logger }

func newCronZapLogger(log *zap.Logger) *cronZapLogger { return &cronZapLogger{log} }

func (l *cronZapLogger) Printf(format string, args ...interface{}) {
	l.log.Sugar().Debugf("[cron] "+format, args...)
}
