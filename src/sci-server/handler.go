package sciserver

import (
	"fmt"
	"net/http"

	"github.com/uoregon-libraries/gopkg/tmpl"
	"github.com/uoregon-libraries/student-course-integrator/src/data/audit"
	"github.com/uoregon-libraries/student-course-integrator/src/data/user"
	"github.com/uoregon-libraries/student-course-integrator/src/data/enrollment"
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
	tmpl *tmpl.Template
}

// response wraps the writer and request to provide us a simpler approach to
// handling whatever we need to send the client
type response struct {
	w    http.ResponseWriter
	req  *http.Request
	user *user.User
	tmpl *tmpl.Template
}

func hHome() *homeHandler {
	var r = layout.Clone()
	return &homeHandler{
		tmpl: r.MustBuild("home.go.html"),
	}
}

// ServeHTTP implements http.Handler for homeHandler
func (h *homeHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var user = getContextUser(req)
	var r = &response{w, req, user, h.tmpl}
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

	if len(f.errors) == 0 {
		err = enrollment.AddGTF(f.CRN, f.DuckID)
		if err != nil {
			render500(r.w, fmt.Errorf("unable to write enrollment data to database: %s", err), pageVars)
			return
		}
		audit.Log(r.user, audit.ActionAssociateStudent, msg)
		return
	}

	// No-go, re-render the form
	msg += "; FAILURE: " + f.errorString()
	audit.Log(r.user, audit.ActionAssociateStudent, msg)
	pageVars.Alert = fmt.Sprintf("The following errors prevented associating %q with CRN %q: %s",
		f.DuckID, f.CRN, f.errorString())
	render(r.tmpl, r.w, pageVars)
}

// serveForm has no logic to handle, just a form to render
func (r *response) serveForm() {
	render(r.tmpl, r.w, &homeVars{commonVars: commonVars{User: r.user}})
}
