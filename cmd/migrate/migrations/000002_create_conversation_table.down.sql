CREATE TABLE conversations (
	id INTEGER PRIMARY KEY AUTOINCREMENT,

	conversation_id TEXT NOT NULL UNIQUE,

	last_message_id TEXT,
	last_customer_message_id TEXT,
	last_staff_message_id TEXT,
	last_message_from TEXT,

	last_customer_message_at DATETIME,
	last_staff_message_at DATETIME,

	human_active BOOLEAN DEFAULT false,
	human_active_at DATETIME,
	human_active_until DATETIME,

	conversation_version INTEGER DEFAULT 0,

	ai_reply_pending BOOLEAN DEFAULT false,
	ai_reply_generating BOOLEAN DEFAULT false,

	has_escalated BOOLEAN DEFAULT false,

	auto_reply_off_until DATETIME,
	pending_message_id TEXT,

	created_at DATETIME,
	updated_at DATETIME,
	deleted_at DATETIME
);

CREATE INDEX idx_conversations_human_active_at ON conversations(human_active_at);
CREATE INDEX idx_conversations_human_active_until ON conversations(human_active_until);
CREATE INDEX idx_conversations_ai_reply_pending ON conversations(ai_reply_pending);