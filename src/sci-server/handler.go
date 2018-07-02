package sciserver

import (
	"fmt"
	"net/http"

	"github.com/uoregon-libraries/gopkg/tmpl"
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
	var f, err = r.getForm()
	if err != nil {
		render500(r.w, fmt.Errorf("unable to instantiate form: %s", err), &commonVars{})
		return
	}

	if len(f.errors) == 0 {
		// TODO: Do stuff
		return
	}

	// No-go, re-render the form
	var pageVars = &homeVars{}
	pageVars.User = r.user
	pageVars.Alert = fmt.Sprintf("The following errors prevented associating %q with CRN %q: %s",
		f.DuckID, f.CRN, f.errorString())
	pageVars.Form = f
	render(r.tmpl, r.w, pageVars)
}

// serveForm has no logic to handle, just a form to render
func (r *response) serveForm() {
	render(r.tmpl, r.w, &homeVars{commonVars: commonVars{User: r.user}})
}
