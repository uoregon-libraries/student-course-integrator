package enrollment

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"time"

	"github.com/Nerdmaster/magicsql"
	"github.com/uoregon-libraries/student-course-integrator/src/global"
)

// AddGE creates a new GE record for a course, ready to be exported on the next canvas export job
func AddGE(courseID, userID string) error {
	var sql = "INSERT INTO enrollments (`course_id`, `user_id`, `role`, `section_id`, `status`)" +
		"VALUES(?, ?, 'GE', '', 'active')"
	var _, err = global.DB.Exec(sql, courseID, userID)
	return err
}

// ExportCSV finds all unexported enrollments in the database and writes them
// out in CSV format to w, including the CSV header (course_id, user_id, role,
// section_id, status).  If successful, the enrollments records will be tied to
// a new canvas_exports record.
func ExportCSV(w io.WriteCloser, path string) (rows int, err error) {
	var buf = new(bytes.Buffer)
	var cw = csv.NewWriter(buf)

	err = cw.Write([]string{"course_id", "user_id", "role", "section_id", "status"})
	if err != nil {
		return 0, fmt.Errorf("enrollment: cannot write CSV header: %s", err)
	}

	// MagicSQL wrapper lets us handle error checking in "clumps"
	var db = magicsql.Wrap(global.DB)

	// Create an export record and tie all unexported enrollments to it.
	// Assuming the CSV writer doesn't fail, the database manipulation is done.
	var op = db.Operation()
	op.BeginTransaction()

	// Queue up a transaction rollback; this ensures we rollback on any return
	// unless we've explicitly commited
	defer op.Rollback()
	var result = op.Exec("INSERT INTO canvas_exports (exported_at, path) VALUES (?, ?)",
		time.Now(), path)
	var exportID = result.LastInsertId()
	op.Exec("UPDATE enrollments SET canvas_export_id = ? WHERE canvas_export_id = 0", exportID)
	if op.Err() != nil {
		return 0, fmt.Errorf("enrollment: database prep error: %s", op.Err())
	}

	var exportRows = getStringRowsForID(op, exportID)
	if op.Err() != nil {
		return 0, fmt.Errorf("enrollment: error reading enrollments table: %s", op.Err())
	}

	err = cw.WriteAll(exportRows)
	if err != nil {
		return 0, fmt.Errorf("enrollment: error writing enrollments csv: %s", err)
	}

	_, err = io.Copy(w, buf)
	if err != nil {
		return 0, fmt.Errorf("enrollment: error writing enrollments csv: %s", err)
	}

	err = w.Close()
	if err != nil {
		return 0, fmt.Errorf("enrollment: error closing enrollments csv: %s", err)
	}

	// Let's not waste a DB record on an empty export
	if len(exportRows) == 0 {
		return 0, nil
	}

	op.EndTransaction()

	return len(exportRows), nil
}

func getStringRowsForID(op *magicsql.Operation, id int64) [][]string {
	var rows = op.Query(`
		SELECT course_id, user_id, role, section_id, status
		FROM enrollments
		WHERE canvas_export_id = ?
	`, id)

	var csvRows [][]string
	var courseID, userID, role, sectionID, status string
	for rows.Next() {
		rows.Scan(&courseID, &userID, &role, &sectionID, &status)
		csvRows = append(csvRows, []string{courseID, userID, role, sectionID, status})
	}
	return csvRows
}
