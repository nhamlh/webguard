-- +goose Up
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	email TEXT NOT NULL UNIQUE,
	password TEXT,
	is_admin TEXT NOT NULL,
	auth_type INTEGER NOT NULL
);

-- +goose Down
DROP TABLE users;
