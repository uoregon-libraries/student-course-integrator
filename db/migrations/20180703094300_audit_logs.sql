-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `audit_logs` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `created_at` datetime,
  `ip` varchar(255) COLLATE utf8_bin,
  `login` varchar(255) COLLATE utf8_bin,
  `action` varchar(255) COLLATE utf8_bin,
  `message` mediumtext COLLATE utf8_bin,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

ALTER TABLE `audit_logs` ADD INDEX audit_logs_created_at (`created_at`);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `audit_logs`;
