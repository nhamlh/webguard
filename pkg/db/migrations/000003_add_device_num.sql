-- +goose Up
ALTER TABLE devices
ADD COLUMN num INTEGER NOT NULL;

CREATE UNIQUE INDEX devices_num_unique ON devices(num);

-- +goose Down
DROP INDEX devices_num_unique;

ALTER TABLE devices
DROP num;
