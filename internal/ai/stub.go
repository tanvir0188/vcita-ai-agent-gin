// Package ai contains the stub implementation of webhook.AIService.
// Replace this with the real implementation.
package ai

import (
	"context"
	"errors"

	"github.com/tanvir0188/vcita-ai-agent/internal/store"
	"github.com/tanvir0188/vcita-ai-agent/internal/vcita"
	"github.com/tanvir0188/vcita-ai-agent/internal/webhook"
)

// Stub satisfies webhook.AIService at compile time.
// Every method returns an error until the AI developer implements the real service.
type Stub struct{}

var _ webhook.AIService = (*Stub)(nil)

func (s *Stub) ProcessMessage(
	_ context.Context, _ string, _ []store.ConversationMessage, _ []vcita.Note, _ string,
) (string, bool, error) {
	return "", false, errors.New("ai.Stub: ProcessMessage not implemented")
}

func (s *Stub) ExtractMedication(
	_ context.Context, _ string,
) (string, int, error) {
	return "", 0, errors.New("ai.Stub: ExtractMedication not implemented")
}

func (s *Stub) SuggestAppointment(
	_ context.Context, _ webhook.AppointmentRequest,
) (*webhook.AppointmentSuggestion, error) {
	return nil, errors.New("ai.Stub: SuggestAppointment not implemented")
}
