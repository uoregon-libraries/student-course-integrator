package audit

import (
	"time"

	"github.com/uoregon-libraries/gopkg/logger"
	"github.com/uoregon-libraries/student-course-integrator/src/data/user"
	"github.com/uoregon-libraries/student-course-integrator/src/global"
)

// Action is just a string, but typed specially for make it easier to have a
// known list of action types
type Action string

// Our full list of valid actions
const (
	ActionAssociateGE Action = "associate GE to course"
)

// Log writes an audit log to the database.  If the database connection fails
// for any reason, we complain loudly to the standard logger ("CRIT" level).
func Log(u *user.User, action Action, message string) {
	var sql = "INSERT INTO audit_logs (`created_at`, `ip`, `login`, `action`, `message`) VALUES(?, ?, ?, ?, ?)"
	var _, err = global.DB.Exec(sql, time.Now(), u.IP, u.Login, action, message)
	if err != nil {
		logger.Criticalf("Unable to write audit log to the database (error %q).  Message follows.", err.Error())
		logger.Criticalf("user: %q, action: %q, message: %q", u, action, message)
	}
}
