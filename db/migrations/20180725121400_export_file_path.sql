-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE `canvas_exports` ADD COLUMN `path` varchar(255) COLLATE utf8_bin;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
ALTER TABLE `canvas_exports` DROP COLUMN `path`;
