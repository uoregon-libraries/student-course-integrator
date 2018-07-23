package importcsv

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Nerdmaster/magicsql"
	"github.com/uoregon-libraries/gopkg/logger"
	"github.com/uoregon-libraries/student-course-integrator/src/global"
	"github.com/uoregon-libraries/student-course-integrator/src/service"
)

// Run implements the CSV import for our main multi-binary
func Run() {
	var courses = readCSV("courses.csv")
	var enrollments = readCSV("enrollments.csv")

	// MagicSQL wrapper lets us defer error checking
	var db = magicsql.Wrap(global.DB)
	var op = db.Operation()
	op.BeginTransaction()
	var err = buildData(op, courses, enrollments)
	if err != nil {
		op.Rollback()
		logger.Fatalf("Error building data: %s", err)
	}
	op.EndTransaction()

	if op.Err() != nil {
		logger.Fatalf("Unable to begin DB transaction: %s", op.Err())
	}
}

func buildData(op *magicsql.Operation, courses, enrollments [][]string) error {
	// courseMap maps the string id to the database id so we can attach faculty to a course
	var courseMap = make(map[string]int64)

	// Delete and repopulate courses, mapping string ids to db ids
	op.Exec("DELETE FROM courses")
	var st = op.Prepare("INSERT INTO courses (canvas_id, description) VALUES (?, ?)")

	var expectedLen = 6
	for i, record := range courses {
		if len(record) != expectedLen {
			return fmt.Errorf("courses.csv: record %d doesn't have %d columns", i, expectedLen)
		}
		var courseID, desc, status = record[0], record[2], record[5]
		if status != "active" {
			continue
		}
		var res = st.Exec(courseID, desc)
		var dbid = res.LastInsertId()
		if dbid > 0 {
			courseMap[courseID] = dbid
		}
	}

	// Delete and repopulate faculty_courses, using the above-generated map to get db ids for courses
	op.Exec("DELETE FROM faculty_courses")
	st = op.Prepare("INSERT INTO faculty_courses (login, course_id) VALUES (?, ?)")

	// duckidMap reduces API hits by storing duckids for previously-seen banner ids
	var duckidMap = make(map[string]string)
	expectedLen = 5

	for i, record := range enrollments {
		if len(record) != expectedLen {
			return fmt.Errorf("enrollments.csv: record %d doesn't have %d columns", i, expectedLen)
		}

		// We only process active enrollments for teachers of the main course (i.e., no section)
		var courseID, userID, role, sectionID, status = record[0], record[1], record[2], record[3], record[4]
		if status != "active" || role != "teacher" || sectionID != "" {
			continue
		}

		var courseDBID, ok = courseMap[courseID]
		if !ok {
			return fmt.Errorf("enrollments: record %d's course id (%s) is unknown", i, courseID)
		}

		var duckid string
		duckid, ok = duckidMap[userID]
		if !ok {
			var err error
			duckid, err = service.BannerIDToDuckID(userID)
			if err != nil {
				return fmt.Errorf("unable to look up duckid for %s: %s", userID, err)
			}
			duckidMap[userID] = duckid
		}
		st.Exec(duckid, courseDBID)
	}

	return nil
}

func readCSV(fname string) (records [][]string) {
	var csvPath = global.Conf.BannerCSVPath
	var f, err = os.Open(filepath.Join(csvPath, fname))
	if err != nil {
		logger.Fatalf("Unable to open %q: %s", fname, err)
	}
	var r = csv.NewReader(f)
	records, err = r.ReadAll()
	if err != nil {
		logger.Fatalf("Unable to read %q: %s", fname, err)
	}

	// In our setup, we *always* have a header row, so we want to ignore it
	return records[1:]
}
