package sciserver

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/uoregon-libraries/gopkg/logger"
	"github.com/uoregon-libraries/student-course-integrator/src/data/user"
	"github.com/uoregon-libraries/student-course-integrator/src/global"
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
	var u = req.Header.Get(global.Conf.AuthHeader)
	if u == "" {
		return "N/A"
	}
	return u
}

// findUser uses the configured auth header to find the user in our database,
// and associates the current IP address
func findUser(req *http.Request) (*user.User, error) {
	var login = getRemoteUser(req)
	var user, err = user.Find(login)
	if err != nil {
		logger.Criticalf("Unable to get user data for %q: %s", login, err)
		return nil, err
	}
	if user.Login == "" {
		logger.Warnf("User with no login tried to get authorization")
	}
	user.IP = getIP(req)
	return user, nil
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
		// Get user session
		var s = getSession(w, req)

		// First check query args for a username
		var username = req.FormValue("debuguser")
		if username != "" {
			s.SetString("debuguser", username)
		}

		// Now check session; if we had an arg above, it will be stored here.  This
		// extra round-trip makes sure our session logic works.  It can be very
		// annoying to debug session data when it seems like it worked, but turns
		// out it was just the hit that had the query arg which worked, while the
		// session had some odd misconfiguration.
		username = s.GetString("debuguser")

		// Still no name?  This app requires a user, so we will just set up a fake name.
		if username == "" {
			username = "dummyuser"
		}
		req.Header.Set(global.Conf.AuthHeader, username)
		next.ServeHTTP(w, req)
	})
}

// getUser pulls the user from our data/user package and stores it on the
// request via context
func getUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var user, err = findUser(req)
		if err != nil {
			logger.Criticalf("Database error trying to look up user: %s", err)
			var r = &responder{w: w, vars: &homeVars{}}
			r.render500(err)
			return
		}
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
			insufficientPrivileges.Execute(w, &homeVars{User: user})
			return
		}
		next.ServeHTTP(w, req)
	})
}
