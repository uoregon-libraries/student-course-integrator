// Package global holds top-level variables that should be initialized once and
// then treated as read-only data
package global

import (
	"database/sql"

	"github.com/gorilla/sessions"
	"github.com/uoregon-libraries/student-course-integrator/src/config"
)

// Conf is the global configuration exposed to the entire app, and as such
// should be built precisely once, and never be modified
var Conf *config.Config

// DB is our persistent database connection
var DB *sql.DB

// Store is our session storage backend
var Store sessions.Store
