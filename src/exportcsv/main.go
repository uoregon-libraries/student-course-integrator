package exportcsv

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/uoregon-libraries/gopkg/fileutil"
	"github.com/uoregon-libraries/gopkg/logger"
	"github.com/uoregon-libraries/student-course-integrator/src/data/enrollment"
	"github.com/uoregon-libraries/student-course-integrator/src/global"
)

// Run implements the CSV export for our main multi-binary
func Run() {
	var csvPath = global.Conf.CanvasCSVPath
	var fname = fmt.Sprintf("enrollments-%s.csv", time.Now().Format("2006-01-02"))
	var fullPath = filepath.Join(csvPath, fname)
	logger.Infof("Writing export to %q", fullPath)
	if !fileutil.MustNotExist(fullPath) {
		logger.Fatalf("%q already exists!", fullPath)
	}
	var f, err = os.Create(fullPath)
	if err != nil {
		logger.Fatalf("Unable to create enrollments output file %q: %s", fullPath, err)
	}

	var rows int
	rows, err = enrollment.ExportCSV(f)
	if err != nil {
		logger.Fatalf("Unable to write enrollments CSV to file %q: %s", fullPath, err)
	}

	logger.Infof("Done: %d row(s) written", rows)
}
