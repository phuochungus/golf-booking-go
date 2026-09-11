-- +goose Up
ALTER TABLE admins MODIFY COLUMN organization_id INT NULL;
-- +goose Down
ALTER TABLE admins MODIFY COLUMN organization_id INT NOT NULL;