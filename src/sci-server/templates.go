package sciserver

import (
	"strconv"
	"strings"

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
		"Version":   version.Version,
		"Debug":     func() bool { return debug },
		"stripterm": stripterm,
	}

	// Set up the layout and then our global templates
	layout = tmpl.Root("layout", templatePath)
	layout.Funcs(webutil.FuncMap)
	layout.Funcs(templateFunctions)
	layout.MustReadPartials("layout.go.html")
	insufficientPrivileges = layout.MustBuild("insufficient-privileges.go.html")
	empty = layout.MustBuild("empty.go.html")
}

func stripterm(courseid string) string {
	// Sanity check: the course id *must* be a four-digit year combined with a
	// two-digit term code
	var termcrn = strings.SplitN(courseid, ".", 2)
	if len(termcrn) != 2 {
		return courseid
	}
	var term, crn = termcrn[0], termcrn[1]
	if len(term) != 6 {
		return courseid
	}
	var termnum, err = strconv.Atoi(term)
	if err != nil || termnum < 201600 || termnum > 210000 {
		return courseid
	}

	return crn
}
