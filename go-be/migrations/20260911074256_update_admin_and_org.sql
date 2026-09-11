-- +goose Up
ALTER TABLE `admins` DROP COLUMN `root`;
ALTER TABLE `organizations` ADD COLUMN `root_admin_id` INT;

-- +goose Down
ALTER TABLE `admins` ADD COLUMN `root` BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE `organizations` DROP COLUMN `root_admin_id`;
