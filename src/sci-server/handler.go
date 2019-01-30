package sciserver

import (
	"fmt"
	"net/http"
	"strings"

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
	Roles []string   // Roles to be displayed in dropdown
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

// audit writes an audit log in the database which includes the raw form data
// and any form errors
func (r *responder) audit(f *form, a audit.Action) {
	var data = make(map[string]string)
	data["crn"] = f.CRN
	data["duckid"] = f.DuckID
	data["confirm"] = f.Confirm
	data["role"] = "GE"
	if len(f.errors) > 0 {
		data["errors"] = f.errorString()
	}

	audit.Log(r.vars.User, a, data)
}

// ServeHTTP implements http.Handler for homeHandler.  It warms the central IS
// cache, generates a responder, and tells it to figure out what to do
func (h *homeHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	warmCache()
	var r = respond(w, req, h)
	r.route()
}

// route uses the form data to determine where to go, as we actually handle
// every action in a single URI due to my quick-and-dirty approach on the
// initial rush to get the code done.
//
// All audit logging happens here to help with code re-use (e.g., serveForm is
// used in several situations that log different actions)
func (r *responder) route() {
	// The easy case: if it's not a POST request, we simply show the main form
	if strings.ToUpper(r.req.Method) != "POST" {
		r.serveForm()
		return
	}

	// It's a POST.  All posts have the same form data right now, so we parse it
	// before anything else.
	var f, err = r.getForm()
	if err != nil {
		r.render500(fmt.Errorf("unable to instantiate form data: %s", err))
		return
	}

	// One or more form errors: re-render and return
	if len(f.errors) > 0 {
		r.audit(f, audit.ActionInvalidSubmission)
		r.vars.Alert = "Error: " + f.errorString()
		r.serveForm()
		return
	}

	// All other routing is dependent on the value of the "confirm" arg:
	// - 0 means the user said "go back" on the confirmation page
	// - 1 means the user said "confirm" on the confirmation page
	// - no value (or anything unknown) means the user submitted the main form,
	//   but hasn't seen the confirmation page yet
	switch f.Confirm {
	case "0":
		r.audit(f, audit.ActionRejectSubmission)
		r.serveForm()
	case "1":
		r.audit(f, audit.ActionConfirmSubmission)
		r.addGE(f)
	default:
		r.audit(f, audit.ActionSubmissionPending)
		r.serveConfirmForm()
	}
}

func (r *responder) addGE(f *form) {
	var err = enrollment.AddGE(f.CRN, f.GE.BannerID)
	if err != nil {
		r.render500(fmt.Errorf("unable to write enrollment data to database: %s", err))
		return
	}
	var s = getSession(r.w, r.req)
	s.SetInfoFlash(fmt.Sprintf(`%s (%s) added to %s`, f.GE.DisplayName, f.GE.DuckID, f.Course.Description))
	http.Redirect(r.w, r.req, webutil.FullPath(""), http.StatusFound)
}

// serveForm has no logic to handle, just a form to render
func (r *responder) serveForm() {
	r.render(r.hh.formTemplate)
}

// serveConfirmForm just renders the confirmation page
func (r *responder) serveConfirmForm() {
	r.render(r.hh.confirmTemplate)
}
