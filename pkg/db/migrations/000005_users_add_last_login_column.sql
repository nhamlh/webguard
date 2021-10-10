-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN last_login TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN last_login;
-- +goose StatementEnd
