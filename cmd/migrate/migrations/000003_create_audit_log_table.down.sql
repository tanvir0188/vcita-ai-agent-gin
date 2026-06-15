CREATE TABLE audit_logs (
	id INTEGER PRIMARY KEY AUTOINCREMENT,

	client_id TEXT NOT NULL,
	action TEXT NOT NULL,
	event_type TEXT NOT NULL DEFAULT '',
	outcome TEXT NOT NULL,
	actor_type TEXT NOT NULL DEFAULT 'system',
	error_msg TEXT DEFAULT '',

	created_at DATETIME,
	updated_at DATETIME,
	deleted_at DATETIME
);

CREATE INDEX idx_audit_logs_client_id ON audit_logs(client_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_event_type ON audit_logs(event_type);