// Package global holds top-level variables that should be initialized once and
// then treated as read-only data
package global

import (
	"github.com/uoregon-libraries/student-course-integrator/src/config"
)

// Conf is the global configuration exposed to the entire app, and as such
// should be built precisely once, and never be modified
var Conf *config.Config
