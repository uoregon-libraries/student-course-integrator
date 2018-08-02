package exportcsv

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/uoregon-libraries/gopkg/fileutil"
	"github.com/uoregon-libraries/gopkg/logger"
	"github.com/uoregon-libraries/student-course-integrator/src/data/enrollment"
	"github.com/uoregon-libraries/student-course-integrator/src/global"
	"github.com/uoregon-libraries/student-course-integrator/src/service"
)

// writeCapturer just lets us grab the data being written to the file and store
// it in memory so we can send it off to the API endpoint without re-reading
// the written file.
type writeCapturer struct {
	fileutil.WriteCancelCloser
	buf []byte
}

// Write stores the bytes locally while also writing to the WriteCancelCloser
func (w *writeCapturer) Write(p []byte) (n int, err error) {
	w.buf = append(w.buf, p...)
	return w.WriteCancelCloser.Write(p)
}

// RunAPIExport generates a CSV export and auto-submits it to Canvas
func RunAPIExport() {
	var fullPath = csvFilename()
	var f = fileutil.NewSafeFile(fullPath)

	logger.Infof("Writing export to %q", fullPath)
	var w = &writeCapturer{f, nil}
	var export, err = enrollment.ExportCSV(w, fullPath)
	if err != nil {
		logger.Fatalf("Unable to generate CSV export: %s", err)
	}
	if export.Rows == 0 {
		logger.Infof("Done: no data was found to export")
		export.Close()
		return
	}

	logger.Infof("CSV file done: %d rows written", export.Rows)
	logger.Infof("Preparing to connect to Canvas endpoint")

	var fatalf = func(format string, args ...interface{}) {
		export.Cancel()
		logger.Fatalf(format, args...)
	}

	logger.Infof("Building multipart form")
	var formbuf bytes.Buffer
	var form = multipart.NewWriter(&formbuf)
	var fwriter io.Writer
	fwriter, err = form.CreateFormFile("attachment", "enrollments.csv")
	if err != nil {
		fatalf("Unable to build form for Canvas API request: %s", err)
	}
	_, err = io.Copy(fwriter, bytes.NewReader(w.buf))
	if err != nil {
		fatalf("Unable to build form for Canvas API request: %s", err)
	}
	err = form.Close()
	if err != nil {
		fatalf("Unable to build form for Canvas API request: %s", err)
	}

	logger.Infof("Building HTTP request")
	var c = new(http.Client)
	var r *http.Request
	r, err = http.NewRequest("POST", global.Conf.CanvasAPIURL, &formbuf)
	if err != nil {
		fatalf("Unable to set up Canvas export http request: %s", err)
	}
	r.Header.Set("Content-Type", form.FormDataContentType())
	err = service.ApplyHeaders(r, global.Conf.CanvasAPIHeaders)
	if err != nil {
		fatalf("Unable to parse Canvas API call headers: %s", err)
	}

	logger.Infof("Sending CSV to Canvas")
	var resp *http.Response
	resp, err = c.Do(r)
	if err != nil {
		fatalf("Error in Canvas API call: %s", err)
	}

	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fatalf("Unable to read Canvas API response: %s", err)
	}

	logger.Infof("Canvas API response: %q", string(data))
	resp.Body.Close()

	err = export.Close()
	if err != nil {
		fatalf("Unable to finalize Canvas data locally (but it may still have been sent successfully): %s", err)
	}

	logger.Infof("Done")
}
