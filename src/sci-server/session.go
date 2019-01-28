package sciserver

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/uoregon-libraries/gopkg/logger"
	"github.com/uoregon-libraries/student-course-integrator/src/global"
)

// session wraps the gorilla session to provide a better API and to ensure an
// empty session doesn't panic.  The session API exposed here only uses
// strings, not any kind of complex structures.  As such, so long as nobody
// mucks about with the session data outside these APIs, there's no risk of
// pulling a non-string data type and panicking when converting to string.
type session struct {
	s   *sessions.Session
	w   http.ResponseWriter
	req *http.Request
}

// getSession pulls the gorilla session data from our global session store
func getSession(w http.ResponseWriter, req *http.Request) *session {
	var s, err = global.Store.Get(req, "sci")
	if err != nil {
		logger.Errorf("Unable to retrieve session: %s", err)
	}
	return &session{s, w, req}
}

// GetString returns the string value for the given key.  If the value doesn't
// exist in our session data, an empty string is returned.
func (s *session) GetString(key string) string {
	var val = s.s.Values[key]
	if val == nil {
		return ""
	}
	var sval = val.(string)
	return sval
}

// SetString stores the key/value string pair to the session
func (s *session) SetString(key, val string) {
	s.s.Values[key] = val
	var err = s.s.Save(s.req, s.w)
	if err != nil {
		logger.Errorf("Unable to save session data: %s", err)
	}
}

// getFlash grabs the given key as a "flash" value from the session store.  If
// the key doesn't exist, or its value isn't a string, an empty string is
// returned.
func (s *session) getFlash(key string) string {
	var data = s.s.Flashes(key)
	if len(data) > 0 {
		if str, ok := data[0].(string); ok {
			s.s.Save(s.req, s.w)
			return str
		}
	}

	return ""
}

// setFlash stores the key/value pair in the session store as a "flash" value
func (s *session) setFlash(key, val string) {
	s.s.AddFlash(val, key)
	s.s.Save(s.req, s.w)
}

// GetInfoFlash returns the one-time "flash" value for any "info" data set up
func (s *session) GetInfoFlash() string {
	return s.getFlash("info")
}

// SetInfoFlash stores val as the one-time "flash" value for "info" data
func (s *session) SetInfoFlash(val string) {
	s.setFlash("info", val)
}

// GetAlertFlash returns the one-time "flash" value for any "alert" data set up
func (s *session) GetAlertFlash() string {
	return s.getFlash("alert")
}

// SetAlertFlash stores val as the one-time "flash" value for "alert" data
func (s *session) SetAlertFlash(val string) {
	s.setFlash("alert", val)
}
