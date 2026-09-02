-- +goose Up
ALTER TABLE facilities ADD COLUMN rating FLOAT DEFAULT 0;

-- +goose Down
ALTER TABLE facilities DROP COLUMN rating;
