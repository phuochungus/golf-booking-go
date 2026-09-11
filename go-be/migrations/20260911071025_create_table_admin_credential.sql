-- +goose Up
ALTER TABLE `admins` DROP COLUMN password;
CREATE TABLE `admin_credentials` (
    admin_id INT PRIMARY KEY REFERENCES admin(id) ON DELETE CASCADE,
    hashed_password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
-- +goose Down
ALTER TABLE `admins` ADD COLUMN password VARCHAR(255) NOT NULL;
DROP TABLE `admin_credentials`;