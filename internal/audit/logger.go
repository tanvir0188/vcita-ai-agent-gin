// Package audit provides an append-only structured audit logger for HIPAA compliance.
//
// HIPAA requires every access and modification of PHI to be traceable.
// This logger writes structured JSON lines to a dedicated audit log file
// in addition to inserting records into the audit_log database table.
//
// CRITICAL: Never pass PHI (message content, patient names, medication details)
// to audit log methods. Log event types, IDs, and actions only.
package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// Event is a single audit log entry.
type Event struct {
	Timestamp string `json:"timestamp"`
	EventType string `json:"event_type"`
	ClientID  string `json:"client_id,omitempty"` // pseudonymised reference
	Actor     string `json:"actor"`
	Detail    string `json:"detail,omitempty"`
}

// Logger writes audit events to a file and optionally a DB sink.
type Logger struct {
	mu   sync.Mutex
	file *os.File
	dbFn func(clientID, action, eventType, outcome, actorType, errMsg string) error
}

// New creates a Logger that appends to the given file path.
func New(filePath string, dbFn func(clientID, action, eventType, outcome, actorType, errMsg string) error) (*Logger, error) {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("audit: open log file %q: %w", filePath, err)
	}
	return &Logger{file: f, dbFn: dbFn}, nil
}

// Log records an audit event. Thread-safe.
func (l *Logger) Log(eventType, clientID, actor, detail string) {
	evt := Event{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		EventType: eventType,
		ClientID:  clientID,
		Actor:     actor,
		Detail:    detail,
	}
	line, _ := json.Marshal(evt)
	l.mu.Lock()
	_, _ = l.file.Write(append(line, '\n'))
	l.mu.Unlock()
	if l.dbFn != nil {
		_ = l.dbFn(clientID, "audit", eventType, "success", actor, detail)
	}
}

// Close flushes and closes the underlying file.
func (l *Logger) Close() error {
	return l.file.Close()
}
