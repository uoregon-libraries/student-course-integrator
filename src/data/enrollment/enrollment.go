package enrollment

import "github.com/uoregon-libraries/student-course-integrator/src/global"

// AddGTF creates a new GTF record for a course, ready to be exported on the next canvas export job
func AddGTF(courseID, userID string) error {
	var sql = "INSERT INTO enrollments (`course_id`, `user_id`, `role`, `section_id`, `status`)" +
		"VALUES(?, ?, 'gtf', '', 'active')"
	var _, err = global.DB.Exec(sql, courseID, userID)
	return err
}
