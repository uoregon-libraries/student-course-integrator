package sciserver

import (
	"fmt"
	"net/http"

	"github.com/uoregon-libraries/gopkg/tmpl"
	"github.com/uoregon-libraries/gopkg/webutil"
	"github.com/uoregon-libraries/student-course-integrator/src/data/audit"
	"github.com/uoregon-libraries/student-course-integrator/src/data/enrollment"
	"github.com/uoregon-libraries/student-course-integrator/src/data/user"
)

// homeVars stores data for the home template to use
type homeVars struct {
	Alert string     // Alert, if set, is displayed as a bootstrappy alert on the page
	Info  string     // Info, if set, is displayed as a bootstrappy info section on the page
	User  *user.User // User is set to the logged-in user
	Form  *form      // Form stores all the submitted form data, if any
}

// homeHandler encapsulates basic data and functionality for handling input and
// rendering output
type homeHandler struct {
	formTemplate    *tmpl.Template
	confirmTemplate *tmpl.Template
}

// hHome generates a homeHandler structure for serving pages via http.ServeHTTP
func hHome() *homeHandler {
	var r = layout.Clone()
	return &homeHandler{
		formTemplate:    r.MustBuild("home.go.html"),
		confirmTemplate: r.MustBuild("confirm_duckid.go.html"),
	}
}

// ServeHTTP implements http.Handler for homeHandler
func (h *homeHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	warmCache()
	var r = respond(w, req, h)
	if req.Method == "POST" {
		r.processSubmission()
		return
	}
	r.serveForm()
}

func (r *responder) processSubmission() {
	var f, err = r.getForm()
	if err != nil {
		r.render500(fmt.Errorf("unable to instantiate form data: %s", err))
		return
	}

	// Set up all the key/value pairs we want to store in the audit log
	var auditVals = map[string]string{
		"crn":     f.CRN,
		"duckid":  f.DuckID,
		"confirm": f.Confirm,
		"role":    "GE",
	}

	// Explicit rejection of duckid was requested: re-render the form
	if f.Confirm == "0" {
		audit.Log(r.vars.User, audit.ActionRejectSubmission, auditVals)
		r.render(r.hh.formTemplate)
		return
	}

	// Errors: re-render the form
	if len(f.errors) > 0 {
		auditVals["error"] = f.errorString()
		audit.Log(r.vars.User, audit.ActionInvalidSubmission, auditVals)
		r.vars.Alert = fmt.Sprintf("Error: %s", f.errorString())
		r.render(r.hh.formTemplate)
		return
	}

	// Require "confirm" to be exactly the string "1" so that we err on the side of not adding GEs
	if f.Confirm == "1" {
		err = enrollment.AddGE(f.CRN, f.GE.BannerID)
		if err != nil {
			r.render500(fmt.Errorf("unable to write enrollment data to database: %s", err))
			return
		}
		audit.Log(r.vars.User, audit.ActionConfirmSubmission, auditVals)
		var s = getSession(r.w, r.req)
		s.SetInfoFlash(fmt.Sprintf(`%s (%s) added to %s`,
			f.GE.DisplayName, f.GE.DuckID, f.Course.Description))
		http.Redirect(r.w, r.req, webutil.FullPath(""), http.StatusFound)
		return
	}

	// No valid "confirm" value, so we need to render the confirmation page
	audit.Log(r.vars.User, audit.ActionSubmissionPending, auditVals)
	r.render(r.hh.confirmTemplate)
}

// serveForm has no logic to handle, just a form to render
func (r *responder) serveForm() {
	r.render(r.hh.formTemplate)
}
