-- +goose Up
CREATE TABLE IF NOT EXISTS devices (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,
	name TEXT NOT NULL,
	private_key TEXT NOT NULL,
	allowed_ips TEXT NOT NULL,

	FOREIGN KEY(user_id) REFERENCES users(id)
);

-- +goose Down
DROP TABLE devices;
