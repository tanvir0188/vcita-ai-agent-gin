CREATE TABLE users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,

	staff_uid TEXT NOT NULL UNIQUE,
	full_name TEXT,
	email TEXT NOT NULL UNIQUE,
	phone_number TEXT,
	password TEXT NOT NULL,

	is_active BOOLEAN DEFAULT false,
	is_verified BOOLEAN DEFAULT false,

	otp_code TEXT,
	otp_expires_at DATETIME,

	created_at DATETIME,
	updated_at DATETIME,
	deleted_at DATETIME
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_staff_uid ON users(staff_uid);