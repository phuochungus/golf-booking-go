-- +goose Up
ALTER TABLE admins ADD COLUMN root BOOLEAN DEFAULT false;
ALTER TABLE organizations DROP COLUMN phone;
ALTER TABLE organizations DROP COLUMN email;
-- +goose Down
ALTER TABLE admins DROP COLUMN root;
ALTER TABLE organizations ADD COLUMN phone VARCHAR(20);
ALTER TABLE organizations ADD COLUMN email VARCHAR(100);
