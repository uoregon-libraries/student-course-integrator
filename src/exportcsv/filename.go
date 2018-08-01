package exportcsv

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/uoregon-libraries/gopkg/fileutil"
	"github.com/uoregon-libraries/student-course-integrator/src/global"
)

// csvSafeFile returns a filename for use in the CSV export using the
// configured Canvas CSV output path, the current date and time, and a sequence
// number.  If after multiple sequenced version of the file, all tested
// filenames exist, this function will panic.
func csvFilename() string {
	var csvPath = global.Conf.CanvasCSVPath
	for x := 0; x < 10; x++ {
		var fname = fmt.Sprintf("enrollments-%s-%d.csv", time.Now().Format("2006-01-02-150405"), x)
		var fullPath = filepath.Join(csvPath, fname)
		if !fileutil.Exists(fullPath) {
			return fullPath
		}
	}
	panic("unable to create a unique filename")
}
