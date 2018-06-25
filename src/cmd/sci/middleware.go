package main

import (
	"net/http"

	"github.com/uoregon-libraries/gopkg/logger"
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

func requestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var addr = getIP(req)
		var user = req.Header.Get("X-Remote-User")

		// Initialize the fake writer to a status of 200 - if a status isn't
		// explicitly written, the http library assumes 200
		var sr = &statusRecorder{w, 200}
		next.ServeHTTP(sr, req)

		var ident = addr
		if user != "" {
			ident = user + " - " + addr
		}

		logger.Infof("Request: [%s] %s - %d", ident, req.URL, sr.status)
	})
}
