package sciserver

import (
	"net/http"

	"github.com/uoregon-libraries/gopkg/tmpl"
	"github.com/uoregon-libraries/student-course-integrator/src/data/user"
)

type commonVars struct {
	Alert string
	Info  string
}

// SetAlert implements alertable for the template rendering function
func (v *commonVars) SetAlert(val string) {
	v.Alert = val
}

type homeVars struct {
	commonVars
	User *user.User
}

// homeHandler encapsulates basic data and functionality for handling input and
// rendering output
type homeHandler struct {
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
	var pageVars = &homeVars{User: user}
	render(h.tmpl, w, pageVars)
}
