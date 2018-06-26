package main

import (
	"net/http"

	"github.com/uoregon-libraries/gopkg/logger"
	"github.com/uoregon-libraries/gopkg/tmpl"
	"github.com/uoregon-libraries/student-course-integrator/src/data/user"
)

type commonVars struct {
	Title string
	Alert string
	Info  string
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
	var pageVars = homeVars{User: user}
	var err = h.tmpl.BufferedExecute(w, pageVars)
	logger.Errorf("Error serving homepage: %s", err)
	if err != nil {
		w.WriteHeader(500)
		pageVars.Alert = "Server error encountered.  Try again or contact support."
		err = empty.Execute(w, pageVars)
		if err != nil {
			logger.Criticalf("Error rendering error page: %s", err)
		}
	}
}
