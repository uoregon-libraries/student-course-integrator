package sciserver

import (
	"net/http"

	"github.com/uoregon-libraries/gopkg/logger"
	"github.com/uoregon-libraries/gopkg/tmpl"
	"github.com/uoregon-libraries/student-course-integrator/src/roles"
)

// responder wraps the writer and request to provide us a simpler approach to
// handling whatever we need to send the client
type responder struct {
	w    http.ResponseWriter
	req  *http.Request
	hh   *homeHandler
	vars *homeVars
}

func respond(w http.ResponseWriter, req *http.Request, hh *homeHandler) *responder {
	var user = getContextUser(req)
	var r = &responder{w: w, req: req, hh: hh, vars: &homeVars{User: user, Form: &form{}, Roles: roles.Roles}}
	var s = getSession(w, req)
	r.vars.Info = s.GetInfoFlash()
	r.vars.Alert = s.GetAlertFlash()

	return r
}

func (r *responder) render(t *tmpl.Template) {
	var err = t.BufferedExecute(r.w, r.vars)
	if err == nil {
		return
	}

	logger.Errorf("Error serving %q: %s", t.Name, err)
	r.render500(err)
}

func (r *responder) render500(err error) {
	r.w.WriteHeader(500)
	r.vars.Alert = "Server error encountered.  Try again or contact support."
	logger.Errorf("Server error: %s", err)
	err = empty.Execute(r.w, r.vars)
	if err != nil {
		logger.Criticalf("Error rendering error page: %s", err)
	}
}
