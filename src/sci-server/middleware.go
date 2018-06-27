package sciserver

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/uoregon-libraries/gopkg/logger"
	"github.com/uoregon-libraries/student-course-integrator/src/data/user"
	"github.com/uoregon-libraries/student-course-integrator/src/statusrecorder"
)

// nocache is a Middleware function to send back no-cache header
func nocache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Cache-Control", "max-age=0, must-revalidate")
		next.ServeHTTP(w, req)
	})
}

// getIP is a naive implementation to grab a user's IP address.  This won't
// work with all web servers, as the request headers seem to differ widely.
func getIP(req *http.Request) string {
	var addr = req.Header.Get("X-Forwarded-For")
	if addr == "" {
		addr = req.RemoteAddr
	}
	return addr
}

func getRemoteUser(req *http.Request) string {
	var u = req.Header.Get("X-Remote-User")
	if u == "" {
		return "N/A"
	}
	return u
}

// findUser uses the X-Remote-User header to find the user in our database, and
// associates the current IP address
func findUser(req *http.Request) *user.User {
	var login = getRemoteUser(req)
	var user = user.Find(login)
	user.IP = getIP(req)
	return user
}

// getContextUser returns the user retrieved from the request context
func getContextUser(req *http.Request) *user.User {
	var data = context.Get(req, "user")
	var user, ok = data.(*user.User)
	if !ok {
		// If this happens, there's a major logic error somewhere
		panic("unable to read 'user' from context")
	}

	return user
}

func requestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var sr = statusrecorder.New(w)
		next.ServeHTTP(sr, req)
		logger.Infof("Request: [%s] %s - %d", getContextUser(req), req.URL, sr.Status)
	})
}

func requestStaticAssetLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var sr = statusrecorder.New(w)
		next.ServeHTTP(sr, req)
		var login = getRemoteUser(req)
		logger.Debugf("Asset Request: [%s - %s] %s - %d", login, getIP(req), req.URL, sr.Status)
	})
}

func fakeUserLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		req.Header.Set("X-Remote-User", "dummyuser")
		next.ServeHTTP(w, req)
	})
}

// getUser pulls the user from our data/user package and stores it on the
// request via context
func getUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var user = findUser(req)
		context.Set(req, "user", user)
		next.ServeHTTP(w, req)
	})
}

// mustAuth makes sure the logged-in user is allowed to be here, otherwise a
// 403 is returned
func mustAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var user = getContextUser(req)
		if !user.Authorized {
			w.WriteHeader(http.StatusForbidden)
			insufficientPrivileges.Execute(w, nil)
			return
		}
		next.ServeHTTP(w, req)
	})
}
