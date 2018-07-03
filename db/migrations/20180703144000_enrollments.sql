-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `enrollments` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `canvas_export_id` int(11) NOT NULL DEFAULT 0,
  `course_id` varchar(255) COLLATE utf8_bin,
  `user_id` varchar(255) COLLATE utf8_bin,
  `role` varchar(255) COLLATE utf8_bin,
  `section_id` varchar(255) COLLATE utf8_bin,
  `status` varchar(255) COLLATE utf8_bin,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

CREATE TABLE `canvas_exports` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `exported_at` datetime,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `enrollments`;
DROP TABLE `canvas_exports`;
