package main

import (
	"net/http"

	"github.com/uoregon-libraries/gopkg/logger"
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

// getIdent combined getIP and the X-Remote-User header for logging
func getIdent(req *http.Request) string {
	var addr = getIP(req)
	var user = req.Header.Get("X-Remote-User")
	if user == "" {
		return addr
	}
	return user + " - " + addr
}

func requestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var sr = statusrecorder.New(w)
		next.ServeHTTP(sr, req)
		logger.Infof("Request: [%s] %s - %d", getIdent(req), req.URL, sr.Status)
	})
}

func requestStaticAssetLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var sr = statusrecorder.New(w)
		next.ServeHTTP(sr, req)
		logger.Debugf("Asset Request: [%s] %s - %d", getIdent(req), req.URL, sr.Status)
	})
}

func fakeUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		req.Header.Set("X-Remote-User", "dummyuser")
		next.ServeHTTP(w, req)
	})
}
