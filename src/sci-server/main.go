package sciserver

import (
	"path/filepath"

	"github.com/uoregon-libraries/student-course-integrator/src/global"
)

// Run is the top-level method to execute when wanting to run this as the web service
func Run() {
	initRootTemplates(filepath.Join(global.Opts.Approot, "templates"), global.Conf.Debug)

	var s = &server{
		Approot: global.Opts.Approot,
		Bind:    ":8080",
		Debug:   global.Conf.Debug,
	}
	s.Listen()
}
