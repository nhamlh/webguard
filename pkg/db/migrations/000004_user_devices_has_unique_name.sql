-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX user_device_name_unique_index
ON devices(user_id, name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
--SELECT 'down SQL query';
-- +goose StatementEnd
