-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `faculty_courses` (
  `login` varchar(255) COLLATE utf8_bin,
  `course_id` int(11) NOT NULL
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

ALTER TABLE `faculty_courses` ADD INDEX faculty_courses_login (`login`);

CREATE TABLE `courses` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `canvas_id` tinytext COLLATE utf8_bin,
  `description` tinytext COLLATE utf8_bin,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `faculty_courses`;
DROP TABLE `courses`;
