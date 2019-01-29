package audit

import (
	"encoding/json"
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
	ActionInvalidSubmission Action = "invalid_form_submission"
	ActionSubmissionPending Action = "submission_pending_confirmation"
	ActionRejectSubmission  Action = "submission_rejected"
	ActionConfirmSubmission Action = "submission_confirmed"
)

func Log(u *user.User, action Action, data map[string]string) {
	var content, err = json.Marshal(data)
	// This supposedly isn't possible when marshaling a string->string map, but I
	// would feel dirty not checking errors just in case
	if err != nil {
		logger.Criticalf("Unable to marshal audit message to JSON (error %q).  Message follows.", err.Error())
		logger.Criticalf("user: %q, action: %q, data: %#v", u, action, data)
		return
	}

	var message = string(content)
	var sql = "INSERT INTO audit_logs (`created_at`, `ip`, `login`, `action`, `message`) VALUES(?, ?, ?, ?, ?)"
	_, err = global.DB.Exec(sql, time.Now(), u.IP, u.Login, action, message)
	if err != nil {
		logger.Criticalf("Unable to write audit log to the database (error %q).  Message follows.", err.Error())
		logger.Criticalf("user: %q, action: %q, message: %q", u, action, message)
	}
}
