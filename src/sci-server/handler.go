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

type homeVars struct {
	Alert string
	Info  string
	User  *user.User
	Form  *form
}

// homeHandler encapsulates basic data and functionality for handling input and
// rendering output
type homeHandler struct {
	formTemplate    *tmpl.Template
	confirmTemplate *tmpl.Template
}

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

	var msg = fmt.Sprintf("GE %q -> course %q", f.DuckID, f.CRN)

	// Explicit rejection of duckid was requested: re-render the form
	if f.Confirm == "0" {
		msg += `; rejected by user ("Go Back" clicked), re-rendering form`
		audit.Log(r.vars.User, audit.ActionAssociateGE, msg)
		r.render(r.hh.formTemplate)
		return
	}

	// Errors: re-render the form
	if len(f.errors) > 0 {
		msg += "; FAILURE: " + f.errorString()
		audit.Log(r.vars.User, audit.ActionAssociateGE, msg)
		r.vars.Alert = fmt.Sprintf("Error: %s", f.errorString())
		r.render(r.hh.formTemplate)
		return
	}

	// Require "confirm" to be exactly the string "1" so that we err on the side of not adding GEs
	if f.Confirm == "1" {
		msg += "; CONFIRMED"
		err = enrollment.AddGE(f.CRN, f.GE.BannerID)
		if err != nil {
			r.render500(fmt.Errorf("unable to write enrollment data to database: %s", err))
			return
		}
		audit.Log(r.vars.User, audit.ActionAssociateGE, msg)
		var s = getSession(r.w, r.req)
		s.SetInfoFlash(fmt.Sprintf(`%s (%s) added to %s`,
			f.GE.DisplayName, f.GE.DuckID, f.Course.Description))
		http.Redirect(r.w, r.req, webutil.FullPath(""), http.StatusFound)
		return
	}

	// No valid "confirm" value, so we need to render the confirmation page
	msg += "; requesting confirmation"
	audit.Log(r.vars.User, audit.ActionAssociateGE, msg)
	r.render(r.hh.confirmTemplate)
}

// serveForm has no logic to handle, just a form to render
func (r *responder) serveForm() {
	r.render(r.hh.formTemplate)
}
