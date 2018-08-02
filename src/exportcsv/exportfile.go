package exportcsv

import (
	"github.com/uoregon-libraries/gopkg/fileutil"
	"github.com/uoregon-libraries/gopkg/logger"
	"github.com/uoregon-libraries/student-course-integrator/src/data/enrollment"
)

// RunFileExport implements the CSV file export (without a canvas upload) for
// our main multi-binary
func RunFileExport() {
	var fullPath = csvFilename()
	var f = fileutil.NewSafeFile(fullPath)

	logger.Infof("Writing export to %q", fullPath)
	var export, err = enrollment.ExportCSV(f, fullPath)
	if err == nil {
		err = export.Close()
	}
	if err != nil {
		logger.Fatalf("Unable to generate CSV export: %s", err)
	}

	if export.Rows == 0 {
		logger.Infof("Done: no data was found to export")
	} else {
		logger.Infof("Done: %d row(s) written", export.Rows)
	}
}
