package sciserver

import (
	"fmt"
	"net/http"

	"github.com/uoregon-libraries/gopkg/tmpl"
	"github.com/uoregon-libraries/student-course-integrator/src/data/audit"
	"github.com/uoregon-libraries/student-course-integrator/src/data/enrollment"
	"github.com/uoregon-libraries/student-course-integrator/src/data/user"
)

type commonVars struct {
	Alert string
	Info  string
	User  *user.User
}

// SetAlert implements alertable for the template rendering function
func (v *commonVars) SetAlert(val string) {
	v.Alert = val
}

type homeVars struct {
	commonVars
	Form form
}

// homeHandler encapsulates basic data and functionality for handling input and
// rendering output
type homeHandler struct {
	formTemplate    *tmpl.Template
	confirmTemplate *tmpl.Template
	successTemplate *tmpl.Template
}

// response wraps the writer and request to provide us a simpler approach to
// handling whatever we need to send the client
type response struct {
	w    http.ResponseWriter
	req  *http.Request
	user *user.User
	hh   *homeHandler
}

func hHome() *homeHandler {
	var r = layout.Clone()
	return &homeHandler{
		formTemplate:    r.MustBuild("home.go.html"),
		confirmTemplate: r.MustBuild("confirm_duckid.go.html"),
		successTemplate: r.MustBuild("enroll_success.go.html"),
	}
}

// ServeHTTP implements http.Handler for homeHandler
func (h *homeHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var user = getContextUser(req)
	var r = &response{w, req, user, h}
	if req.Method == "POST" {
		r.processSubmission()
		return
	}
	r.serveForm()
}

func (r *response) processSubmission() {
	var pageVars = &homeVars{commonVars: commonVars{User: r.user}}
	var f, err = r.getForm()
	if err != nil {
		render500(r.w, fmt.Errorf("unable to instantiate form data: %s", err), pageVars)
		return
	}

	pageVars.Form = f
	var msg = fmt.Sprintf("student %q -> course %q", f.DuckID, f.CRN)

	// Explicit rejection of duckid was requested: re-render the form
	if f.Confirm == "0" {
		msg += "; wrong duckid, re-rendering form"
		audit.Log(r.user, audit.ActionAssociateStudent, msg)
		render(r.hh.formTemplate, r.w, pageVars)
		return
	}

	// Errors: re-render the form
	if len(f.errors) > 0 {
		msg += "; FAILURE: " + f.errorString()
		audit.Log(r.user, audit.ActionAssociateStudent, msg)
		pageVars.Alert = fmt.Sprintf("The following errors prevented associating %q with CRN %q: %s",
			f.DuckID, f.CRN, f.errorString())
		render(r.hh.formTemplate, r.w, pageVars)
		return
	}

	// Require "confirm" to be exactly the string "1" so that we err on the side of not adding students
	if f.Confirm == "1" {
		msg += "; CONFIRMED"
		err = enrollment.AddGTF(f.CRN, f.DuckID)
		if err != nil {
			render500(r.w, fmt.Errorf("unable to write enrollment data to database: %s", err), pageVars)
			return
		}
		audit.Log(r.user, audit.ActionAssociateStudent, msg)
		render(r.hh.successTemplate, r.w, pageVars)
		return
	}

	// No valid "confirm" value, so we need to render the confirmation page
	msg += "; requesting confirmation"
	audit.Log(r.user, audit.ActionAssociateStudent, msg)
	render(r.hh.confirmTemplate, r.w, pageVars)
}

// serveForm has no logic to handle, just a form to render
func (r *response) serveForm() {
	render(r.hh.formTemplate, r.w, &homeVars{commonVars: commonVars{User: r.user}})
}
