package enrollment

import (
	"encoding/csv"
	"errors"
	"fmt"
	"time"

	"github.com/Nerdmaster/magicsql"
	"github.com/uoregon-libraries/gopkg/fileutil"
	"github.com/uoregon-libraries/student-course-integrator/src/global"
	"github.com/uoregon-libraries/student-course-integrator/src/roles"
)

// Add creates a new record for a course, ready to be exported on the next canvas export job
func Add(courseID, userID, role string) error {
	if !roles.IsValid(role) {
		return errors.New("did not attempt to write to database, " + role + " is not a valid role.")
	}
	var sql = "INSERT INTO enrollments (`course_id`, `user_id`, `role`, `section_id`, `status`)" +
		"VALUES(?, ?, ?, '', 'active')"
	var _, err = global.DB.Exec(sql, courseID, userID, role)
	return err

}

// Export wraps the database and filesystem information surrounding a CSV
// export of enrollments
type Export struct {
	op   *magicsql.Operation
	w    fileutil.WriteCancelCloser
	ID   int64
	Rows int
}

// Close ends the database transaction and finalizes the export's file if there
// were no database errors.  As a special case, if no data was available to
// export, no file is writen and the database changes are reverted to avoid
// empty files and db export records.
func (e *Export) Close() error {
	if e.Rows == 0 {
		e.Cancel()
		return nil
	}

	e.op.EndTransaction()
	if e.op.Err() != nil {
		e.w.Cancel()
		return e.op.Err()
	}

	return e.w.Close()
}

// Cancel cancels the export file's write and rolls back the database changes
func (e *Export) Cancel() {
	e.op.Rollback()
	e.w.Cancel()
}

// ExportCSV finds all unexported enrollments in the database and writes them
// out in CSV format to w, including the CSV header (course_id, user_id, role,
// section_id, status).
//
// The path is only necessary to associate the filename with the export in the
// database.  the file is not accessed or manipulated except through the
// WriteCancelCloser.
//
// On succes, the enrollments records will be tied to a new canvas_exports
// record and the returned Export will be ready for finalizing.  The caller
// must either commit or roll back the database transaction, as makes sense for
// whatever other operations are happening (e.g., the Canvas API call).
func ExportCSV(w fileutil.WriteCancelCloser, path string) (*Export, error) {
	var export = &Export{w: w}
	var cw = csv.NewWriter(w)

	var err = cw.Write([]string{"course_id", "user_id", "role", "section_id", "status"})
	if err != nil {
		return nil, fmt.Errorf("enrollment: cannot write CSV header: %s", err)
	}

	// MagicSQL wrapper lets us handle error checking in "clumps"
	var db = magicsql.Wrap(global.DB)

	// Create an export record and tie all unexported enrollments to it
	export.op = db.Operation()
	export.op.BeginTransaction()

	var fail = func(err error) (*Export, error) {
		export.Cancel()
		return nil, err
	}

	var result = export.op.Exec("INSERT INTO canvas_exports (exported_at, path) VALUES (?, ?)",
		time.Now(), path)
	export.ID = result.LastInsertId()
	export.op.Exec("UPDATE enrollments SET canvas_export_id = ? WHERE canvas_export_id = 0", export.ID)
	if export.op.Err() != nil {
		return fail(fmt.Errorf("enrollment: database prep error: %s", export.op.Err()))
	}

	var rows = export.rowsStrings()
	if export.op.Err() != nil {
		return fail(fmt.Errorf("enrollment: error reading enrollments table: %s", export.op.Err()))
	}

	err = cw.WriteAll(rows)
	if err != nil {
		return fail(fmt.Errorf("enrollment: error writing enrollments csv: %s", err))
	}

	export.Rows = len(rows)
	return export, nil
}

func (e *Export) rowsStrings() [][]string {
	var rows = e.op.Query(`
		SELECT course_id, user_id, role, section_id, status
		FROM enrollments
		WHERE canvas_export_id = ?
	`, e.ID)

	var csvRows [][]string
	var courseID, userID, role, sectionID, status string
	for rows.Next() {
		rows.Scan(&courseID, &userID, &role, &sectionID, &status)
		// HACK: At some point GEs become TAs.  This allows us to keep the old
		// version going, but with the export turned off, until Canvas is ready.
		if role == "GE" {
			role = roles.TA
		}
		csvRows = append(csvRows, []string{courseID, userID, role, sectionID, status})
	}
	return csvRows
}
