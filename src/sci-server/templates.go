package sciserver

import (
	"net/http"

	"github.com/uoregon-libraries/gopkg/logger"
	"github.com/uoregon-libraries/gopkg/tmpl"
	"github.com/uoregon-libraries/gopkg/webutil"
	"github.com/uoregon-libraries/student-course-integrator/src/version"
)

var (
	// layout holds the base site layout template.  Handlers should clone and use
	// this for parsing their specific page templates
	layout *tmpl.TRoot

	// insufficientPrivileges is a simple page to declare to a user they are not
	// allowed to visit a certain page or perform a certain action
	insufficientPrivileges *tmpl.Template

	// empty holds a simple blank page for rendering the header/footer and often
	// a simple alert-style message
	empty *tmpl.Template
)

// initRootTemplates sets up pre-parsed templates
func initRootTemplates(templatePath string, debug bool) {
	var templateFunctions = tmpl.FuncMap{
		"Version": func() string { return version.Version },
		"Debug":   func() bool { return debug },
	}

	// Set up the layout and then our global templates
	layout = tmpl.Root("layout", templatePath)
	layout.Funcs(webutil.FuncMap)
	layout.Funcs(templateFunctions)
	layout.MustReadPartials("layout.go.html")
	insufficientPrivileges = layout.MustBuild("insufficient-privileges.go.html")
	empty = layout.MustBuild("empty.go.html")
}

type alertable interface {
	SetAlert(string)
}

func render(t *tmpl.Template, w http.ResponseWriter, data alertable) {
	var err = t.BufferedExecute(w, data)
	if err == nil {
		return
	}

	logger.Errorf("Error serving %q: %s", t.Name, err)
	render500(w, err, data)
}

func render500(w http.ResponseWriter, err error, data alertable) {
	w.WriteHeader(500)
	data.SetAlert("Server error encountered.  Try again or contact support.")
	logger.Errorf("Server error: %s", err)
	err = empty.Execute(w, data)
	if err != nil {
		logger.Criticalf("Error rendering error page: %s", err)
	}
}
