package sciserver

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/uoregon-libraries/gopkg/logger"
	"github.com/uoregon-libraries/student-course-integrator/src/global"
)

// session wraps the gorilla session to provide a better API and to ensure an
// empty session doesn't panic
type session struct {
	s   *sessions.Session
	w   http.ResponseWriter
	req *http.Request
}

func getSession(w http.ResponseWriter, req *http.Request) *session {
	var s, err = global.Store.Get(req, "sci")
	if err != nil {
		logger.Errorf("Unable to retrieve session: %s", err)
	}
	return &session{s, w, req}
}

func (s *session) GetString(key string) string {
	var val = s.s.Values[key]
	if val == nil {
		return ""
	}
	var sval = val.(string)
	return sval
}

func (s *session) SetString(key, val string) {
	s.s.Values[key] = val
	var err = s.s.Save(s.req, s.w)
	if err != nil {
		logger.Errorf("Unable to save session data: %s", err)
	}
}